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
	for appModuleName, appModuleInst := range j.AppModules {
		if appModuleInst == nil {
			continue
		}
		path := strings.ReplaceAll(appModuleName, " ", "-")
		path = "/" + path
		routes = append(routes, rest.Post(path, j.createHTTPHandler(appModuleName, appModuleInst)))
	}
	return routes
}

func (j *JadeSDK) createHTTPHandler(moduleName string, module ModuleInstance) func(w rest.ResponseWriter, r *rest.Request) {
	httpHandler := func(w rest.ResponseWriter, r *rest.Request) {
		timeArrive := time.Now()
		input := module.NewInput()
		req, err := decodeRequest(moduleName, input, r)
		if err != nil {
			j.log.Println(err.Error())
			w.WriteJson(newErrorResponse(err.Error()))
			return
		}
		if !j.AllowSelfCycle && req.From != nil && req.From.Equal(j.GetSelfInterfaceOfModule(moduleName)) {
			w.WriteJson(newErrorResponse("self-cycle is not allowed"))
			return
		}
		w.WriteJson(newSuccessResponse("received"))
		timeDecoded := time.Now()
		j.Stats[moduleName].Decoding.AddDuration(timeDecoded.Sub(timeArrive))
		// process the task
		taskThread := func() {
			result, err := module.Handler(input)
			timeExecuted := time.Now()
			j.Stats[moduleName].Task.AddDuration(timeExecuted.Sub(timeDecoded))
			j.log.Println(fmt.Sprintf("[%v] done", moduleName))
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
					errorCache = j.sendMessages(req.Task, j.GetSelfInterfaceOfModule(moduleName), req.To, result)
					j.log.Println(fmt.Sprintf("[%v] forward to next hop completed", moduleName))
					timeForwarded := time.Now()
					j.Stats[moduleName].Forwarding.AddDuration(timeForwarded.Sub(timePoint))
					timePoint = timeForwarded
				}
				// send message to aggregator
				if aggNodeInterface := j.GetAggregator(); aggNodeInterface != nil {
					j.log.Println(fmt.Sprintf("[%v] forwarding to the aggregator", moduleName))
					to := []*Interface{aggNodeInterface}
					aggErrCache := j.sendMessages(req.Task, j.GetSelfInterfaceOfModule(moduleName), to, result)
					if aggErrCache != nil && errorCache != nil {
						for key, value := range aggErrCache {
							errorCache[key] = value
						}
					} else if aggErrCache != nil {
						errorCache = aggErrCache
					}
					timeReportedToAgg := time.Now()
					j.Stats[moduleName].ReportToAggregater.AddDuration(timeReportedToAgg.Sub(timePoint))
					timePoint = timeReportedToAgg
					j.log.Println(fmt.Sprintf("[%v] forward to the aggregator completed", moduleName))
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
					j.Stats[moduleName].ReportToMaster.AddDuration(timeReportedToMaster.Sub(timePoint))
					timePoint = timeReportedToMaster
					j.log.Println(fmt.Sprintf("[%v] report to the master completed", moduleName))
				}
			}
		}
		go taskThread()
	}
	return httpHandler
}
