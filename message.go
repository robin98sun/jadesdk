package jadesdk

import (
// "encoding/json"
)

type RequestOptions struct {
	EstimatedServiceTime float64 `json:"estimatedServiceTime,omitempty"`
}

// Request message of request
type Request struct {
	Task    *Task           `json:"task,omitempty"`
	From    *Interface      `json:"from,omitempty"`
	To      []*Interface    `json:"to,omitempty"`
	Payload interface{}     `json:"payload,omitempty"`
	Options *RequestOptions `json:"options,omitempty"`
}

// AggregatorEnqueuingMessage message of enqueuing a task
type AggregatorEnqueuingMessage struct {
	TaskKey    string       `json:"taskId,omitempty"`
	SubtaskKey string       `json:"subtaskId,omitempty"`
	ModuleName string       `json:"moduleName,omitempty"`
	Subtasks   []string     `json:"subtasks,omitempty"`
	ReportTo   []*Interface `json:"reportTo,omitempty"`
}

// Response message respond to the request
type Response struct {
	Status  string      `json:"status,omitempty"`
	Error   string      `json:"error,omitempty"`
	Payload interface{} `json:"payload,omitempty"`
}

func newSuccessResponse(payload interface{}) *Response {
	return &Response{
		Status:  "OK",
		Payload: payload,
	}
}

func newErrorResponse(errMsg string) *Response {
	return &Response{
		Error: errMsg,
	}
}

func (j *JadeSDK) sendMessages(task *Task, from *Interface, to []*Interface, message interface{}, options *RequestOptions) (int, map[string]error) {
	errorCache := map[string]error{}
	var requestLength int
	for _, i := range to {
		if !i.IsValid() {
			continue
		}
		_, reqlen, err := j.HTTPCommunicate(
			"send message to "+i.Key(),
			"http", "POST", "/"+i.ModuleName, i.Node,
			&Request{
				Task:    task,
				From:    from,
				Payload: message,
				Options: options,
			},
			0, 30,
		)
		if err != nil {
			errorCache[i.Key()] = err
			j.log.Println("error when sending message sent to "+i.Key(), err.Error())
		} else {
			j.log.Println("message sent to " + i.Key())
			requestLength += reqlen
		}
	}
	if len(errorCache) == 0 {
		return requestLength, nil
	}
	return requestLength, errorCache
}
