package main


import (
	"errors"
	"github.com/nextthingco/libgadget"
	log "github.com/sirupsen/logrus"
)

// Process the build arguments and execute build
func GadgetDelete(args []string, g *libgadget.GadgetContext) error {
	
	libgadget.EnsureKeys()

	client, err := libgadget.GadgetLogin(libgadget.GadgetPrivKeyLocation)

	if err != nil {
		return err
	}

	log.Info("Deleting:")

	stagedContainers,_ := libgadget.FindStagedContainers(args, append(g.Config.Onboot, g.Config.Services...))
	
	deleteFailed := false
	
	for _, container := range stagedContainers {
		log.Infof("  %s", container.ImageAlias)
		
		stdout, stderr, err := libgadget.RunRemoteCommand(client, "docker", "rmi", container.ImageAlias)
		
		log.WithFields(log.Fields{
			"function": "GadgetDelete",
			"name": container.Alias,
			"delete-stage": "rmi",
		}).Debug(stdout)
		log.WithFields(log.Fields{
			"function": "GadgetDelete",
			"name": container.Alias,
			"delete-stage": "rmi",
		}).Debug(stderr)
		
		if err != nil {
			
			log.WithFields(log.Fields{
				"function": "GadgetDelete",
				"name": container.Alias,
				"delete-stage": "rmi",
			}).Debug("This is likely due to specifying containers for a previous stage, but trying to delete all")


			log.Error("Failed to delete container on Gadget")
			log.Warn("Was the container ever deployed?")
			
			deleteFailed = true
		}
		
	}
	
	if deleteFailed {
		err = errors.New("Failed to delete one or more containers")
	}
	
	return err
}
