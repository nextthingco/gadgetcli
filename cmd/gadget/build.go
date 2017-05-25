package main

import (
	"fmt"
	"os/exec"
)

// Process the build arguments and execute build
func build(args []string, g *GadgetContext) {

	g.loadConfig()

	// find docker binary in path
	binary, err := exec.LookPath("docker")
	if err != nil {
		panic(err)
	}

	fmt.Println("[BUILD]  Building:")

	stagedContainers,_ := findStagedContainers(args, append(g.Config.Onboot, g.Config.Services...))

	for _, container := range stagedContainers {
		fmt.Printf("[BUILD]    %s ", container.ImageAlias)

		// use local directory for build
		if container.Directory != "" {
			containerDirectory := fmt.Sprintf("%s/%s", g.WorkingDirectory, container.Directory)
			runLocalCommand(binary,
				"build",
				"--tag",
				container.ImageAlias,
				containerDirectory)
		} else {
			runLocalCommand(binary,
				"pull",
				container.Image)
			runLocalCommand(binary,
				"tag",
				container.Image,
				container.ImageAlias)
		}
		
		fmt.Printf("✔\n")
	}
}
