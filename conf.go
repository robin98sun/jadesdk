package jadesdk

import (
	// "fmt"
	"encoding/json"
	"os"
	"strconv"
	"strings"
)

type VideoStream struct {
	URL string
}

type Conf struct {
	TaskID         string
	SubtaskID      string
	SelfNode       *Node
	AggregatorNode *Node
	MasterNode     *Node // to report the completion of sub-task
	TTL            int64 // in milliseconds, if job can not done in TTL, the pod would be forced to clean
	Subtasks       []string
	Module         string
}

// ReadConfFromEnv Read configuration from environment variables
func (j *JadeSDK) ReadConfFromEnv() *Conf {
	conf := &Conf{
		SelfNode:       &Node{},
		AggregatorNode: &Node{},
		MasterNode:     &Node{},
		Subtasks:       []string{},
	}
	conf.TaskID = os.Getenv("JADE_TASKID")
	conf.SubtaskID = os.Getenv("JADE_SUBTASKID")
	conf.SelfNode.Protocol = os.Getenv("JADE_SELFNODE_PROTOCOL")
	conf.SelfNode.Addr = os.Getenv("JADE_SELFNODE_ADDR")
	conf.SelfNode.Port, _ = strconv.Atoi(os.Getenv("JADE_SELFNODE_PORT"))
	conf.MasterNode.Protocol = os.Getenv("JADE_MASTERNODE_PROTOCOL")
	conf.MasterNode.Addr = os.Getenv("JADE_MASTERNODE_ADDR")
	conf.MasterNode.Port, _ = strconv.Atoi(os.Getenv("JADE_MASTERNODE_PORT"))
	conf.AggregatorNode.Protocol = os.Getenv("JADE_AGGREGATORNODE_PROTOCOL")
	conf.AggregatorNode.Addr = os.Getenv("JADE_AGGREGATORNODE_ADDR")
	conf.AggregatorNode.Port, _ = strconv.Atoi(os.Getenv("JADE_AGGREGATORNODE_PORT"))
	conf.TTL, _ = strconv.ParseInt(os.Getenv("JADE_TTL"), 10, 64)
	subtaskstr := os.Getenv("JADE_SUBTASKS")
	conf.Subtasks = strings.Split(subtaskstr, ",")
	conf.Module = os.Getenv("JADE_MODULE")
	j.Conf = conf
	return conf
}

func (j *JadeSDK) PrintConfig() {
	confstr, _ := json.Marshal(j.Conf)
	j.log.Println("Conf:", string(confstr))
}
