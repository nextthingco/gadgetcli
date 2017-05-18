package main

import (
	"fmt"
	"os/exec"
)

// Process the build arguments and execute build
func build(args []string, g *GadgetContext) {

	loadConfig(g)

	// find docker binary in path
	binary, err := exec.LookPath("docker")
	if err != nil {
		panic(err)
	}

	// loop through 'onboot' config and build containers
	for _, onboot := range append(g.Config.Onboot, g.Config.Services...) {
		fmt.Println(" ==> Building:", onboot.Name)

		// use local directory for build
		if onboot.Directory != "" {
			containerDirectory := fmt.Sprintf("%s/%s", g.WorkingDirectory, onboot.Directory)
			runLocalCommand(binary,
				"build",
				"--tag",
				onboot.ImageAlias,
				containerDirectory)
		} else {
			runLocalCommand(binary,
				"pull",
				onboot.Image)
			runLocalCommand(binary,
				"tag",
				onboot.Image,
				onboot.ImageAlias)
		}
	}
}
