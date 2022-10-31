package data_structure

import (
	"strconv"
	// "uta.edu/aces/jadesdk"
)

type Node struct {
	Addr         string `json:"address,omitempty"`
	Port            int    `json:"port,omitempty"`
	Protocol        string `json:"protocol,omitempty"`
	Token           string `json:"token,omitempty"`
	Namespace       string `json:"namespace,omitempty"`
	PodName         string `json:"podName,omitempty"`
	Hostname        string `json:"hostname,omitempty"`
	ServiceExternal string `json:"serviceExternal,omitempty"`
	// Roles           []string `json:"roles,omitempty"` // master(the top most), linker(only one subnode), aggregator(has subnodes), worker(no subnode)
}

// NewNode construct a new node instance with default values
func NewNode() *Node {
	return &Node{
		Addr:   "",
		Port:      0,
		Protocol:  "http",
		Token:     "",
		// Namespace: "jade-app",
		// PodName:   "",
		Hostname:  "",
	}
}

func NodeFromMap(nodeMap map[string]interface{}) *Node {
	if len(nodeMap) == 0 {
		return nil
	}
	node := NewNode()
	if v, e := nodeMap["address"]; e {
		node.Addr = v.(string)
	}
	if v, e := nodeMap["port"]; e {
		node.Port= v.(int)
	}
	if v, e := nodeMap["protocol"]; e {
		node.Protocol = v.(string)
	}
	if v, e := nodeMap["token"]; e {
		node.Token = v.(string)
	}
	if v, e := nodeMap["namespace"]; e {
		node.Namespace = v.(string)
	}
	if v, e := nodeMap["podName"]; e {
		node.PodName = v.(string)
	}
	if v, e := nodeMap["hostname"]; e {
		node.Hostname = v.(string)
	}
	if v, e := nodeMap["serviceExternal"]; e {
		node.ServiceExternal = v.(string)
	}
	return node
}

func (n *Node) GetSDKNode() *Node {
	sdkNode := &Node{
		Addr:     n.Addr,
		Port:     n.Port,
		Protocol: n.Protocol,
	}
	return sdkNode
}

// Key is used to store node in cache
func (n *Node) Key() string {
	if !n.IsAddrEmpty() {
		return n.URL()
	} else if n.Hostname != "" && n.Namespace != "" && n.PodName != "" {
		return n.Hostname + ":" + n.Namespace + ":" + n.PodName
	}
	return ""
}

func (n *Node) GetKey() string {
	return n.Key()
}

// IsAddrEmpty tells whether a node address is meaningless
func (n *Node) IsAddrEmpty() bool {
	return n.Addr == "" || n.Port == 0
}

// URL is the base http/https url for the node to access
func (n *Node) URL() string {
	if n.Addr != "" && n.Port != 0 && n.Protocol != "" {
		return n.Protocol + "://" + n.Addr + ":" + strconv.Itoa(n.Port)
	} else if n.Addr != "" && n.Port != 0 {
		return n.Addr + ":" + strconv.Itoa(n.Port)
	}
	return ""
}

func (n *Node) GetAddr() string {
	if n.Addr != "" && n.Port != 0 {
		return n.Addr + ":" + strconv.Itoa(n.Port)
	}
	return ""
}

// MiniNode is to get a copy of minimum content to transfer on the network
func (n *Node) MiniNode() *Node {
	return &Node{
		Addr:  n.Addr,
		Port:     n.Port,
		Protocol: n.Protocol,
		Token:    n.Token,
	}
}

func (n *Node) Merge(newNode *Node) {
	if n == nil || newNode == nil {
		return
	}
	if newNode.Addr != "" {
		n.Addr = newNode.Addr
	}
	if newNode.Port != 0 {
		n.Port = newNode.Port
	}
	if newNode.Protocol != "" {
		n.Protocol = newNode.Protocol
	}
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
