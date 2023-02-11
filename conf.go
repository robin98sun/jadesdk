package jadesdk

import (
	"encoding/json"
	"os"
	"strconv"
	ds "uta.edu/aces/jadesdk/data_structure"
)

type VideoStream struct {
	URL string
}

type SDKConf struct {
	SelfNode     *ds.Node
	MasterNode   *ds.Node         // to report the completion of sub-task
	AppModule    string        `json:"appModule,omitempty"`
	AppName      string        `json:"appName,omitempty"`
	AppVersion   string        `json:"appVersion,omitempty"`
	Capabilities []*ds.Capability `json:"capabilities,omitempty"`
}

// ReadConfFromEnv Read configuration from environment variables
func (j *JadeSDK) ReadConfFromEnv() *SDKConf {
	conf := &SDKConf{
		SelfNode:   &ds.Node{},
		MasterNode: &ds.Node{},
	}
	conf.SelfNode.Protocol = os.Getenv("JADE_SELFNODE_PROTOCOL")
	conf.SelfNode.Addr = os.Getenv("JADE_SELFNODE_ADDR")
	conf.SelfNode.Port, _ = strconv.Atoi(os.Getenv("JADE_SELFNODE_PORT"))
	conf.MasterNode.Protocol = os.Getenv("JADE_MASTERNODE_PROTOCOL")
	conf.MasterNode.Addr = os.Getenv("JADE_MASTERNODE_ADDR")
	conf.MasterNode.Port, _ = strconv.Atoi(os.Getenv("JADE_MASTERNODE_PORT"))
	conf.AppName = os.Getenv("JADE_APP_NAME")
	conf.AppModule = os.Getenv("JADE_APP_MODULE")
	conf.AppVersion = os.Getenv("JADE_APP_VERSION")

	all_capabilities := []*ds.Capability{}
	for _, list := range ds.ReadCapabilitiesFromEnv() {
		all_capabilities = append(all_capabilities, list...)
	}
	conf.Capabilities = all_capabilities
	j.Conf = conf
	j.UpdateAPIs(nil)
	return conf
}

func (j *JadeSDK) PrintConfig() {
	confstr, _ := json.Marshal(j.Conf)
	j.log.Println("Conf:", string(confstr))
}

func (c *SDKConf) Merge(newConf *SDKConf) {
	if c == nil || newConf == nil {
		return
	}
	if newConf.SelfNode != nil {
		if c.SelfNode == nil {
			c.SelfNode = newConf.SelfNode
		} else {
			c.SelfNode.Merge(newConf.SelfNode)
		}
	}
	if newConf.MasterNode != nil {
		if c.MasterNode == nil {
			c.MasterNode = newConf.MasterNode
		} else {
			c.MasterNode.Merge(newConf.MasterNode)
		}
	}
	if newConf.AppName != "" {
		c.AppName = newConf.AppName
	}
	if newConf.AppModule != "" {
		c.AppModule = newConf.AppModule
	}
	if newConf.AppVersion != "" {
		c.AppVersion = newConf.AppVersion
	}
	if len(newConf.Capabilities) != 0 {
		c.Capabilities = newConf.Capabilities
	}
}
