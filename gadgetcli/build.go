/*
This file is part of the Gadget command-line tools.
Copyright (C) 2017 Next Thing Co.

Gadget is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 2 of the License, or
(at your option) any later version.

Gadget is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with Gadget.  If not, see <http://www.gnu.org/licenses/>.
*/

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
	
	err = libgadget.EnsureDocker(binary, g)
	if err != nil {
		log.Errorf("Failed to contact the docker daemon.")
		log.Warnf("Is it installed and running with appropriate permissions?")
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

				log.WithFields(log.Fields{
					"function": "GadgetBuild",
					"name":     container.Alias,
				}).Debug("The build command returned an error, possible sources are any docker failure scenario")

			} else {
				log.Info("    Done ✔")
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
				log.Warn("Are you sure '%s' is a valid image [and tag]?")

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

				log.WithFields(log.Fields{
					"function": "GadgetBuild",
					"name":     container.Alias,
				}).Debug("The build command returned an error, possible sources are any docker failure scenario")

			} else {
				log.Info("    Done ✔")
			}
		}

	}

	if buildFailed {
		err = errors.New("Failed to build one or more containers")
	}

	return err
}
