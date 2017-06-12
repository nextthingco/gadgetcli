package main

import (
	"fmt"
	"errors"
	"github.com/satori/go.uuid"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	log "github.com/sirupsen/logrus"
)

func addUsage() error {
	log.Info("Usage: gadget [flags] add [type] [name]")
	log.Info("               *opt        *req   *req ")
	log.Info("Type: service | onboot                 ")
	log.Info("Name: friendly name for container      ")
	
	return errors.New("Incorrect add usage")
}

// Process the build arguments and execute build
func GadgetAdd(args []string, g *GadgetContext) error {

	addUu := uuid.NewV4()
	
	if len(args) != 2 {
		return addUsage()
	}
	
	log.Infof("Adding new %s: \"%s\" ", args[0], args[1])
	
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
		log.Errorf("  %q is not valid command.", args[0])
		return addUsage()
	}
	
	g.Config = CleanConfig(g.Config)
	
	fileLocation := fmt.Sprintf("%s/gadget.yml", g.WorkingDirectory)
	
	outBytes, err := yaml.Marshal(g.Config)
	if err != nil {
			
		log.WithFields(log.Fields{
			"function": "GadgetAdd",
			"location": fileLocation,
			"init-stage": "parsing",
		}).Debug("The config file is probably malformed")


		log.Errorf("Failed to parse config file [%s]", fileLocation)
		log.Warn("Is this a valid gadget.yaml?")
		return err
	}

	err = ioutil.WriteFile(fmt.Sprintf("%s/gadget.yml", g.WorkingDirectory), outBytes, 0644)
	if err != nil {
			
		log.WithFields(log.Fields{
			"function": "GadgetAdd",
			"location": fileLocation,
			"init-stage": "writing file",
		}).Debug("This is likely due to a problem with permissions")


		log.Errorf("Failed to edit config file [%s]", fileLocation)
		log.Warn("Do you have permission to modify this file?")
		
		return err
	}
	
	return err
}
