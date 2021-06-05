package jadesdk

import (
// "encoding/json"
	"sync"
	"time"
	"strings"
	"encoding/json"
	"net/http"
)

type Addons struct {
	MetricsEnvData *MetricsEnv
	MetricsEnvAPI *Capability
	mutex *sync.Mutex
	started bool
}

func (a *Addons) Lock() {
	a.mutex.Lock()
}

func (a *Addons) Unlock() {
	a.mutex.Unlock()
}

func NewAddons() *Addons {
	return &Addons{
		mutex: &sync.Mutex{},
		started: false,
	}
}

func (a *Addons) GetMetricsEnv() *MetricsEnv {
	return a.MetricsEnvData
}

func (a *Addons) UpdateMetricsEnvAPI( api *Capability) {
	a.Lock()
	a.MetricsEnvAPI = api
	a.Unlock()

	if api != nil && !a.started {
		a.StartDaemon()
	}
}

func (a *Addons) StartDaemon() {
	a.Lock()
	if !a.started {
		a.started = true
		a.Unlock()
		go a.daemon()	
	} else {
		a.Unlock()
	}
}

func (a *Addons) daemon() {
	for {
		time.Sleep(time.Duration(1000 * time.Millisecond))
		a.FetchMetricsEnv()
	}
}

func (a *Addons) FetchMetricsEnv() {
	a.Lock()
	defer a.Unlock()

	if a.MetricsEnvAPI == nil {
		return
	}


	if strings.ToLower(a.MetricsEnvAPI.Type) == "http" {
		// Send the register information to MetricsEnv service
		req, err := http.NewRequest(strings.ToUpper(a.MetricsEnvAPI.Action), a.MetricsEnvAPI.URL, nil)
		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		res, err := client.Do(req)
		if err != nil {
			a.MetricsEnvData = nil
		}

		if res == nil || res.Body == nil {
			a.MetricsEnvData = nil
		}

		// parse the response message of upper node for registering
		resMsg := MetricsEnv{}
		json.NewDecoder(res.Body).Decode(&resMsg)

		a.MetricsEnvData = &resMsg
	}
}
