package main

import (
	"fmt"
	"github.com/satori/go.uuid"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
)

// Process the build arguments and execute build
func gadgetInit(args []string, g *GadgetContext) {

	initUu1 := uuid.NewV4()
	initUu2 := uuid.NewV4()
	initUu3 := uuid.NewV4()

	g.WorkingDirectory, _ = filepath.Abs(g.WorkingDirectory)
	initName := filepath.Base(g.WorkingDirectory)

	initConfig := templateConfig(initName, fmt.Sprintf("%s", initUu1), fmt.Sprintf("%s", initUu2), fmt.Sprintf("%s", initUu3))

	outBytes, err := yaml.Marshal(initConfig)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = ioutil.WriteFile(fmt.Sprintf("%s/gadget.yml", g.WorkingDirectory), outBytes, 0644)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}
