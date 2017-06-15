package main

import (
	"errors"
	"fmt"
	"github.com/nextthingco/libgadget"
	log "github.com/sirupsen/logrus"
)

// Process the build arguments and execute build
func GadgetLogs(args []string, g *libgadget.GadgetContext) error {

	libgadget.EnsureKeys()

	client, err := libgadget.GadgetLogin(libgadget.GadgetPrivKeyLocation)

	if err != nil {
		return err
	}

	log.Info("Retrieving logs:")

	logsFailed := false

	stagedContainers, _ := libgadget.FindStagedContainers(args, append(g.Config.Onboot, g.Config.Services...))
	for _, container := range stagedContainers {
		commandFormat := `docker logs %s`
		cmd := fmt.Sprintf(commandFormat, container.Alias)

		stdout, stderr, err := libgadget.RunRemoteCommand(client, cmd)

		if err != nil {

			// fail loudly, but continue

			logsFailed = true

			log.Errorf("Failed to fetch '%s' logs on Gadget", container.Name)
			log.Warn("Was the container ever deployed?")

			log.WithFields(log.Fields{
				"function":    "GadgetLogs",
				"name":        container.Alias,
				"start-stage": "docker logs",
			}).Debug("This is likely due to specifying containers for deploying, but trying to fetch all logs")

		} else {

			log.Infof("  Begin: %s\n", container.Name)
			log.Infof("\n%s", stdout)
			log.Warnf("\n%s", stderr)
			log.Infof("  End: %s", container.Name)

		}
	}

	if logsFailed {
		err = errors.New("Failed to fetch logs for one or more containers")
	}

	return err
}
