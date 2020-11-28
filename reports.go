package jadesdk

import (
// "encoding/json"
)

type ReportMessage struct {
	TaskKey    string    `json:"taskId,omitempty"`
	SubtaskKey string    `json:"subtaskId,omitempty"`
	Updates    *Response `json:"updates,omitempty"`
	Stat       *StatItem `json:"stat,omitempty"`
	Status     string    `json:"status,omitempty"`
}

func (j *JadeSDK) reportToMaster(task *Task, updates *Response, doneOrFail bool, stat *StatItem) {
	if !j.Conf.MasterNode.IsValid() {
		j.log.Println("unable to report to master due to no valid master node")
		return
	}
	msg := &ReportMessage{
		TaskKey:    task.TaskID,
		SubtaskKey: task.SubtaskID,
		Updates:    updates,
		Stat:       stat,
	}
	if doneOrFail {
		msg.Status = "done"
	} else {
		msg.Status = "failed"
	}

	j.log.Println("reporting to master:", msg)
	res, err := j.HTTPCommunicate(
		"report to master", "http", "put",
		"/$jade$/app/listener", j.Conf.MasterNode, msg,
		0, 10)
	if err != nil {
		j.log.Println("error when reporting to master:", err.Error())
	} else {
		j.log.Println("reported, master response:", res)
	}
}
