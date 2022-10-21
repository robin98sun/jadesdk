package data_structure

import (
	// "time"
)

type SubTask struct {
	TaskKey           string    `json:"taskId,omitempty"`
	Key               string    `json:"key,omitempty"`
	AppKey            string    `json:"appKey,omitempty"`
	ModuleName        string    `json:"moduleName,omitempty"`
	Fanout            int       `json:"fanout,omitempty"`
	NodeKey           string    `json:"nodeId,omitempty"`
	ResourceKey 	  string    `json:"resourceId,omitempty"`
}

func (t *SubTask) GetKey() string {
	if t.Key == "" {
		t.Key = t.TaskKey + ":" + string(t.ModuleName) + ":" + t.NodeKey + ":" + RandomString()
	}
	return t.Key
}

func NewSubtask(taskKey string, appKey string, moduleName string, nodeKey string, resourceKey string, subtaskKey string) *SubTask {
	newSubtask := &SubTask{
		TaskKey:    taskKey,
		AppKey:     appKey,
		ModuleName: moduleName,
		NodeKey:    nodeKey,
		Fanout:     1,
	}
	if subtaskKey == "" {
		newSubtask.Key = newSubtask.GetKey()
	} else {
		newSubtask.Key = subtaskKey
	}
	newSubtask.ResourceKey = resourceKey
	return newSubtask
}
