package jadesdk

import (
// "encoding/json"
)

type MetricsEnv struct {
	Temperature *struct{
		Cpu float64 `json:"cpu,omitempty"`
		Device float64 `json:"device,omitempty"`
		Disk float64 `json:"disk,omitempty"`
	} `json:"TEMPERATURE,omitempty"`

	CPU *struct{
		User int `json:"us,omitempty"`
		Sys int `json:"sy,omitempty"`
		Idle int `json:"id,omitempty"`
		Wait int `json:"wa,omitempty"`
		Stolen int `json:"st,omitempty"`
		Frequency float64 `json:"frequency,omitempty"`
	} `json:"CPU,omitempty"`

	RAM *struct{
		Swapped int64 `json:"swpd,omitempty"`
		Free int64 `json:"free,omitempty"`
		Buffer int64 `json:"buffer,omitempty"`
		Cache int64 `json:"cache,omitempty"`
	} `json:"RAM,omitempty"`

	Swap *struct{
		SwappedIn int64 `json:"si,omitempty"`
		SwappedOut int64 `json:"so,omitempty"`
	} `json:"SWAP,omitempty"`

	IO *struct{
		BlocksReceived int64 `json:"bi,omitempty"`
		BlocksSent int64 `json:"bo,omitempty"`
	} `json:"IO,omitempty"`

	System *struct{
		Interrupts int64 `json:"in,omitempty"`
		ContextSwitches int64 `json:"cs,omitempty"`
	} `json:"SYSTEM,omitempty"`

	Processes *struct{
		Runnable int64 `json:"r,omitempty"`
		Sleeping int64 `json:"cpu,omitempty"`
	} `json:"PROCS,omitempty"`
}