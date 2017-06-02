package main

import (
	"fmt"
	"os/exec"
	log "github.com/sirupsen/logrus"
)

// Process the build arguments and execute build
func GadgetBuild(args []string, g *GadgetContext) error {

	// find docker binary in path
	binary, err := exec.LookPath("docker")
	if err != nil {
		return err
	}

	log.Info("Building:")

	stagedContainers,_ := FindStagedContainers(args, append(g.Config.Onboot, g.Config.Services...))

	for _, container := range stagedContainers {
		log.Infof("  %s", container.ImageAlias)

		// use local directory for build
		if container.Directory != "" {
			containerDirectory := fmt.Sprintf("%s/%s", g.WorkingDirectory, container.Directory)
			RunLocalCommand(binary,
				"build",
				"--tag",
				container.ImageAlias,
				containerDirectory)
		} else {
			RunLocalCommand(binary,
				"pull",
				container.Image)
			RunLocalCommand(binary,
				"tag",
				container.Image,
				container.ImageAlias)
		}
		
	}
	
	return err
}
