package jadesdk

import (
	"sync"
)

type JadeSDKStatus string

const (
	JadeSDKStatusReady         JadeSDKStatus = "ready"
	JadeSDKStatusReporting                   = "reporting"
	JadeSDKStatusError                       = "error"
	JadeSDKStatusTaskDone                    = "taskDone"
	JadeSDKStatusTaskInitiated               = "taskInitiated"
	JadeSDKStatusTaskRunning                 = "taskRunning"
	JadeSDKStatusTaskIdle                    = "idle"
)

type JadeSDK struct {
	Conf           *Conf
	Status         JadeSDKStatus
	AppModules     map[string]ModuleInstance
	mutex          *sync.Mutex
	log            *Logger
	AllowSelfCycle bool
	Stats          map[string]*Stat
}

func NewJadeSDK() *JadeSDK {
	return &JadeSDK{
		Status:     JadeSDKStatusReady,
		mutex:      &sync.Mutex{},
		log:        &Logger{},
		AppModules: map[string]ModuleInstance{},
		Stats:      map[string]*Stat{},
	}
}

func (j *JadeSDK) Verbose(on bool) {
	j.log.Enable = on
}

func (j *JadeSDK) PreventSelfCycle(on bool) {
	j.AllowSelfCycle = !on
}

func (j *JadeSDK) Lock() {
	j.mutex.Lock()
}

func (j *JadeSDK) Unlock() {
	j.mutex.Unlock()
}

func (j *JadeSDK) SetModule(moduleName string, inst ModuleInstance) {
	if moduleName == "" {
		return
	}
	if inst != nil {
		j.AppModules[moduleName] = inst
		j.Stats[moduleName] = newStat()
	} else if _, exists := j.AppModules[moduleName]; exists {
		delete(j.AppModules, moduleName)
	}
}

func (j *JadeSDK) SetWorkerModule(inst ModuleInstance) {
	j.SetModule("worker", inst)
}
func (j *JadeSDK) SetAggregatorModule(inst ModuleInstance) {
	j.SetModule("aggregator", inst)
}

func (j *JadeSDK) GetModule(moduleName string) ModuleInstance {
	if moduleName == "" || len(j.AppModules) == 0 {
		return nil
	}
	if handler, exists := j.AppModules[moduleName]; exists {
		return handler
	}
	return nil
}

func (j *JadeSDK) GetWorkerModule() ModuleInstance {
	return j.GetModule("worker")
}
func (j *JadeSDK) GetAggregatorModule() ModuleInstance {
	return j.GetModule("aggregator")
}

func (j *JadeSDK) GetAggregator() *Interface {
	if aggModule := j.GetAggregatorModule(); j.Conf.AggregatorNode.IsValid() && aggModule != nil {
		aggNodeInterface := &Interface{
			Node:       j.Conf.AggregatorNode,
			ModuleName: "aggregator",
		}
		return aggNodeInterface
	}
	return nil
}
func (j *JadeSDK) GetSelfInterfaceOfModule(moduleName string) *Interface {
	if m := j.GetModule(moduleName); j.Conf.SelfNode.IsValid() && m != nil {
		i := &Interface{
			Node:       j.Conf.SelfNode,
			ModuleName: moduleName,
		}
		return i
	}
	return nil
}

func (j *JadeSDK) ModuleCount() int {
	return len(j.AppModules)
}
