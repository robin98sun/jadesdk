package data_structure

import (
	"strconv"
	"fmt"
)

type Pod struct {
	NodeKey    string          `json:"nodeId,omitempty"`
	Namespace  string          `json:"namespace,omitempty"`
	PodName    string          `json:"podName,omitempty"`
	Protocol   string          `json:"protocol,omitempty"`
	Addr       string          `json:"addr,omitempty"`
	Port       int             `json:"port,omitempty"`
	Allocation *AllocationUnit `json:"allocation,omitempty"`
	Capacity   *Capacity       `json:"capacity,omitempty"`
	AppKey     string          `json:"appId,omitempty"`
	Container  *Container      `json:"container,omitempty"`
	ModuleName string          `json:"moduleName,omitempty"`
	UID        string		   `json:"uid,omitempty"`
	Key        string          `json:"id,omitempty"`
}

type AllocationUnit struct {
	MinimumCapacity *Capacity `json:"minimumCapacity,omitempty"`
	MaximumCapacity *Capacity `json:"maximumCapacity,omitempty"`
}

func (a *AllocationUnit) GetAvgCPUCores() float64 {
	minCap := int64(0)
	maxCap := int64(0)
	if a.MinimumCapacity != nil {
		minCap = a.MinimumCapacity.CPU
	}
	if a.MaximumCapacity != nil {
		maxCap = a.MaximumCapacity.CPU
	}
	return float64(minCap+maxCap)/float64(2000)
}

func (a *AllocationUnit) Describe() string {
	if a == nil {
		return "nil"
	}

	return fmt.Sprintf("min: (%v), max: (%v)", a.MinimumCapacity.Describe(), a.MaximumCapacity.Describe())
}

func (a *AllocationUnit) Valid() bool {
	return a != nil && a.MaximumCapacity != nil && a.MinimumCapacity != nil && a.MaximumCapacity.GE(a.MinimumCapacity)
}

func (p *Pod) GetKey() string {
	if p.Key == "" {
		p.Key = GenPodKey(p.AppKey, p.ModuleName, p.Addr, p.Port)
	}
	return p.Key
}

func GenPodKey(appKey string, moduleName string, podAddr string, podPort int) string {
	return appKey + ":" + moduleName + "@" + podAddr + ":" + strconv.Itoa(podPort)
}

func NewPod(appKey string, moduleName string, nodeKey string) *Pod {
	return &Pod{
		NodeKey:    nodeKey,
		AppKey:     appKey,
		ModuleName: moduleName,
	}
}

func (p *Pod) GetNodeRepresentation(protocol string) *Node {
	node := NewNode()
	node.Addr = p.Addr
	node.Port = p.Port
	node.Protocol = protocol
	return node
}

func (p *Pod) GetNodeAddr() *Node {
	node := NewNode()
	node.Addr = p.Addr
	node.Port = p.Port
	node.Protocol = p.Protocol
	return node
}

func (p *Pod) CopyForReportTo() *Pod {
	return &Pod{
		Addr: p.Addr,
		Port: p.Port,
	}
}
