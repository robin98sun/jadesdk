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
			conf := &SDKConf{}
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
			j.UpdateAPIs()
			if j.Conf != nil && j.Conf.Capabilities != nil && len(j.Conf.Capabilities) > 0 {
				
			} else if j.Conf == nil {
				j.log.Println("ERROR: configuration is nil when updated")
			} else if j.Conf.Capabilities == nil {
				j.log.Println("ERROR: capabilities is nil when configuration is updated")
			} else if len(j.Conf.Capabilities) == 0 {
				j.log.Println("WARNING: capabilities is empty when configuration is updated")
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
