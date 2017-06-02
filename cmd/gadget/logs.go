package main

import (
	"fmt"
	"errors"
	log "github.com/sirupsen/logrus"
)

// Process the build arguments and execute build
func GadgetLogs(args []string, g *GadgetContext) error {
	
	EnsureKeys()

	client, err := GadgetLogin(gadgetPrivKeyLocation)

	if err != nil {
		return err
	}

	log.Info("[GADGT]  Retrieving logs:")
	
	logsFailed := false
	
	stagedContainers,_ := FindStagedContainers(args, append(g.Config.Onboot, g.Config.Services...))
	for _, container := range stagedContainers {
		commandFormat := `docker logs %s`
		cmd := fmt.Sprintf(commandFormat, container.Alias)
		
		stdout, stderr, err := RunRemoteCommand(client, cmd)
		
		if err != nil {
			
			// fail loudly, but continue
			
			log.WithFields(log.Fields{
				"function": "GadgetLogs",
				"name": container.Alias,
				"start-stage": "docker logs",
			}).Debug("This is likely due to specifying containers for deploying, but trying to fetch all logs")

			log.Error("Failed to fetch container logs on Gadget")
			log.Warn("Was the container ever deployed?")
			
			logsFailed = true
		}
		
		log.Infof("  Begin: %s", container.Name)
		log.Infof("\n%s\n", stdout)
		log.Debugf("\n%s\n", stderr)
		log.Infof("  End: %s", container.Name)
		
	}
	
	if logsFailed {
		err = errors.New("Failed to fetch logs for one or more containers")
	}
	
	return err
}
