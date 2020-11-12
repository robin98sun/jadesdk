package jadesdk

import (
	// RESTful Server
	"github.com/ant0ine/go-json-rest/rest"
	// "net/http"
	// others
	// "encoding/json"
	// "flag"
	"fmt"
	// "strconv"
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
		timeDecoded := time.Now()
		j.Stats[moduleName].Decoding.AddDuration(timeDecoded.Sub(timeArrive))
		// process the task
		taskThread := func() {
			var result interface{}
			var err error
			if moduleType == AppModuleWorker {
				result, err = moduleInst.(WorkerModuleInstance).Handler(input)
			} else if moduleType == AppModuleAggregator &&
				!j.AggregativeTaskCache.DoesSubtaskExists(req.Task.TaskID, req.Task.SubtaskID) {
				cumulation, previousResults := j.AggregativeTaskCache.GetCumulation(req.Task.TaskID)
				result, err = moduleInst.(AggregatorModuleInstance).Handler(cumulation, previousResults, input)
				j.AggregativeTaskCache.SetSubtaskResult(req.Task.TaskID, req.Task.SubtaskID, input)
				j.AggregativeTaskCache.SetCumulation(req.Task.TaskID, result)
			} else {
				return
			}
			timeExecuted := time.Now()
			executionDuration := timeExecuted.Sub(timeDecoded)
			j.Stats[moduleName].Task.AddDuration(executionDuration)
			j.log.Println(fmt.Sprintf("[%v] done in %v microseconds", moduleName, executionDuration))
			if moduleType == AppModuleAggregator && !j.AggregativeTaskCache.IsTaskDone(req.Task.TaskID) {
				return
			}
			if err != nil {
				errMsg := fmt.Sprintf("[%v] error when processing request: %v", moduleName, err.Error())
				j.log.Println(errMsg)
				j.reportToMaster(req.Task, &Response{
					Error: errMsg,
				}, false)
				timeReportedToMasterAboutError := time.Now()
				j.Stats[moduleName].ReportToMaster.AddDuration(timeReportedToMasterAboutError.Sub(timeExecuted))
			} else {
				// send messages to next hop
				var errorCache map[string]error
				timePoint := time.Now()
				if len(req.To) > 0 {
					j.log.Println(fmt.Sprintf("[%v] forwarding to next hop (%v modules)", moduleName, len(req.To)))
					errorCache = j.sendMessages(req.Task, j.GetSelfInterface(moduleName), req.To, result)
					timeForwarded := time.Now()
					forwardingDuration := timeForwarded.Sub(timePoint)
					j.Stats[moduleName].Forwarding.AddDuration(forwardingDuration)
					j.log.Println(fmt.Sprintf("[%v] forward to next hop completed in %v microseconds", moduleName, forwardingDuration))
					timePoint = timeForwarded
				}
				// report to master
				if j.Conf.MasterNode.IsValid() {
					j.log.Println(fmt.Sprintf("[%v] reporting to the master", moduleName))
					if errorCache == nil {
						j.reportToMaster(req.Task, &Response{
							Status:  "OK",
							Payload: result,
						}, true)
					} else {
						errMsg := ""
						for i, e := range errorCache {
							errMsg += i + "=" + e.Error() + "; "
						}
						errMsg = strings.Trim(errMsg, " ;")
						j.reportToMaster(req.Task, &Response{
							Error:   errMsg,
							Status:  "error when sending messages",
							Payload: result,
						}, true)
					}
					timeReportedToMaster := time.Now()
					reportingDuration := timeReportedToMaster.Sub(timePoint)
					j.Stats[moduleName].ReportToMaster.AddDuration(reportingDuration)
					timePoint = timeReportedToMaster
					j.log.Println(fmt.Sprintf("[%v] report to the master completed in %v microseconds", moduleName, reportingDuration))
				}
			}
		}
		go taskThread()
	}
	return httpHandler
}
