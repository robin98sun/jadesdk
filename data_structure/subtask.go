package data_structure

import (
	// "time"
)

type SubTask struct {
	TaskKey           string    `json:"taskId,omitempty"`
	Key               string    `json:"key,omitempty"`
	AppName           string    `json:"appName,omitempty"`
	ModuleName        string `json:"moduleName,omitempty"`
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

func NewSubtask(taskKey string, appName string, moduleName string, nodeKey string, subtaskKey string, resourceKey string) *SubTask {
	newSubtask := &SubTask{
		TaskKey:    taskKey,
		AppName:    appName,
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
