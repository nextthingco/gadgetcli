package main

import (
	"errors"
	"fmt"
	"github.com/nextthingco/libgadget"
	"github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"path/filepath"
)

// Process the build arguments and execute build
func GadgetInit(args []string, g *libgadget.GadgetContext) error {

	initUu1 := uuid.NewV4()
	initUu2 := uuid.NewV4()

	log.Info("Creating new project:")

	g.WorkingDirectory, _ = filepath.Abs(g.WorkingDirectory)
	initName := filepath.Base(g.WorkingDirectory)

	log.Infof("  in %s", g.WorkingDirectory)

	initConfig := libgadget.TemplateConfig(initName, fmt.Sprintf("%s", initUu1), fmt.Sprintf("%s", initUu2))

	outBytes, err := yaml.Marshal(initConfig)
	if err != nil {
		return err
	}

	fileLocation := fmt.Sprintf("%s/gadget.yml", g.WorkingDirectory)

	gadgetFileExists, err := libgadget.PathExists(fileLocation)

	if gadgetFileExists {

		log.WithFields(log.Fields{
			"function":   "GadgetInit",
			"location":   fileLocation,
			"init-stage": "overwriting file",
		}).Debug("gadget.yml already exists in this location")

		log.Errorf("There's already a config file here [%s]", fileLocation)
		log.Warnf("Remove %s if you'd like to init again", fileLocation)

		err = errors.New("Tried to overwrite pre-existing configuration file")
		return err
	}

	err = ioutil.WriteFile(fileLocation, outBytes, 0644)
	if err != nil {

		log.WithFields(log.Fields{
			"function":   "GadgetInit",
			"location":   fileLocation,
			"init-stage": "writing file",
		}).Debug("This is likely due to a problem with permissions")

		log.Errorf("Failed to create config file [%s]", fileLocation)
		log.Warn("Do you have permission to create a file here?")

		return err
	}

	return err
}
