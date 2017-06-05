package main

import (
	"fmt"
	"github.com/satori/go.uuid"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"path/filepath"
	log "github.com/sirupsen/logrus"
)

// Process the build arguments and execute build
func GadgetInit(args []string, g *GadgetContext) error {

	initUu1 := uuid.NewV4()
	initUu2 := uuid.NewV4()


	log.Info("[INIT ]  Creating new project:")

	g.WorkingDirectory, _ = filepath.Abs(g.WorkingDirectory)
	initName := filepath.Base(g.WorkingDirectory)
	
	log.Infof("[INIT ]    in %s", g.WorkingDirectory)

	initConfig := TemplateConfig(initName, fmt.Sprintf("%s", initUu1), fmt.Sprintf("%s", initUu2))

	outBytes, err := yaml.Marshal(initConfig)
	if err != nil {
		return err
	}
	
	fileLocation := fmt.Sprintf("%s/gadget.yml", g.WorkingDirectory)
	
	err = ioutil.WriteFile(fileLocation, outBytes, 0644)
	if err != nil {
			
		log.WithFields(log.Fields{
			"function": "GadgetInit",
			"location": fileLocation,
			"init-stage": "writing file",
		}).Debug("This is likely due to a problem with permissions")


		log.Errorf("Failed to create config file [%s]", fileLocation)
		log.Warn("Do you have permission to create a file here?")
		
		return err
	}
	
	return err
}
