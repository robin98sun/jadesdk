package jadesdk

import (
// "encoding/json"
)

func (j *JadeSDK) reportToMaster(task *Task, updates *Response, doneOrFail bool) {
	if !j.Conf.MasterNode.IsValid() {
		j.log.Println("unable to report to master due to no valid master node")
		return
	}
	msg := map[string]interface{}{
		"taskId":    task.TaskID,
		"subtaskId": task.SubtaskID,
		"updates":   updates,
	}
	if doneOrFail {
		msg["status"] = "done"
	} else {
		msg["status"] = "failed"
	}

	j.log.Println("reporting to master:", msg)
	res, err := j.HTTPCommunicate(
		"report to master", "http", "put",
		"/$jade$/app/listener", j.Conf.MasterNode, msg,
		0, 0)
	if err != nil {
		j.log.Println("error when reporting to master:", err.Error())
	} else {
		j.log.Println("reported, master response:", res)
	}
}
