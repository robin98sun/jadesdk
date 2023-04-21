package jadesdk

import (
// "encoding/json"
	ds "uta.edu/aces/jadesdk/data_structure"
)

type RequestOptions struct {
	EstimatedServiceTime float64 `json:"estimatedServiceTime,omitempty"`
}

// Request message of request
type Request struct {
	Task    *SDKTask           `json:"task,omitempty"`
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
	SLO        *ds.TaskDispatchingItemSLO `json:"slo,omitempty"`
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

func (j *JadeSDK) sendMessages(task *SDKTask, from *Interface, to []*Interface, message interface{}, options *RequestOptions) (int, map[string]error) {
	errorCache := map[string]error{}
	var requestLength int
	for _, i := range to {
		if !i.IsValid() {
			j.log.Printf("the interface is not valid, module: %v, node addr: %v, port: %v, protocol: %v", i.ModuleName, i.Node.Addr, i.Node.Port, i.Node.Protocol)
			continue
		}
		_, reqlen, _, err := j.HTTPCommunicate(
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

func NewAggregatorEnqueuingMessage(taskItem *ds.TaskDispatchingItem, subtasks []*ds.SubtaskOnNode) *AggregatorEnqueuingMessage {
	inst := &AggregatorEnqueuingMessage{
		TaskKey:  taskItem.Task.GetKey(),
		Subtasks: []string{},
		ReportTo: []*Interface{},
	}
	reportTo := taskItem.GetReportToForModule(string(ds.AppModuleAggregator))
	if reportTo != nil && reportTo.Pod != nil {
		inst.ReportTo = append(inst.ReportTo, &Interface{
			Node: reportTo.Pod.GetNodeAddr(),
			ModuleName: string(ds.AppModuleAggregator),
		})
	}
	for _, item := range subtasks {
		inst.Subtasks = append(inst.Subtasks, item.Subtask.GetKey())
	}

	return inst
}

func NewAggregativeWorkerTask(
	taskItem *ds.TaskDispatchingItem,
	worker *ds.SubtaskOnNode,
	protocol string, input interface{}, estimatedServiceTime float64,
) *Request {
	task := taskItem.Task
	reportTo := taskItem.GetReportToForModule(string(ds.AppModuleWorker))
	if reportTo == nil || reportTo.Pod == nil {
		return nil
	}
	req := &Request{
		Task: &SDKTask{
			ModuleName: string(ds.AppModuleWorker),
			TaskID:     task.GetKey(),
			SubtaskID:  worker.Subtask.GetKey(),
		},
		To: []*Interface{
			&Interface{
				Node: &ds.Node{
					Addr:     reportTo.Pod.Addr,
					Port:     reportTo.Pod.Port,
					Protocol: protocol,
				},
				ModuleName: string(ds.AppModuleAggregator),
			},
		},
		Payload: input,
		Options: &RequestOptions{
			EstimatedServiceTime: estimatedServiceTime,
		},
	}
	return req
}


