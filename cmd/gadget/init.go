package main

import (
	"fmt"
	"github.com/satori/go.uuid"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	//~ "os"
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

	err = ioutil.WriteFile(fmt.Sprintf("%s/gadget.yml", g.WorkingDirectory), outBytes, 0644)
	
	//~ fmt.Printf("✔\n")
	return err
}
