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
func GadgetRmi(container libgadget.GadgetContainer, g *libgadget.GadgetContext) error {

	libgadget.EnsureKeys()

	client, err := libgadget.GadgetLogin(libgadget.GadgetPrivKeyLocation)

	if err != nil {
		return err
	}

	log.Infof("Removing image:")

	log.Infof("  %s", container.ImageAlias)

	stdout, stderr, err := libgadget.RunRemoteCommand(client, "docker", "rmi", container.ImageAlias)

	log.WithFields(log.Fields{
		"function":     "GadgetRmi",
		"name":         container.Alias,
		"delete-stage": "rmi",
	}).Debug(stdout)

	log.WithFields(log.Fields{
		"function":     "GadgetRmi",
		"name":         container.Alias,
		"delete-stage": "rmi",
	}).Debug(stderr)

	if err != nil {

		log.WithFields(log.Fields{
			"function":     "GadgetRmi",
			"name":         container.Alias,
			"delete-stage": "rmi",
		}).Debug("This is likely due to specifying containers for a previous stage, but trying to remove all")

		log.Error("Failed to remove image on Gadget")
		log.Warn("Was the image ever deployed?")

		return err
	}

	return nil
}

// Process the build arguments and execute build
func GadgetRmiDanglers(g *libgadget.GadgetContext) error {

	libgadget.EnsureKeys()

	client, err := libgadget.GadgetLogin(libgadget.GadgetPrivKeyLocation)

	if err != nil {
		return err
	}

	log.Debug("Removing danglers:")

	// not checking for error, as it's bound to fail when there are no dangles
	stdout, stderr, _ := libgadget.RunRemoteCommand(client, "docker", "rmi", `$(docker images -q --filter "dangling=true")`)

	log.WithFields(log.Fields{
		"function":     "GadgetRmiDangle",
		"delete-stage": "rmi",
	}).Debug(stdout)
	log.WithFields(log.Fields{
		"function":     "GadgetRmiDangle",
		"delete-stage": "rmi",
	}).Debug(stderr)

	return err
}

// Process the build arguments and execute build
func GadgetPurge(garbage []string, g *libgadget.GadgetContext) error {

	libgadget.EnsureKeys()

	client, err := libgadget.GadgetLogin(libgadget.GadgetPrivKeyLocation)

	if err != nil {
		return err
	}

	// not checking for error, as it's bound to fail when there are no dangles

	log.Info("Removing containers ..")
	stdout, stderr, _ := libgadget.RunRemoteCommand(client, "docker", "rm", "-f", `$(docker ps -aq)`)

	log.WithFields(log.Fields{
		"function":     "GadgetPurge",
		"delete-stage": "rm -f",
	}).Debug(stdout)
	log.WithFields(log.Fields{
		"function":     "GadgetPurge",
		"delete-stage": "rm -f",
	}).Debug(stderr)

	log.Info("Removing images..")
	stdout, stderr, _ = libgadget.RunRemoteCommand(client, "docker", "rmi", "-f", `$(docker images -aq)`)

	log.WithFields(log.Fields{
		"function":     "GadgetPurge",
		"delete-stage": "rmi -f",
	}).Debug(stdout)
	log.WithFields(log.Fields{
		"function":     "GadgetPurge",
		"delete-stage": "rmi -f",
	}).Debug(stderr)

	log.Info("Removing volumes..")
	stdout, stderr, _ = libgadget.RunRemoteCommand(client, "docker", "volume", "rm", `$(docker volume ls -q)`)

	log.WithFields(log.Fields{
		"function":     "GadgetPurge",
		"delete-stage": "rmi -f",
	}).Debug(stdout)
	log.WithFields(log.Fields{
		"function":     "GadgetPurge",
		"delete-stage": "rmi -f",
	}).Debug(stderr)

	return err
}

// Process the build arguments and execute build
func GadgetDelete(args []string, g *libgadget.GadgetContext) error {

	libgadget.EnsureKeys()

	log.Info("Deleting:")

	stagedContainers, _ := libgadget.FindStagedContainers(args, append(g.Config.Onboot, g.Config.Services...))

	deleteFailed := false

	for _, container := range stagedContainers {
		log.Infof("  %s", container.ImageAlias)

		storedLevel := log.GetLevel()
		if !g.Verbose {
			log.SetLevel(log.ErrorLevel)
		}

		// stop container
		log.Debug("Stopping container: %s", container.Alias)
		toStop := []string{container.Name}
		err := GadgetStop(toStop, g)
		if err != nil {
			deleteFailed = true
		}

		// remove container
		log.Debug("Removing container: %s", container.Alias)
		err = GadgetRm(container, g)
		if err != nil {
			deleteFailed = true
		}

		// delete image
		log.Debug("Deleting image: %s", container.ImageAlias)
		err = GadgetRmi(container, g)
		if err != nil {
			deleteFailed = true
		}

		// delete image danglers
		log.Debug("Deleting image danglers: %s", container.ImageAlias)
		_ = GadgetRmiDanglers(g)

		if !g.Verbose {
			log.SetLevel(storedLevel)
		}

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

	var err error = nil
	if deleteFailed {
		err = errors.New("Failed to delete one or more containers")
	}

	// copy the config file over for autostarts
	libgadget.GadgetInstallConfig(g)

	return err
}
