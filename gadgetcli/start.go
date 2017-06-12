package main

import (
	"errors"
	"strings"
	"../libgadget"
	log "github.com/sirupsen/logrus"
)

// Process the build arguments and execute build
func GadgetStart(args []string, g *libgadget.GadgetContext) error {
	
	libgadget.EnsureKeys()

	client, err := libgadget.GadgetLogin(libgadget.GadgetPrivKeyLocation)

	if err != nil {
		return err
	}
	
	var startFailed bool = false
	
	log.Info("Starting:")
	stagedContainers,_ := libgadget.FindStagedContainers(args, append(g.Config.Onboot, g.Config.Services...))
	for _, container := range stagedContainers {
		
		log.Infof("  %s", container.Alias)
		binds := strings.Join( libgadget.PrependToStrings(container.Binds[:],"-v "), " ")
		commands := strings.Join(container.Command[:]," ")
		
		stdout, stderr, err := libgadget.RunRemoteCommand(client, "docker create --name", container.Alias, binds, container.ImageAlias, commands)
		
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


			log.Debugf("Failed to create %s on Gadget,", container.Alias)
			log.Debug("it might have already been deployed,")
			log.Debug("Or creation otherwise failed")
			
			startFailed = true
		}

		stdout, stderr, err = libgadget.RunRemoteCommand(client, "docker start", container.Alias)
		
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
			
			startFailed = true
		} else {
			log.Info("    - started")
		}

	}
	
	if startFailed {
		err = errors.New("Failed to create or start one or more containers")
	}
	
	return err
}
