package jadesdk

import (
	"fmt"
)

type WorkerHandler func([]byte) (interface{}, error)
type WorkerModuleInstance interface {
	Handler(interface{}) (interface{}, error)
	ShapeInput() interface{}
}

type AggregatorHandler func(cumulation []byte, resultsOfPreviousSubtasks []byte, resultOfCurrentSubtask []byte) (interface{}, error)
type AggregatorModuleInstance interface {
	ShapeResultOfSubtask() interface{}
	ShapeCumulation() interface{}
	Handler(cumulation interface{}, resultsOfPreviousSubtasks []interface{}, resultOfCurrentSubtask interface{}) (interface{}, error)
}

type Task struct {
	ModuleName string `json:"moduleName,omitempty"`
	TaskID     string `json:"taskId,omitempty"`
	SubtaskID  string `json:"subtaskId,omitempty"`
	TTL        int    `json:"ttl,omitempty"` // in milliseconds
}

type Interface struct {
	Node       *Node  `json:"node,omitempty"`
	ModuleName string `json:"moduleName,omitempty"`
}

func (i *Interface) IsValid() bool {
	if i == nil || i.Node == nil {
		return false
	}
	return i.ModuleName != "" && i.Node.IsValid()
}

func (i *Interface) Key() string {
	return i.Node.Key() + "/" + i.ModuleName
}

func (i *Interface) Equal(j *Interface) bool {
	if i == nil || j == nil {
		return false
	}
	return i.Node.Equal(j.Node) && i.ModuleName == j.ModuleName
}

type Node struct {
	Addr     string `json:"addr,omitempty"`
	Port     int    `json:"port,omitempty"`
	Protocol string `json:"protocol,omitempty"`
}

func (n *Node) Key() string {
	return fmt.Sprintf("%v:%v", n.Addr, n.Port)
}

func (n *Node) IsValid() bool {
	return n != nil && n.Addr != "" && n.Port > 0
}

func (n *Node) Equal(m *Node) bool {
	if n == nil || m == nil {
		return false
	} else if n == m {
		return true
	} else if n.Addr == m.Addr &&
		n.Port == m.Port &&
		n.Protocol == m.Protocol {
		return true
	}
	return false
}

type AppModule string

const (
	AppModuleAggregator AppModule = "aggregator"
	AppModuleWorker               = "worker"
)
