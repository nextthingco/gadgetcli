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
	From         string
	Net          string
	PID          string
	Readonly     bool
	Command      []string 	`yaml:",flow"`
	Binds        []string 	`yaml:",flow"`
	Capabilities []string 	`yaml:",flow"`
	Alias        string		"alias,omitempty"
	ImageAlias   string		"imagealias,omitempty"
}

func templateConfig(gName, gUu1, gUu2, gUu3 string) GadgetConfig {
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
		Services: []GadgetContainer{
			{
				Name:    "gadget-dmesg",
				Image:   "gadget/dmesg",
				UUID:    gUu3,
				From:    "armhf/alpine",
				Command: []string{"dmesg", "-wH"},
			},
		},
	}
}

func parseConfig(config []byte) (GadgetConfig, error) {
	g := GadgetConfig{}

	// Parse yaml
	err := yaml.Unmarshal(config, &g)
	if err != nil {
		return g, err
	}

	return g, nil
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
func walkUp(bottom_dir string) (string, error) {

	var rc error = nil
	// TODO: error checking on path
	//~ bottom_dir,_ = filepath.Abs(bottom_dir)
	//~ ^ moved to loadConfig -- only runs once, usable later

	if _, err := os.Stat(fmt.Sprintf("%s/gadget.yml", bottom_dir)); err != nil {

		// haven't found it
		if isRoot(bottom_dir) || isDriveLetter(bottom_dir) {
			return "", errors.New("config: could not find configuration file")
		} else {
			bottom_dir, rc = walkUp(filepath.Dir(bottom_dir))
		}
	}

	return bottom_dir, rc
}

func loadConfig(g *GadgetContext) {

	g.WorkingDirectory, _ = filepath.Abs(g.WorkingDirectory)

	// find and read gadget.yml
	// TODO: this should probably get moved into the NewConfig function
	// TODO: better error checking/reporting. WHY can't the config file be opened?
	var config []byte
	var parseerr error = nil
	var cwderr error = nil

	g.WorkingDirectory, cwderr = walkUp(g.WorkingDirectory)
	if cwderr == nil {
		// found the config
		fmt.Printf("Running in directory: %s\n", g.WorkingDirectory)

		config, parseerr = ioutil.ReadFile(fmt.Sprintf("%s/gadget.yml", g.WorkingDirectory))
		if parseerr != nil {
			// couldn't read it
			fmt.Printf("Cannot open config file: %v\n", parseerr)
		}
	} else {
		panic(cwderr)
	}

	// create new config class from gadget.yml output
	// TODO: add error checking here.
	g.Config, parseerr = parseConfig(config)

	for index, onboot := range g.Config.Onboot {
		onboot.Alias = fmt.Sprintf("%s_%s", onboot.Name, onboot.UUID)
		onboot.ImageAlias = fmt.Sprintf("%s-img", onboot.Alias)
		g.Config.Onboot[index] = onboot
	}
}
