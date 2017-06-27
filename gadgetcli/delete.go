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

	stagedContainers, _ := libgadget.FindStagedContainers(args, append(g.Config.Onboot, g.Config.Services...))

	deleteFailed := false

	for _, container := range stagedContainers {
		log.Infof("  %s", container.ImageAlias)

		stdout, stderr, err := libgadget.RunRemoteCommand(client, "docker rm", container.Alias)

		log.WithFields(log.Fields{
			"function":   "GadgetStart",
			"name":       container.Alias,
			"stop-stage": "rm",
		}).Debug(stdout)
		log.WithFields(log.Fields{
			"function":   "GadgetStart",
			"name":       container.Alias,
			"stop-stage": "rm",
		}).Debug(stderr)

		//~ if err != nil {

		//~ stopFailed = true

		//~ log.WithFields(log.Fields{
		//~ "function": "GadgetStop",
		//~ "name": container.Alias,
		//~ "stop-stage": "rm",
		//~ }).Debug("This is likely due to specifying containers for a previous operation, but trying to stop all")

		//~ log.Errorf("Failed to stop '%s' on Gadget", container.Name)
		//~ log.Warn("Was it ever started?")

		//~ } else {
		//~ log.Info("  - stopped")
		//~ }

		stdout, stderr, err = libgadget.RunRemoteCommand(client, "docker", "rmi", container.ImageAlias)

		log.WithFields(log.Fields{
			"function":     "GadgetDelete",
			"name":         container.Alias,
			"delete-stage": "rmi",
		}).Debug(stdout)
		log.WithFields(log.Fields{
			"function":     "GadgetDelete",
			"name":         container.Alias,
			"delete-stage": "rmi",
		}).Debug(stderr)

		if err != nil {

			log.WithFields(log.Fields{
				"function":     "GadgetDelete",
				"name":         container.Alias,
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

	// copy the config file over for autostarts
	libgadget.GadgetInstallConfig(g)

	return err
}
