package main

import (
	"errors"
	"fmt"
	"github.com/nextthingco/libgadget"
	log "github.com/sirupsen/logrus"
	"os/exec"
)

// Process the build arguments and execute build
func GadgetBuild(args []string, g *libgadget.GadgetContext) error {

	// find docker binary in path
	binary, err := exec.LookPath("docker")
	if err != nil {
		log.Error("Failed to find local docker binary")
		log.Warn("Is docker installed?")

		log.WithFields(log.Fields{
			"function": "GadgetBuild",
			"stage":    "LookPath(docker)",
		}).Debug("Couldn't find docker in the $PATH")
		return err
	}

	log.Info("Building:")

	stagedContainers, _ := libgadget.FindStagedContainers(args, append(g.Config.Onboot, g.Config.Services...))

	buildFailed := false

	for _, container := range stagedContainers {
		log.Infof("  '%s'", container.Name)

		// use local directory for build
		if container.Directory != "" {
			containerDirectory := fmt.Sprintf("%s/%s", g.WorkingDirectory, container.Directory)
			stdout, stderr, err := libgadget.RunLocalCommand(binary,
				g,
				"build",
				"--tag",
				container.ImageAlias,
				containerDirectory)

			log.WithFields(log.Fields{
				"function": "GadgetBuild",
				"name":     container.Alias,
				"stage":    "docker build",
			}).Debug(stdout)
			log.WithFields(log.Fields{
				"function": "GadgetBuild",
				"name":     container.Alias,
				"stage":    "docker build",
			}).Debug(stderr)

			if err != nil {

				buildFailed = true

				log.Errorf("Failed to build '%s'", container.Name)
				log.Warn("Is the docker daemon installed and running?")

				log.WithFields(log.Fields{
					"function": "GadgetBuild",
					"name":     container.Alias,
				}).Debug("The build command returned an error, possible sources are any docker failure scenario")

			}

		} else {
			stdout, stderr, err := libgadget.RunLocalCommand(binary,
				g,
				"pull",
				container.Image)

			log.WithFields(log.Fields{
				"function": "GadgetBuild",
				"name":     container.Alias,
				"stage":    "docker pull",
			}).Debug(stdout)
			log.WithFields(log.Fields{
				"function": "GadgetBuild",
				"name":     container.Alias,
				"stage":    "docker pull",
			}).Debug(stderr)

			if err != nil {

				buildFailed = true

				log.Errorf("Failed to build '%s'", container.Name)
				log.Warn("Is the docker daemon installed and running?")

				log.WithFields(log.Fields{
					"function": "GadgetBuild",
					"name":     container.Alias,
				}).Debug("The build command returned an error, possible sources are any docker failure scenario")

				continue

			}

			stdout, stderr, err = libgadget.RunLocalCommand(binary,
				g,
				"tag",
				container.Image,
				container.ImageAlias)

			log.WithFields(log.Fields{
				"function": "GadgetBuild",
				"name":     container.Alias,
				"stage":    "docker tag",
			}).Debug(stdout)
			log.WithFields(log.Fields{
				"function": "GadgetBuild",
				"name":     container.Alias,
				"stage":    "docker tag",
			}).Debug(stderr)

			if err != nil {

				buildFailed = true

				log.Errorf("Failed to build '%s'", container.Name)
				log.Warn("Is the docker daemon installed and running?")

				log.WithFields(log.Fields{
					"function": "GadgetBuild",
					"name":     container.Alias,
				}).Debug("The build command returned an error, possible sources are any docker failure scenario")

			}
		}

	}

	if buildFailed {
		err = errors.New("Failed to build one or more containers")
	}

	return err
}
