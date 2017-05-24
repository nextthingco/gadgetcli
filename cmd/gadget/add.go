package main

import (
	"fmt"
	"github.com/satori/go.uuid"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

func addUsage() {
	fmt.Println("Usage: gadget [flags] add [type] [name]")
	fmt.Println("               *opt        *req   *req ")
	fmt.Println("Type: service | onboot                 ")
	fmt.Println("Name: friendly name for container      ")
	os.Exit(1)
}

// Process the build arguments and execute build
func gadgetAdd(args []string, g *GadgetContext) {

	loadConfig(g)

	addUu := uuid.NewV4()
	
	if len(args) != 2 {
		addUsage()
	}
	
	
	fmt.Printf("[SETUP]  Adding new %s: \"%s\" ", args[0], args[1])
	
	addGadgetContainer := GadgetContainer {	
		Name: 	args[1], 
		Image: 	fmt.Sprintf("%s/%s", g.Config.Name, args[1]),
		UUID: 	fmt.Sprintf("%s", addUu),
	}
	
	// parse arguments
	switch args[0] {
	case "service":
		g.Config.Services = append(g.Config.Services, addGadgetContainer)
	case "onboot":
		g.Config.Onboot = append(g.Config.Onboot, addGadgetContainer)
	default:
		fmt.Printf("%q is not valid command.\n\n", args[0])
		addUsage()
	}
	
	g.Config = cleanConfig(g.Config)

	outBytes, err := yaml.Marshal(g.Config)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = ioutil.WriteFile(fmt.Sprintf("%s/gadget.yml", g.WorkingDirectory), outBytes, 0644)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	
	fmt.Printf("✔\n")
	
}
