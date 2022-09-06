package data_structure

import (
	// "strings"
)

// Application the application specficiations
type Application struct {
	EnvName string                `json:"envName,omitempty"`
	Name    string                `json:"name,omitempty"`
	Version string                `json:"version,omitempty"`
	Owner   string                `json:"owner,omitempty"`
	Modules map[string]*Container `json:"modules,omitempty"`
}

func (a *Application) GetModule(moduleName string) *Container {
	if len(a.Modules) == 0 {
		return nil
	}
	if m, exists := a.Modules[moduleName]; exists {
		return m
	} else {
		return nil
	}
}

func (a *Application) Key() string {
	return a.Owner + ":" + a.Name + ":" + a.Version
}

func (a *Application) valid() bool {
	if a.Name == "" || a.Version == "" || a.Owner == "" || len(a.Modules) == 0 {
		return false
	}
	valid := true
	for _, m := range a.Modules {
		if !m.valid() {
			valid = false
			break
		}
	}
	return valid
}

// predefined module names
type AppModule string
const (
	AppModuleAggregator AppModule = "aggregator"
	AppModuleWorker     AppModule = "worker"
)
