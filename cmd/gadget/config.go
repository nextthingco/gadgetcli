package main

import (
	"gopkg.in/yaml.v2"
    "fmt"
    "os"
    "runtime"
    "path/filepath"
    "errors"
)

type GadgetContext struct {
	Config				GadgetConfig
	WorkingDirectory	string
}
type GadgetConfig struct {
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

func NewConfig(config []byte) (GadgetConfig, error) {
	g := GadgetConfig{}

	// Parse yaml
	err := yaml.Unmarshal(config, &g)
	if err != nil {
		return g, err
	}

	return g,nil
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
func walkUp ( bottom_dir string ) (string, error) {
	
	var rc error = nil
	
	if _, err := os.Stat(fmt.Sprintf("%s/gadget.yaml", bottom_dir)); err != nil {
		
		// haven't found it
		if isRoot(bottom_dir) || isDriveLetter(bottom_dir) {
			return "", errors.New("config: could not find configuration file")
		} else {
			bottom_dir, rc = walkUp(filepath.Dir(bottom_dir))
		}
	}
	
	return bottom_dir, rc
}
