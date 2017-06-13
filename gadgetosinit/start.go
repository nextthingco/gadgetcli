package main

import (
	"errors"
	"github.com/nextthingco/libgadget"
	log "github.com/sirupsen/logrus"
)

// Process the build arguments and execute build
func GadgetOsInit(args []string, g *libgadget.GadgetContext) error {
	
	libgadget.EnsureKeys()

	client, err := libgadget.GadgetLogin(libgadget.GadgetPrivKeyLocation)

	if err != nil {
		return err
	}
	
	var initFailed bool = false
	
	log.Info("Starting:")
	
	for _, container := range g.Config.Onboot {
		
		log.Infof("  %s", container.Alias)
		stdout, stderr, err := libgadget.RunRemoteCommand(client, "docker run --restart=onfailure:3", container.Alias)
		
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
			// fail loudly, but continue
			
			log.WithFields(log.Fields{
				"function": "GadgetStart",
				"name": container.Alias,
				"start-stage": "create",
			}).Debug("This is likely due to specifying containers for deploying, but trying to start all")


			log.Errorf("Failed to start '%s' on Gadget", container.Name)
			log.Warn("Was it ever deployed?")
			
			initFailed = true
		} else {
			log.Info("    - started")
		}

	}
	
	for _, container := range g.Config.Services {
		
		log.Infof("  %s", container.Alias)
		stdout, stderr, err := libgadget.RunRemoteCommand(client, "docker run --restart=on-failure", container.Alias)
		
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
			// fail loudly, but continue
			
			log.WithFields(log.Fields{
				"function": "GadgetStart",
				"name": container.Alias,
				"start-stage": "create",
			}).Debug("This is likely due to specifying containers for deploying, but trying to start all")


			log.Errorf("Failed to start '%s' on Gadget", container.Name)
			log.Warn("Was it ever deployed?")
			
			initFailed = true
		} else {
			log.Info("    - started")
		}

	}
	
	if initFailed {
		err = errors.New("Failed to create or start one or more containers")
	}
	
	return err
}
