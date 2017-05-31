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
func GadgetInit(args []string, g *GadgetContext) error {

	initUu1 := uuid.NewV4()
	initUu2 := uuid.NewV4()


	fmt.Println("[INIT ]  Creating new project:")

	g.WorkingDirectory, _ = filepath.Abs(g.WorkingDirectory)
	initName := filepath.Base(g.WorkingDirectory)
	
	fmt.Printf("[INIT ]    in %s ", g.WorkingDirectory)

	initConfig := TemplateConfig(initName, fmt.Sprintf("%s", initUu1), fmt.Sprintf("%s", initUu2))

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
	
	fmt.Printf("✔\n")
	return nil	
}
