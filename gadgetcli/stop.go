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
	"github.com/nextthingco/libgadget"
	log "gopkg.in/sirupsen/logrus.v1"
)

// Process the build arguments and execute build
func GadgetRm(args []string, g *libgadget.GadgetContext) error {

	libgadget.EnsureKeys()

	client, err := libgadget.GadgetLogin(libgadget.GadgetPrivKeyLocation)

	if err != nil {
		return err
	}

	log.Info("Removing:")
	stagedContainers, _ := libgadget.FindStagedContainers(args, append(g.Config.Onboot, g.Config.Services...))

	var rmFailed bool = false

	for _, container := range stagedContainers {
		log.Infof("  %s", container.Alias)

		stdout, stderr, err := libgadget.RunRemoteCommand(client, "docker rm", container.Alias)

		log.WithFields(log.Fields{
		"function": "GadgetRm",
		"name": container.Alias,
		"stop-stage": "rm",
		}).Debug(stdout)
		log.WithFields(log.Fields{
		"function": "GadgetRm",
		"name": container.Alias,
		"stop-stage": "rm",
		}).Debug(stderr)

		if err != nil {

		rmFailed = true

		log.WithFields(log.Fields{
		"function": "GadgetRm",
		"name": container.Alias,
		"stop-stage": "rm",
		}).Debug("This is likely due to specifying containers for a previous operation, but trying to stop all")

		log.Errorf("Failed to stop '%s' on Gadget", container.Name)
		log.Warn("Was it ever started?")

		} else {
		log.Info("  - stopped")
		}
	}

	if rmFailed {
		err = errors.New("A problem was encountered in GadgetStop")
	}

	return err
}

// Process the build arguments and execute build
func GadgetStop(args []string, g *libgadget.GadgetContext) error {

	libgadget.EnsureKeys()

	client, err := libgadget.GadgetLogin(libgadget.GadgetPrivKeyLocation)

	if err != nil {
		return err
	}

	log.Info("Stopping:")
	stagedContainers, _ := libgadget.FindStagedContainers(args, append(g.Config.Onboot, g.Config.Services...))

	var stopFailed bool = false

	for _, container := range stagedContainers {
		log.Infof("  %s", container.Alias)

		stdout, stderr, err := libgadget.RunRemoteCommand(client, "docker stop", container.Alias)

		log.WithFields(log.Fields{
			"function":   "GadgetStart",
			"name":       container.Alias,
			"stop-stage": "stop",
		}).Debug(stdout)
		log.WithFields(log.Fields{
			"function":   "GadgetStart",
			"name":       container.Alias,
			"stop-stage": "stop",
		}).Debug(stderr)

		if err != nil {

			stopFailed = true

			log.WithFields(log.Fields{
				"function":   "GadgetStop",
				"name":       container.Alias,
				"stop-stage": "stop",
			}).Debug("This is likely due to specifying containers for a previous operation, but trying to stop all")

			log.Debug("Failed to stop container on Gadget,")
			log.Debug("it might have never been deployed,")
			log.Debug("Or stop otherwise failed")

		}
	}

	if stopFailed {
		err = errors.New("A problem was encountered in GadgetStop")
	}

	return err
}
