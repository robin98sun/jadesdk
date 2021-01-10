package jadesdk

import (
	"errors"
	"fmt"
	"github.com/ant0ine/go-json-rest/rest"
	"strings"
	"time"
)

func (j *JadeSDK) createAppRoutes() []*rest.Route {
	routes := []*rest.Route{}

	for appModuleName, appModuleInst := range j.AggregatorModules {
		if appModuleInst == nil {
			continue
		}
		path := strings.ReplaceAll(appModuleName, " ", "-")
		path = "/" + path
		routes = append(routes, rest.Post(path, j.createHTTPHandler(appModuleName, appModuleInst, AppModuleAggregator)))
	}

	for appModuleName, appModuleInst := range j.WorkerModules {
		if appModuleInst == nil {
			continue
		}
		path := strings.ReplaceAll(appModuleName, " ", "-")
		path = "/" + path
		routes = append(routes, rest.Post(path, j.createHTTPHandler(appModuleName, appModuleInst, AppModuleWorker)))
	}

	return routes
}

func (j *JadeSDK) createHTTPHandler(moduleName string, moduleInst interface{}, moduleType AppModule) func(w rest.ResponseWriter, r *rest.Request) {
	httpHandler := func(w rest.ResponseWriter, r *rest.Request) {
		timeArrive := time.Now()
		var req *Request
		var err error
		var input interface{}
		var statItem *StatItem = &StatItem{}
		if moduleType == AppModuleWorker {
			workerInput := moduleInst.(WorkerModuleInstance).ShapeInput()
			req, err = decodeRequest(moduleName, workerInput, r)
			input = workerInput
		} else if moduleType == AppModuleAggregator {
			aggregatorInput := moduleInst.(AggregatorModuleInstance).ShapeResultOfSubtask()
			req, err = decodeRequest(moduleName, aggregatorInput, r)
			input = aggregatorInput
		}

		if err != nil {
			j.log.Println(err.Error())
			w.WriteJson(newErrorResponse(err.Error()))
			return
		}
		if !j.AllowSelfCycle && req.From != nil &&
			j.WorkerModuleExists(moduleName) && req.From.Equal(j.GetSelfInterface(moduleName)) {
			w.WriteJson(newErrorResponse("self-cycle is not allowed"))
			return
		}
		w.WriteJson(newSuccessResponse("received"))
		// process the task
		taskThread := func() {
			var result interface{}
			var err error
			if moduleType == AppModuleWorker {
				if req.Options != nil && req.Options.EstimatedServiceTime > 0 {
					// backdoor for fake service time
					time.Sleep(time.Duration(req.Options.EstimatedServiceTime) * time.Millisecond)
					j.log.Printf("running in simlation mode, sleep %v milliseconds according to the service time setting", req.Options.EstimatedServiceTime)
				} else {
					result, err = moduleInst.(WorkerModuleInstance).Handler(input)
				}
			} else if moduleType == AppModuleAggregator {
				j.log.Printf("checking subtask[%v] of task[%v] in cache", req.Task.SubtaskID, req.Task.TaskID)
				if j.AggregativeTaskCache.DoesSubtaskExist(req.Task.TaskID, req.Task.SubtaskID) {
					j.log.Printf("received the result of subtask [%v] of task[%v]", req.Task.SubtaskID, req.Task.TaskID)
					cumulation, previousResults := j.AggregativeTaskCache.GetCumulation(req.Task.TaskID)
					// j.log.Printf("for subtask [%v] of task[%v], cumulation: %v, with {%v} previous results", req.Task.SubtaskID, req.Task.TaskID, cumulation, len(previousResults))
					j.AggregativeTaskCache.SetSubtaskResult(req.Task.TaskID, req.Task.SubtaskID, input)
					j.log.Printf("saved result of subtask [%v] of task[%v]", req.Task.SubtaskID, req.Task.TaskID)

					if req.Options != nil && req.Options.EstimatedServiceTime > 0 {
						// backdoor for fake service time
						j.log.Printf("the subtask [%v] is running in the simulation mode with service time [%v], its result is ignored",
							req.Task.SubtaskID, req.Options.EstimatedServiceTime,
						)
					} else {
						TryCatchBlock{
							Try: func() {
								result, err = moduleInst.(AggregatorModuleInstance).Handler(cumulation, previousResults, input)
							},
							Catch: func(e Exception) {
								errMsg := fmt.Sprintf(
									"Application crashed when executing module[%v] on subtask[%v] of task[%v]: %v",
									moduleName, req.Task.SubtaskID, req.Task.TaskID,
									e,
								)
								err = errors.New(errMsg)
								j.log.Printf("ERROR: %v", errMsg)
							},
						}.Do()

						j.log.Printf("[%v] executed aggregator for subtask[%v] of task[%v], err: %v", moduleName, req.Task.SubtaskID, req.Task.TaskID, err)
					}
					if err == nil {
						j.AggregativeTaskCache.SetCumulation(req.Task.TaskID, result)
					}
				} else {
					j.log.Printf("[%v] ERROR: subtask [%v] of task[%v] does not exist", moduleName, req.Task.SubtaskID, req.Task.TaskID)
					return
				}
			}
			timeExecuted := time.Now()
			executionDuration := timeExecuted.Sub(timeArrive)
			j.Stats[moduleName].Service.AddDuration(executionDuration)
			statItem.Service = executionDuration
			j.log.Println(fmt.Sprintf("[%v] done in %v milliseconds", moduleName, executionDuration/time.Millisecond))
			if moduleType == AppModuleAggregator && !j.AggregativeTaskCache.IsTaskDone(req.Task.TaskID) {
				j.log.Printf("[%v] still waiting for more subtasks of task[%v]", moduleName, req.Task.TaskID)
				return
			} else if moduleType == AppModuleAggregator {
				j.log.Printf("[%v] all subtasks of task[%v] is done", moduleName, req.Task.TaskID)
			}
			// prepare to forward and report
			task := req.Task
			if moduleType == AppModuleAggregator {
				task = j.AggregativeTaskCache.GetAggregatorTask(req.Task.TaskID)
			}
			if err != nil {
				errMsg := fmt.Sprintf("[%v] error when processing result of subtask[%v], ERROR: %v", moduleName, req.Task.SubtaskID, err.Error())
				j.log.Println(errMsg)
				j.reportToMaster(task, &Response{
					Error: errMsg,
				}, false, statItem)
				timeReportedToMasterAboutError := time.Now()
				j.Stats[moduleName].ReportToMaster.AddDuration(timeReportedToMasterAboutError.Sub(timeExecuted))
			} else {
				// send messages to next hop
				var errorCache map[string]error
				var reqlen int
				timePoint := time.Now()
				var reportTo []*Interface
				// report to upper tier aggregators
				if moduleType == AppModuleAggregator {
					reportTo = j.AggregativeTaskCache.GetReportTo(task.TaskID)
				}
				if len(req.To) > 0 {
					reportTo = append(reportTo, req.To...)
				}
				if len(reportTo) > 0 {
					j.log.Println(fmt.Sprintf("[%v] forwarding to next hop (%v nodes)", moduleName, len(reportTo)))
					reqlen, errorCache = j.sendMessages(task, j.GetSelfInterface(moduleName), reportTo, result, req.Options)
					timeForwarded := time.Now()
					forwardingDuration := timeForwarded.Sub(timePoint)
					j.Stats[moduleName].Forwarding.AddDuration(forwardingDuration)
					j.Stats[moduleName].PackageSize.AddNumber(int64(reqlen))
					statItem.Forwarding = forwardingDuration
					statItem.PackageSize = int64(reqlen)
					j.log.Println(fmt.Sprintf("[%v] forward to next hop completed in %v milliseconds", moduleName, forwardingDuration/time.Millisecond))
					timePoint = timeForwarded
				}

				// report to master
				if j.Conf.MasterNode.IsValid() {
					j.log.Println(fmt.Sprintf("[%v] reporting app result to the master", moduleName))
					if errorCache == nil {
						j.reportToMaster(task, &Response{
							Status:  "OK",
							Payload: result,
						}, true, statItem)
					} else {
						errMsg := ""
						for i, e := range errorCache {
							errMsg += i + "=" + e.Error() + "; "
						}
						errMsg = strings.Trim(errMsg, " ;")
						j.reportToMaster(task, &Response{
							Error:   errMsg,
							Status:  "error when sending messages",
							Payload: result,
						}, true, statItem)
					}
					timeReportedToMaster := time.Now()
					reportingDuration := timeReportedToMaster.Sub(timePoint)
					j.Stats[moduleName].ReportToMaster.AddDuration(reportingDuration)
					timePoint = timeReportedToMaster
					j.log.Println(fmt.Sprintf("[%v] report to the master completed in %v milliseconds", moduleName, reportingDuration/time.Millisecond))
				}
			}
			// clear the task cache
			if moduleType == AppModuleAggregator && j.AggregativeTaskCache.IsTaskDone(task.TaskID) {
				j.AggregativeTaskCache.CleanTask(task.TaskID)
			}
		}
		go taskThread()
	}
	return httpHandler
}
