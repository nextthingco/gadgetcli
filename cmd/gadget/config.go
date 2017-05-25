package main

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
)

type GadgetContext struct {
	Config           GadgetConfig
	Verbose          bool
	WorkingDirectory string
}
type GadgetConfig struct {
	Spec     string
	Name     string
	UUID     string
	Type     string
	Onboot   []GadgetContainer
	Services []GadgetContainer
}

type GadgetContainer struct {
	Name         string
	UUID         string
	Image        string
	Directory    string
	Net          string
	PID          string
	Readonly     bool
	Command      []string `yaml:",flow"`
	Binds        []string `yaml:",flow"`
	Capabilities []string `yaml:",flow"`
	Alias        string   "alias,omitempty"
	ImageAlias   string   "imagealias,omitempty"
}

func TemplateConfig(gName, gUu1, gUu2, gUu3 string) GadgetConfig {
	return GadgetConfig{
		Spec: Version,
		Name: gName,
		UUID: gUu1,
		Type: "docker",
		Onboot: []GadgetContainer{
			{
				Name:    "hello-world",
				Image:   "armhf/hello-world",
				UUID:    gUu2,
			},
		},
	}
}

func ParseConfig(config []byte) (GadgetConfig, error) {
	g := GadgetConfig{}

	// Parse yaml
	err := yaml.Unmarshal(config, &g)
	if err != nil {
		return g, err
	}

	return g, nil
}

func CleanConfig( g GadgetConfig ) GadgetConfig {
	
	// helper function to remove hidden config
	// items before writing the struct out
	
	for i := range g.Onboot {
		g.Onboot[i].Alias = ""
		g.Onboot[i].ImageAlias = ""
	}

	for i := range g.Services {
		g.Services[i].Alias = ""
		g.Services[i].ImageAlias = ""
	}
	
	return g	
}

// helper function for walkup, determines if cwd is '/'
func isRoot(path string) bool {
	if runtime.GOOS != "windows" {
		return path == "/"
	}
	switch len(path) {
	case 1:
		return os.IsPathSeparator(path[0])
	case 3:
		return path[1] == ':' && os.IsPathSeparator(path[2])
	}
	return false
}

// isDriveLetter returns true if path is Windows drive letter (like "c:").
func isDriveLetter(path string) bool {
	if runtime.GOOS != "windows" {
		return false
	}
	return len(path) == 2 && path[1] == ':'
}

// recursive function, returns ("", rc) on failure
// returns ("/path/to/dir", rc) on success
func WalkUp(bottom_dir string) (string, error) {

	var rc error = nil
	// TODO: error checking on path
	//~ bottom_dir,_ = filepath.Abs(bottom_dir)
	//~ ^ moved to loadConfig -- only runs once, usable later

	if _, err := os.Stat(fmt.Sprintf("%s/gadget.yml", bottom_dir)); err != nil {

		// haven't found it
		if isRoot(bottom_dir) || isDriveLetter(bottom_dir) {
			return "", errors.New("[SETUP]  could not find configuration file")
		} else {
			bottom_dir, rc = WalkUp(filepath.Dir(bottom_dir))
		}
	}

	return bottom_dir, rc
}

func (g *GadgetContext) LoadConfig() {

	g.WorkingDirectory, _ = filepath.Abs(g.WorkingDirectory)

	// find and read gadget.yml
	// TODO: this should probably get moved into the NewConfig function
	// TODO: better error checking/reporting. WHY can't the config file be opened?
	var config []byte
	var parseerr error = nil
	var cwderr error = nil

	g.WorkingDirectory, cwderr = WalkUp(g.WorkingDirectory)
	if cwderr == nil {
		// found the config
		fmt.Printf("[SETUP]  Running in directory:\n")
		fmt.Printf("[SETUP]    %s\n", g.WorkingDirectory)

		config, parseerr = ioutil.ReadFile(fmt.Sprintf("%s/gadget.yml", g.WorkingDirectory))
		if parseerr != nil {
			// couldn't read it
			fmt.Printf("[SETUP]  Cannot open config file: %v\n", parseerr)
		}
	} else {
		panic(cwderr)
	}

	// create new config class from gadget.yml output
	// TODO: add error checking here.
	g.Config, parseerr = ParseConfig(config)

	for index, onboot := range g.Config.Onboot {
		onboot.Alias = fmt.Sprintf("%s_%s", onboot.Name, onboot.UUID)
		onboot.ImageAlias = fmt.Sprintf("%s-img", onboot.Alias)
		g.Config.Onboot[index] = onboot
	}
	for index, service := range g.Config.Services {
		service.Alias = fmt.Sprintf("%s_%s", service.Name, service.UUID)
		service.ImageAlias = fmt.Sprintf("%s-img", service.Alias)
		g.Config.Services[index] = service
	}
}
type GadgetContainers []GadgetContainer

func (containers GadgetContainers) Find(name string) (GadgetContainer, error) {
for _,container := range containers {
		if container.Name == name {
			return container, nil
		}
	}
	return GadgetContainer{}, errors.New(fmt.Sprintf("[CONFIG]  could not find container: %s", name))
}
