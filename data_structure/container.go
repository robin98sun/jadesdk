package data_structure

import (
	"strings"
)


// Container specify one container of an application
type Container struct {
	Image    string      `json:"image,omitempty"`
	Port     int         `json:"port,omitempty"`
	Protocol string      `json:"protocol,omitempty"`
	Addr     string      `json:"addr,omitempty"`
	Input    interface{} `json:"input,omitempty"`
	ID       string   `json:"id,omitempty"`
}

func (c *Container) valid() bool {
	if c.Image == "" {
		return false
	}
	if c.Port == 0 {
		return false
	}
	if c.Protocol == "" || (strings.ToLower(c.Protocol) != "tcp" && strings.ToLower(c.Protocol) != "udp") {
		return false
	}
	return true
}

func (c *Container) SetISAInImage(isa string) {
	if isa == "" || c.Image == "" {
		return
	}

	separater := "--"
	parts := strings.Split(c.Image, separater)
	new_image := ""
	targetPosition := 1
	if len(parts) > 1 {
		targetPosition = len(parts) - 1
	}
	for i:=0;i<targetPosition;i++ {
		new_image += parts[i] + separater
	}
	new_image += isa
	c.Image = new_image
}

func (c *Container) Copy() *Container {
	return &Container{
		Image: c.Image,
		Port: c.Port,
		Protocol: c.Protocol,
		Addr: c.Addr,
		Input: c.Input,
		ID: c.ID,
	}
}