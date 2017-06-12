package main

import (
	"fmt"
	"errors"
	"../libgadget"
	log "github.com/sirupsen/logrus"
)

// Process the build arguments and execute build
func GadgetStatus(args []string, g *libgadget.GadgetContext) error {
	
	libgadget.EnsureKeys()

	client, err := libgadget.GadgetLogin(libgadget.GadgetPrivKeyLocation)

	if err != nil {
		return err
	}

	log.Info("Retrieving status:")
	
	stagedContainers,_ := libgadget.FindStagedContainers(args, append(g.Config.Onboot, g.Config.Services...))
	
	statusFailed := false
	
	for _, container := range stagedContainers {
		commandFormat := `docker ps -a --filter=ancestor=%s --format "{{.Image}} {{.Command}} {{.Status}}"`
		cmd := fmt.Sprintf(commandFormat, container.ImageAlias)
		
		stdout, stderr, err := libgadget.RunRemoteCommand(client, cmd)
		if err != nil {
			
			log.WithFields(log.Fields{
				"function": "GadgetStatus",
				"name": container.Alias,
				"start-stage": "docker ps -a",
			}).Debug("This is likely due to specifying containers for deploying, but trying get status for all")


			log.Error("Failed to fetch container status on Gadget")
			log.Warn("Was the container ever deployed?")
			
			statusFailed = true
		}
		
		log.Info(stdout)
		log.Debug(stderr)
		
	}
	
	if statusFailed {
		err = errors.New("Failed to get status on one or more containers")
	}
	
	return err
}
