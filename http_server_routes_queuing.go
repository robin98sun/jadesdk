package jadesdk

import (
	"github.com/ant0ine/go-json-rest/rest"
)

func (j *JadeSDK) createQueuingInterfaces() []*rest.Route {
	routes := []*rest.Route{
		rest.Put("/$jade$/enqueueAggregativeTask", func(w rest.ResponseWriter, r *rest.Request) {
			msg := &AggregatorEnqueuingMessage{}
			_, err := decodeRawRequestBody(r, msg)
			if err != nil {
				w.WriteJson(newErrorResponse(err.Error()))
				return
			}
			j.AggregativeTaskCache.EnqueueAggregativeTask(msg)
			w.WriteJson(newSuccessResponse(j.Status))
		}),
	}
	return routes
}
