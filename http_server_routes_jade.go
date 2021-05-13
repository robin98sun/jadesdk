package jadesdk

import (
	// RESTful Server
	"github.com/ant0ine/go-json-rest/rest"
	// "net/http"
	// others
	// "encoding/json"
	// "flag"
	// "fmt"
	// "strconv"
	// "strings"
	"runtime"
)

func (j *JadeSDK) createJadeInterfaces() []*rest.Route {
	routes := []*rest.Route{
		rest.Put("/$jade$/config", func(w rest.ResponseWriter, r *rest.Request) {
			conf := &Conf{}
			_, err := decodeRawRequestBody(r, conf)
			if err != nil {
				w.WriteJson(newErrorResponse(err.Error()))
				return
			}
			j.Lock()
			if j.Conf == nil {
				j.Conf = conf
			} else {
				j.Conf.Merge(conf)
			}
			if j.Conf != nil && j.Conf.Capabilities != nil && len(j.Conf.Capabilities) > 0 {
				for i, cap := range j.Conf.Capabilities {
					j.log.Printf("capability[%v] name: %v, value: %v, api: %v, type: %v, action: %v, url: %v", 
						i, cap.Name, cap.Value, cap.API, cap.Type, cap.Action, cap.URL,
					)
					if cap.Parameters != nil && len(cap.Parameters) > 0 {
						for k, param := range cap.Parameters {
							j.log.Printf("   param[%v] name: %v, type: %v", k, param.Name, param.Type)
						}
					}
				}
			}
			j.Unlock()
			w.WriteJson(newSuccessResponse(j.Status))
		}),
		rest.Get("/$jade$/config", func(w rest.ResponseWriter, r *rest.Request) {
			w.WriteJson(newSuccessResponse(j.Conf))
		}),
		rest.Get("/$jade$/status", func(w rest.ResponseWriter, r *rest.Request) {
			w.WriteJson(newSuccessResponse(j.Status))
		}),
		rest.Get("/$jade$/taskCache", func(w rest.ResponseWriter, r *rest.Request) {
			w.WriteJson(newSuccessResponse(j.AggregativeTaskCache))
		}),
		rest.Get("/$jade$/stat", func(w rest.ResponseWriter, r *rest.Request) {
			w.WriteJson(newSuccessResponse(j.Stats))
		}),
		rest.Delete("/$jade$/GC", func(w rest.ResponseWriter, r *rest.Request) {
			runtime.GC()
			w.WriteJson(newSuccessResponse("OK"))
		}),
	}
	return routes
}
