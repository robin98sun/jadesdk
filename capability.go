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

func ReadCapabilitiesFromEnv() []*Capability {
	capalist := []*Capability{}
	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)
		envName := pair[0]
		envValue := pair[1]

		nameParts := strings.SplitN(envName, "_", -1)
		if len(nameParts) < 3 || nameParts[0] != "JADE" {
			continue
		}
		if nameParts[1] == "CAPABILITY" && len(nameParts) == 4 {
			i, err := strconv.Atoi(nameParts[2])
			if err == nil {
				if i >= len(capalist) {
					for x := len(capalist); x <= i; x++ {
						capability := *NewCapability()
						capalist = append(capalist, &capability)
					}
				}
				switch nameParts[3] {
				case "NAME":
					capalist[i].Name = envValue
				case "API":
					capalist[i].API = envValue
				}
				capalist[i].ParseAPI()
			}
		}
	}
	if len(capalist) > 0 {
		return capalist
	}
	return nil
}
