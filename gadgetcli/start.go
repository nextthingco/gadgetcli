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
	"strings"
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
	stagedContainers, _ := libgadget.FindStagedContainers(args, append(g.Config.Onboot, g.Config.Services...))
	for _, container := range stagedContainers {

		log.Infof("  %s", container.Alias)
		commands := strings.Join(container.Command[:], " ")

		stdout, stderr, err := libgadget.RunRemoteCommand(client, "docker start", container.Alias)

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
