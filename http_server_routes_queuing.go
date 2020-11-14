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
				j.log.Printf("ERROR: failed to dequeue aggregator message: %v", err.Error())
				w.WriteJson(newErrorResponse(err.Error()))
				return
			}
			j.log.Printf("enqueue aggregative task: %v", msg)
			j.AggregativeTaskCache.EnqueueAggregativeTask(msg)
			w.WriteJson(newSuccessResponse(j.Status))
		}),
	}
	return routes
}
