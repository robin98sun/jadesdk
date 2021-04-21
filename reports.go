package jadesdk

import (
	// "encoding/json"
	"errors"
	"time"
)

type ReportMessage struct {
	TaskKey    string    `json:"taskId,omitempty"`
	SubtaskKey string    `json:"subtaskId,omitempty"`
	Updates    *Response `json:"updates,omitempty"`
	Stat       *StatItem `json:"stat,omitempty"`
	Status     string    `json:"status,omitempty"`
	Node       *Node     `json:"node,omitempty"`
}

func (j *JadeSDK) reportToMaster(task *Task, updates *Response, doneOrFail bool, stat *StatItem) (interface{}, int, error) {
	if !j.Conf.MasterNode.IsValid() {
		j.log.Println("unable to report to master due to no valid master node")
		return nil, 0, errors.New("unable to report to master due to no valid master node")
	}
	msg := &ReportMessage{
		TaskKey:    task.TaskID,
		SubtaskKey: task.SubtaskID,
		Updates:    updates,
		Stat:       stat,
		Node:       j.Conf.SelfNode,
	}
	if doneOrFail {
		msg.Status = "done"
	} else {
		msg.Status = "failed"
	}

	j.log.Println("reporting to master:", msg)
	res, reqlen, _, _, err := j.HTTPCommunicate(
		"report to master", "http", "put",
		"/$jade$/app/listener", j.Conf.MasterNode, msg,
		0, 10, time.Time{})
	if err != nil {
		j.log.Println("error when reporting to master:", err.Error())
	} else {
		j.log.Println("reported, master response:", res)
	}
	return res, reqlen, err
}
