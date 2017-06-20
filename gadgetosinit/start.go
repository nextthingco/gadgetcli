package main

import (
	"errors"
	"github.com/nextthingco/libgadget"
	log "github.com/sirupsen/logrus"
	"os/exec"
	"strings"
)

// Process the build arguments and execute build
func GadgetOsInit(args []string, g *libgadget.GadgetContext) error {

	binary, err := exec.LookPath("docker")
	if err != nil {
		log.Error("Failed to find local docker binary")
		log.Warn("Is docker installed?")

		log.WithFields(log.Fields{
			"function": "GadgetOsInit",
			"stage":    "LookPath(docker)",
		}).Debug("Couldn't find docker in the $PATH")
		return err
	}

	var initFailed bool = false

	log.Info("Starting:")

	stagedContainers, _ := libgadget.FindStagedContainers(args, append(g.Config.Onboot, g.Config.Services...))
	for _, container := range stagedContainers {

		log.Infof("  %s", container.Alias)

		commands := strings.Join(container.Command[:], " ")

		stdout, stderr, err := libgadget.RunLocalCommand(binary, g, "start", container.Alias)

		log.WithFields(log.Fields{
			"function":    "GadgetStart",
			"name":        container.Alias,
			"start-stage": "create",
		}).Debug(stdout)
		log.WithFields(log.Fields{
			"function":    "GadgetStart",
			"name":        container.Alias,
			"start-stage": "create",
		}).Debug(stderr)

		if err != nil {
			// fail loudly, but continue

			log.WithFields(log.Fields{
				"function":    "GadgetStart",
				"name":        container.Alias,
				"start-stage": "create",
			}).Debug("This is likely due to specifying containers for deploying, but trying to start all")

			log.Errorf("  Failed to start '%s' on Gadget", container.Name)
			log.Warn("  Potential causes:")
			log.Warn("  - container was never deployed")
			if commands != "" {
				log.Warn("  - conflicting CMD/ENTRYPOINT")
				log.Warnf("    ['%s' was also supplied with the commands '%s']", container.Name, commands)
				log.Warn("    [consult the original Dockerfile to rule out conflicting CMD/ENTRYPOINT]")
			}

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
