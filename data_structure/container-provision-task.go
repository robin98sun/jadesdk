package data_structure

import (
	// "time"
	// "fmt"
	// "strings"
)


type ContainerProvisionTask struct {
	Containers  []*ContainerProvisionItem	`json:"containers,omitempty"`
}


type ContainerProvisionItem struct {

	AppName  	 string   					`json:"appName,omitempty"`
	ModuleName   string						`json:"moduleName,omitempty"`
	Version      string						`json:"version,omitempty"`
	Owner        string 					`json:"owner,omitempty"`

	Container 	 *Container					`json:"container,omitempty"`
	Replicas     *ResourceReplicaOptions	`json:"replicas,omitempty"`
}


type ResourceReplicaOptions struct {
	Count    			int 		`json:"count,omitempty"`
	UnifiedCapacity 	*Capacity 	`json:"unifiedCapacity,omitempty"`
}

func (i *ContainerProvisionItem) GetAppKey() string {
	return i.Owner + ":" + i.AppName + ":" + i.Version
}

func (i *ContainerProvisionItem) GetModuleKey() string {
	return i.Owner + ":" + i.AppName + ":" + i.ModuleName + ":" + i.Version
}

