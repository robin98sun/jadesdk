package jadesdk

import (
	// "encoding/json"
	"errors"
	ds "uta.edu/aces/jadesdk/data_structure"
	"fmt"
)

type ReportMessage struct {
	TaskKey    string    `json:"taskId,omitempty"`
	SubtaskKey string    `json:"subtaskId,omitempty"`
	Updates    *Response `json:"updates,omitempty"`
	Stat       *StatItem `json:"stat,omitempty"`
	Status     string    `json:"status,omitempty"`
	Node       *ds.Node     `json:"node,omitempty"`
	MetricsEnv *MetricsEnv `json:"metricsEnv,omitempty"`
}

func (j *JadeSDK) reportToMaster(task *SDKTask, updates *Response, doneOrFail bool, stat *StatItem) (interface{}, int, error) {
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
		MetricsEnv: j.Addons.GetMetricsEnv(),
	}
	if doneOrFail {
		msg.Status = "done"
	} else {
		msg.Status = "failed"
	}

	return j.SendReportMessageToJadelet(task.TaskID, j.Conf.MasterNode, msg)
	
}

func (j *JadeSDK) SendReportMessageToJadelet(taskKey string, node *ds.Node, msg *ReportMessage) (interface{}, int, error) {
	j.log.Println("reporting to master:", msg)
	res, reqlen, _, err := j.HTTPCommunicate(
		"report to master for task["+taskKey+"]", "http", "put",
		"/$jade$/app/listener", node, msg,
		0, 10,
	)
	if err != nil {
		j.log.Println(fmt.Printf("error when reporting task[%v] to master: %v", taskKey, err.Error()))
	} else {
		j.log.Println(fmt.Printf("reported task[%v], master response: %v", taskKey, res))
	}
	return res, reqlen, err
}
