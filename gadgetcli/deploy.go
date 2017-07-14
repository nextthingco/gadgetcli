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
	"golang.org/x/crypto/ssh"
	"gopkg.in/cheggaaa/pb.v1"
	"io"
	"os/exec"
	"strings"
)

func DeployContainer(client *ssh.Client, container *libgadget.GadgetContainer, g *libgadget.GadgetContext) error {

	binary, err := exec.LookPath("docker")
	if err != nil {
		return err
	}

	log.Infof("Deploying: '%s'", container.Name)
	docker := exec.Command(binary, "save", container.ImageAlias)

	session, err := client.NewSession()
	if err != nil {
		client.Close()
		return err
	}

	// create pipe for local -> remote file transmission
	pr, pw := io.Pipe()
	sessionLogger := log.New()
	if g.Verbose {
		sessionLogger.Level = log.DebugLevel
	}

	bar := pb.New(0)
	bar.SetUnits(pb.U_BYTES)
	bar.ShowSpeed = true
	bar.ShowPercent = false
	bar.ShowTimeLeft = false
	bar.ShowBar = false

	docker.Stdout = pw
	reader := bar.NewProxyReader(pr)
	session.Stdin = reader
	session.Stdout = sessionLogger.WriterLevel(log.DebugLevel)
	session.Stderr = sessionLogger.WriterLevel(log.DebugLevel)

	log.Debug("  Starting session")
	if err := session.Start(`docker load`); err != nil {
		return err
	}

	log.Debug("  Starting docker")
	if err := docker.Start(); err != nil {
		return err
	}

	deployFailed := false

	go func() error {
		defer pw.Close()
		log.Info("  Starting transfer..")
		log.Debug("  Waiting on docker")
		bar.Start()

		if err := docker.Wait(); err != nil {
			deployFailed = true
			// TODO: we should handle this error or report to the log
			log.Errorf("Failed to transfer '%s'", container.Name)
			log.Warn("Was the container ever built?")
			return err
		}
		return err
	}()

	session.Wait()
	bar.Finish()
	if !deployFailed {
		log.Info("Done!")
		log.Debug("Closing session")
	}
	session.Close()
	
	restart := ""
	mode, err := FindRunMode(container.UUID, g.Config.Onboot, g.Config.Services)
	if err != nil {
		log.Debug(err)
		log.Errorf("Failed to find run mode for '%s:%s'", container.Name, container.UUID)
		return err
	} else if mode == GADGETSERVICE {
		restart = "--restart=on-failure"
	}

	tmpString := []string{container.Net}
	net := strings.Join(libgadget.PrependToStrings(tmpString[:], "--net "), " ")

	tmpString = []string{container.PID}
	pid := strings.Join(libgadget.PrependToStrings(tmpString[:], "--pid "), " ")

	readOnly := ""
	if container.Readonly {
		readOnly = "--read-only"
	}

	binds := strings.Join(libgadget.PrependToStrings(container.Binds[:], "-v "), " ")
	caps := strings.Join(libgadget.PrependToStrings(container.Capabilities[:], "--cap-add "), " ")
	devs := strings.Join(libgadget.PrependToStrings(container.Devices[:], "--device "), " ")
	commands := strings.Join(container.Command[:], " ")

	stdout, stderr, err := libgadget.RunRemoteCommand(client, "docker create --name", container.Alias,
		net, pid, readOnly, binds, caps, devs, restart, container.ImageAlias, commands)

	log.Debugf("docker create --name %s %s %s %s %s %s %s %s", container.Alias,
		net, pid, readOnly, binds, caps, devs, restart, container.ImageAlias, commands)
	
	
	// delete image danglers
	err = GadgetRmiDanglers( g)

	log.WithFields(log.Fields{
		"function":     "GadgetDelete",
		"name":         container.Alias,
		"delete-stage": "rmi (danglers)",
	}).Debug(stdout)
	log.WithFields(log.Fields{
		"function":     "GadgetDelete",
		"name":         container.Alias,
		"delete-stage": "rmi (danglers)",
	}).Debug(stderr)
	
	if err != nil {

		log.Errorf("Failed to set %s to always restart on Gadget", container.Alias)
		return err
	}

	log.WithFields(log.Fields{
		"function":     "DeployContainer",
		"name":         container.Alias,
		"deploy-stage": "create restarting",
	}).Debug(stdout)
	log.WithFields(log.Fields{
		"function":     "DeployContainer",
		"name":         container.Alias,
		"deploy-stage": "create restarting",
	}).Debug(stderr)

	// copy the config file over for autostarts
	libgadget.GadgetInstallConfig(g)

	return err
}

const (
	GADGETONBOOT = 1 << iota
	GADGETSERVICE
)

func FindRunMode(uuid string, onboot []libgadget.GadgetContainer, services []libgadget.GadgetContainer) (int, error) {
	for _, container := range onboot {
		if container.UUID == uuid {
			return GADGETONBOOT, nil
		}
	}
	for _, container := range services {
		if container.UUID == uuid {
			return GADGETSERVICE, nil
		}
	}

	return 0, errors.New("Failed to find by UUID")
}

// Process the build arguments and execute build
func GadgetDeploy(args []string, g *libgadget.GadgetContext) error {

	err := libgadget.EnsureKeys()
	if err != nil {
		log.Errorf("Failed to connect to Gadget")
		return err
	}

	client, err := libgadget.GadgetLogin(libgadget.GadgetPrivKeyLocation)
	if err != nil {
		log.Errorf("Failed to connect to Gadget")
		return err
	}

	stagedContainers, err := libgadget.FindStagedContainers(args, append(g.Config.Onboot, g.Config.Services...))

	deployFailed := false

	for _, container := range stagedContainers {

		// stop and delete possible older versions of image/container
		// not collecting the errors, as errors may be returned
		// when trying to delete an img/cntnr that was never deployed
		tmpName := make([]string, 1)
		tmpName[0] = container.Name

		log.Infof("Stopping/deleting older '%s' if applicable", container.Name)

		if !g.Verbose {
			log.SetLevel(log.PanicLevel)
		}

		_ = GadgetStop(tmpName, g)
		_ = GadgetRm(tmpName, g)

		if !g.Verbose {
			log.SetLevel(log.InfoLevel)
		}

		err = DeployContainer(client, &container, g)
		deployFailed = true
	}

	if deployFailed == true {
		err = errors.New("Failed to deploy one or more containers")
	}

	return err
}
