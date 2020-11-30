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
	}
	return routes
}
