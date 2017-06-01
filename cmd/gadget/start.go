package main

import (
	"fmt"
	"errors"
	"strings"
	log "github.com/sirupsen/logrus"
)

// Process the build arguments and execute build
func GadgetStart(args []string, g *GadgetContext) error {
	//~ g.LoadConfig()
	EnsureKeys()

	client, err := GadgetLogin(gadgetPrivKeyLocation)

	if err != nil {
		return err
	}
	
	var startFailed bool = false
	
	log.Info(fmt.Sprintf("[GADGT]  Starting:"))
	stagedContainers,_ := FindStagedContainers(args, append(g.Config.Onboot, g.Config.Services...))
	for _, container := range stagedContainers {
		
		log.Info(fmt.Sprintf("[GADGT]    %s", container.Alias))
		binds := strings.Join( PrependToStrings(container.Binds[:],"-v "), " ")
		commands := strings.Join(container.Command[:]," ")
		
		stdout, stderr, err := RunRemoteCommand(client, "docker create --name", container.Alias, binds, container.ImageAlias, commands)
		
		log.WithFields(log.Fields{
			"function": "GadgetStart",
			"name": container.Alias,
			"start-stage": "create",
		}).Debug(stdout)
		log.WithFields(log.Fields{
			"function": "GadgetStart",
			"name": container.Alias,
			"start-stage": "create",
		}).Debug(stderr)
		
		if err != nil {
			//~ fmt.Printf("✘ ")
			
			// fail loudly, but continue
			
			log.WithFields(log.Fields{
				"function": "GadgetStart",
				"name": container.Alias,
				"start-stage": "create",
			}).Debug("This is likely due to specifying containers for deploying, but trying to start all")


			log.Debug("Failed to create container on Gadget,")
			log.Debug("it might have already been deployed,")
			log.Debug("Or creation otherwise failed")
			
			startFailed = true
			//return err
		} else {
			//~ fmt.Printf("✔ ")
		}

		stdout, stderr, err = RunRemoteCommand(client, "docker start", container.Alias)
		
		log.WithFields(log.Fields{
			"function": "GadgetStart",
			"name": container.Alias,
			"start-stage": "create",
		}).Debug(stdout)
		log.WithFields(log.Fields{
			"function": "GadgetStart",
			"name": container.Alias,
			"start-stage": "create",
		}).Debug(stderr)
		
		if err != nil {
			//~ fmt.Printf("✘\n")
			
			// fail loudly, but continue
			
			log.WithFields(log.Fields{
				"function": "GadgetStart",
				"name": container.Alias,
				"start-stage": "start",
			}).Debug("This is likely due to specifying containers for deploying, but trying to start all")


			log.Error("Failed to start container on Gadget")
			log.Warn("Was the container ever deployed?")
			
			// return err
			startFailed = true
		} else {
			//~ fmt.Printf("✔\n")
		}

	}
	
	if startFailed {
		err = errors.New("Failed to create or start one or more containers")
	}
	
	return err
}
