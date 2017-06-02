package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
)

// Process the build arguments and execute build
func GadgetStatus(args []string, g *GadgetContext) error {
	
	EnsureKeys()

	client, err := GadgetLogin(gadgetPrivKeyLocation)

	if err != nil {
		return err
	}

	log.Info("[GADGT]  Retrieving status:")
	
	stagedContainers,_ := FindStagedContainers(args, append(g.Config.Onboot, g.Config.Services...))
	
	for _, container := range stagedContainers {
		commandFormat := `docker ps -a --filter=ancestor=%s --format "{{.Image}} {{.Command}} {{.Status}}"`
		cmd := fmt.Sprintf(commandFormat, container.ImageAlias)
		RunRemoteCommand(client, cmd)
		if err != nil {
			return err
		}
	}
	return nil
}
