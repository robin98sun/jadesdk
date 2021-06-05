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
	Conf                 *Conf
	Status               JadeSDKStatus
	WorkerModules        map[string]WorkerModuleInstance
	AggregatorModules    map[string]AggregatorModuleInstance
	mutex                *sync.Mutex
	log                  *Logger
	AllowSelfCycle       bool
	Stats                map[string]*Stat
	AggregativeTaskCache *AggregativeTaskCache
	Addons 				 *Addons
}

func NewJadeSDK() *JadeSDK {
	return &JadeSDK{
		Status:               JadeSDKStatusReady,
		mutex:                &sync.Mutex{},
		log:                  &Logger{},
		WorkerModules:        map[string]WorkerModuleInstance{},
		AggregatorModules:    map[string]AggregatorModuleInstance{},
		Stats:                map[string]*Stat{},
		AggregativeTaskCache: NewAggregativeTaskCache(),
		Addons:               NewAddons(),
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

func (j *JadeSDK) GetWorkerModule(moduleName string) WorkerModuleInstance {
	if moduleName == "" || len(j.WorkerModules) == 0 {
		return nil
	}
	if moduleInst, exists := j.WorkerModules[moduleName]; exists {
		return moduleInst
	}
	return nil
}
func (j *JadeSDK) GetDefaultWorkerModule() WorkerModuleInstance {
	return j.GetWorkerModule(string(AppModuleWorker))
}
func (j *JadeSDK) SetWorkerModule(moduleName string, inst WorkerModuleInstance) {
	if inst != nil {
		j.WorkerModules[moduleName] = inst
		j.Stats[moduleName] = NewStat()
	} else if _, exists := j.WorkerModules[moduleName]; exists {
		delete(j.WorkerModules, moduleName)
	}
}
func (j *JadeSDK) SetDefaultWorkerModule(inst WorkerModuleInstance) {
	j.SetWorkerModule(string(AppModuleWorker), inst)
}
func (j *JadeSDK) WorkerModuleExists(moduleName string) bool {
	if moduleName == "" {
		return false
	}
	if _, e := j.WorkerModules[moduleName]; e {
		return true
	}
	return false
}

func (j *JadeSDK) GetAggregatorModule(moduleName string) AggregatorModuleInstance {
	if moduleName == "" || len(j.WorkerModules) == 0 {
		return nil
	}
	if moduleInst, exists := j.AggregatorModules[moduleName]; exists {
		return moduleInst
	}
	return nil
}
func (j *JadeSDK) GetDefaultAggregatorModule() AggregatorModuleInstance {
	return j.GetAggregatorModule(string(AppModuleAggregator))
}
func (j *JadeSDK) SetAggregatorModule(moduleName string, inst AggregatorModuleInstance) {
	if inst != nil {
		j.AggregatorModules[moduleName] = inst
		j.Stats[moduleName] = NewStat()
	} else if _, exists := j.AggregatorModules[moduleName]; exists {
		delete(j.AggregatorModules, moduleName)
	}
}
func (j *JadeSDK) SetDefaultAggregatorModule(inst AggregatorModuleInstance) {
	j.SetAggregatorModule(string(AppModuleAggregator), inst)
}
func (j *JadeSDK) AggregatorModuleExists(moduleName string) bool {
	if moduleName == "" {
		return false
	}
	if _, e := j.AggregatorModules[moduleName]; e {
		return true
	}
	return false
}

func (j *JadeSDK) GetSelfInterface(moduleName string) *Interface {
	if j.Conf.SelfNode.IsValid() {
		return &Interface{
			Node:       j.Conf.SelfNode,
			ModuleName: moduleName,
		}
	}
	return nil
}

func (j *JadeSDK) ModuleCount() int {
	return len(j.WorkerModules) + len(j.AggregatorModules)
}

func (j *JadeSDK) UpdateAPIs() {
	if j == nil {
		return
	}
	var metricsEnvApi *Capability
	if j.Conf != nil && j.Conf.Capabilities != nil && len(j.Conf.Capabilities) > 0 {
		// log & inspect the capabilities (APIs)
		// prepare APIs
		for i, cap := range j.Conf.Capabilities {
			j.log.Printf("capability[%v] name: %v, value: %v, api: %v, type: %v, action: %v, url: %v", 
				i, cap.Name, cap.Value, cap.API, cap.Type, cap.Action, cap.URL,
			)
			if cap.Parameters != nil && len(cap.Parameters) > 0 {
				for k, param := range cap.Parameters {
					j.log.Printf("   param[%v] name: %v, type: %v", k, param.Name, param.Type)
				}
			}
			// capture APIs
			if cap.Name == "jade-addon-env-metrics" {
				metricsEnvApi = cap
				j.log.Printf("==>captured API for MetricsEnv, action: %v, URL: %v", cap.Action, cap.URL)
			}
		}
	}
	// update APIs
	j.Addons.UpdateMetricsEnvAPI(metricsEnvApi)
}
