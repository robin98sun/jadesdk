package jadesdk

import (
	"encoding/json"
	"os"
	"strconv"
)

type VideoStream struct {
	URL string
}

type Conf struct {
	SelfNode   *Node
	MasterNode *Node  // to report the completion of sub-task
	AppModule  string `json:"appModule,omitempty"`
	AppName    string `json:"appName,omitempty"`
	AppVersion string `json:"appVersion,omitempty"`
}

// ReadConfFromEnv Read configuration from environment variables
func (j *JadeSDK) ReadConfFromEnv() *Conf {
	conf := &Conf{
		SelfNode:   &Node{},
		MasterNode: &Node{},
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
	j.Conf = conf
	return conf
}

func (j *JadeSDK) PrintConfig() {
	confstr, _ := json.Marshal(j.Conf)
	j.log.Println("Conf:", string(confstr))
}

func (c *Conf) Merge(newConf *Conf) {
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
}
