package main

import (
	"gopkg.in/yaml.v2"
)

type Gadget struct {
	Spec				string
	Name				string
	UUID				string
	Type				string
	Onboot []GadgetContainer
	Services []GadgetContainer
}

type GadgetContainer struct {
	Name				string
	UUID				string
	Image				string
	From				string
	Net					string
	PID					string
	Readonly			bool
	Command				[]string
	Binds				[]string
	Capabilities		[]string
}

func NewConfig(config []byte) (*Gadget, error) {
	g := Gadget{}

	// Parse yaml
	err := yaml.Unmarshal(config, &g)
	if err != nil {
		return &g, err
	}

	return &g,nil
}
