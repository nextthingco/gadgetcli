/*
This file is part of the Gadget command-line tools.
Copyright (C) 2017 Robert Wolterman.

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
	log "gopkg.in/sirupsen/logrus.v1"
	"os/exec"
)

func tagUsage() error {
	log.Info("Usage:  gadget [flags] tag [repository] [tag]     ")
	log.Info("                *opt           *req      *opt     ")
	log.Info("Value (repository): repository for        ")
	log.Info("Value (rootfs): kernel <more to be added soon> ")

	return errors.New("Incorrect edit usage")
}

// Process the build arguments and execute build
func GadgetTag(args []string, g *libgadget.GadgetContext) error {

	if len(args) < 1 {
		return tagUsage()
	}

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

	log.Info("Tagging:")

	// We're going to tag all the containers in the config
	stagedContainers := append(g.Config.Onboot, g.Config.Services...)

	tagFailed := false

	for _, container := range stagedContainers {
        // Figure out the tag name
        taggedImage := fmt.Sprintf("%s/%s", args[0], container.Name)
        if len(args) > 1 {
            taggedImage = fmt.Sprintf("%s/%s:%s", args[0], container.Name, args[1])
        }

        log.Infof("  '%s' ➡ '%s'", container.ImageAlias, taggedImage)

		stdout, stderr, err := libgadget.RunLocalCommand(binary,
			"", g,
			"tag",
			container.ImageAlias,
			taggedImage)

		log.WithFields(log.Fields{
			"function": "GadgetTag",
			"name":     container.Alias,
			"stage":    "docker tag",
		}).Debug(stdout)
		log.WithFields(log.Fields{
			"function": "GadgetTag",
			"name":     container.Alias,
			"stage":    "docker tag",
		}).Debug(stderr)

		if err != nil {

			tagFailed = true

			log.Errorf("Failed to tag '%s'", container.Name)

			log.WithFields(log.Fields{
				"function": "GadgetTag",
				"name":     container.Name,
			}).Debug("The tag command returned an error, possible sources are any docker failure scenario")

		} else {
			log.Info("    Done ✔")
		}
	}

	if tagFailed {
		err = errors.New("Failed to tag one or more containers")
	}

	return err
}
