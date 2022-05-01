package jadesdk

import (
	"os"
	"strconv"
	"strings"
)

type Capability struct {
	Name       string   `json:"name,omitempty"`
	API        string   `json:"api,omitempty"`
	Type       string   `json:"type,omitempty"`
	Action     string   `json:"action,omitempty"`
	Value      string   `json:"value,omitempty"`
	URL        string   `json:"url,omitempty"`
	Parameters []*Param `json:"parameters,omitempty"`
}

type Param struct {
	Name string `json:"name,omitempty"`
	Type string `json:"type,omitempty"`
}

func (c *Capability) GetKey() string {
	key := c.Name + ";" + c.API + ";" + c.Type + ";" + c.Action + ";" + c.Value + ";" + c.URL
	if len(c.Parameters) > 0 {
		for _, p := range c.Parameters {
			key += ";" + p.Name + ":" + p.Type
		}
	}
	return key
}

// Examples:
// {
// 	"name": "location",
// 	"api": ";value://NYC"
// },
// {
// 	"name": "accuracy",
// 	"api": ";value://city"
// },
// {
// 	"name": "avg_temperature",
// 	"api": "get;http://176.0.0.3/temperature;timespan=int"
// },
// {
// 	"name": "current_temperature",
// 	"api": "get;http://176.0.0.3/temperature"
// }

// NewCapability construct a capability instance with default values
func NewCapability() *Capability {
	return &Capability{
		Name: "",
		API:  "",
	}
}

// IsStatic tells whether a capability contain a static value
func (c *Capability) IsStatic() bool {
	if c.API != "" && c.API[0:9] == ":value://" {
		return true
	}
	return false
}

// MiniCapability generate a mini instance to transfer in the network
func (c *Capability) MiniCapability() *Capability {
	if c == nil {
		return NewCapability()
	}
	if c.Type == "" {
		return &Capability{
			Name: c.Name,
			API:  c.API,
		}
	} else if c.Type == "value" {
		return &Capability{
			Name:  c.Name,
			Value: c.Value,
		}
	} else {
		return &Capability{
			Name: c.Name,
		}
	}
}

// ParseAPI parse the API into structured data
func (c *Capability) ParseAPI() {
	if c.API == "" {
		return
	}
	parts := strings.Split(c.API, ";")
	if len(parts) < 2 {
		return
	}
	tmpUrl := parts[1]
	tmpParts := strings.Split(tmpUrl, ":")

	if len(tmpParts) < 2 {
		return
	}
	c.Type = tmpParts[0]
	if c.Type == "value" {
		c.Value = tmpParts[1]
		if len(c.Value) > 2 && c.Value[0:2] == "//" {
			c.Value = c.Value[2:]
		}
	} else {
		c.Action = parts[0]
		c.URL = tmpUrl
		if len(parts) > 2 {
			params := parts[2]
			paramsParts := strings.Split(params, ",")
			for _, paramStr := range paramsParts {
				paramstrParts := strings.Split(paramStr, "=")
				if len(paramstrParts) == 2 {
					c.Parameters = append(c.Parameters, &Param{
						Name: paramstrParts[0],
						Type: paramstrParts[1],
					})
				}
			}
		}
	}
}

func (c *Capability) Equal(n *Capability) bool {
	if c.Name == n.Name {
		if c.Value == "" && n.Value == "" {
			return true
		} else if c.Value == n.Value {
			return true
		}
	}
	return false
}

func FindMissingCapabilities(availableList []*Capability, requiredList []*Capability) []*Capability {
	if len(availableList) == 0 {
		return requiredList
	}
	if len(requiredList) == 0 {
		return nil
	}
	cache := make(map[string]*Capability)
	for _, avl := range availableList {
		cache[avl.Name] = avl
	}
	missing := []*Capability{}
	for _, req := range requiredList {
		if avl, exists := cache[req.Name]; exists {
			if !avl.Equal(req) {
				missing = append(missing, req)
			}
		} else {
			missing = append(missing, req)
		}
	}
	return missing
}

func AnyCapabilityExists(availableList []*Capability, requiredList []*Capability) bool {
	if len(availableList) == 0 {
		return false
	}
	if len(requiredList) == 0 {
		return true
	}
	cache := make(map[string]*Capability)
	for _, avl := range availableList {
		cache[avl.Name] = avl
	}
	for _, req := range requiredList {
		if avl, exists := cache[req.Name]; exists {
			if avl.Equal(req) {
				return true
			}
		}
	}
	return false
}

func AnyCapabilityMissing(availableList []*Capability, requiredList []*Capability) bool {
	if len(availableList) == 0 {
		return true
	}
	if len(requiredList) == 0 {
		return false
	}
	cache := make(map[string]*Capability)
	for _, avl := range availableList {
		cache[avl.Name] = avl
	}
	for _, req := range requiredList {
		if avl, exists := cache[req.Name]; exists {
			if !avl.Equal(req) {
				return true
			}
		} else {
			return true
		}
	}
	return false
}

func ReadCapabilitiesFromEnv() map[string][]*Capability {
	cap_set := make(map[string][]*Capability)
	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)
		envName := pair[0]
		envValue := pair[1]

		nameParts := strings.SplitN(envName, "_", -1)
		if len(nameParts) < 5 || nameParts[0] != "JADE" {
			continue
		}
		if nameParts[1] == "CAPABILITY" && len(nameParts) == 5 {
			i, err := strconv.Atoi(nameParts[3])
			if err == nil {
				cap_type := strings.ToLower(nameParts[2])
				if _, e := cap_set[cap_type]; !e {
					cap_set[cap_type] = []*Capability{}
				}
				if i >= len(cap_set[cap_type]) {
					for x := len(cap_set[cap_type]); x <= i; x++ {
						capability := *NewCapability()
						cap_set[cap_type] = append(cap_set[cap_type], &capability)
					}
				}
				switch nameParts[4] {
				case "NAME":
					cap_set[cap_type][i].Name = envValue
				case "API":
					cap_set[cap_type][i].API = envValue
				}
				cap_set[cap_type][i].ParseAPI()
			}
		}
	}
	if len(cap_set) > 0 {
		return cap_set
	}
	return nil
}
