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
	"golang.org/x/crypto/ssh"
	//~ "gopkg.in/cheggaaa/pb.v1"
	log "gopkg.in/sirupsen/logrus.v1"
	//~ "io"
	//~ "os/exec"
	//~ "strings"
)

type ArtifactDef struct {
	Board        string
	Artifacts    []string
	ArtifactType []string
}

var (
	ArtDefs = []ArtifactDef{
		ArtifactDef { 
			Board:        "chippro",
			Artifacts:    []string {"zImage", "ntc-gr8-crumb.dtb", "rootfs.ubifs"},
			ArtifactType: []string {"kernel", "fdt", "rootfs"},
		},
		//~ ArtifactDef { 
			//~ Board: "chippro4gb",
			//~ Artifacts: []string {"zImage", "ntc-gr8-crumb.dtb", "rootfs.ubifs"},
		//~ },
		//~ ArtifactDef {
			//~ Board: "chip",
			//~ Artifacts: []string {"zImage", "ntc-r8-chip.dtb", "rootfs.ubifs"},
		//~ },
	}
)

func GadgetFlashFile(client *ssh.Client, artifactLocation string, artifactType string, g *libgadget.GadgetContext) error {
	_ = client
	
	log.Infof("artLoc: %s", artifactLocation)
	log.Infof("artTyp: %s", artifactType)
	
	//~ binary, err := exec.LookPath("docker")
	//~ if err != nil {
		//~ return err
	//~ }

	//~ log.Infof("Deploying: '%s'", container.Name)
	//~ docker := exec.Command(binary, "save", container.ImageAlias)

	//~ session, err := client.NewSession()
	//~ if err != nil {
		//~ client.Close()
		//~ return err
	//~ }

	//~ // create pipe for local -> remote file transmission
	//~ pr, pw := io.Pipe()
	//~ sessionLogger := log.New()
	//~ if g.Verbose {
		//~ sessionLogger.Level = log.DebugLevel
	//~ }

	//~ bar := pb.New(0)
	//~ bar.SetUnits(pb.U_BYTES)
	//~ bar.ShowSpeed = true
	//~ bar.ShowPercent = false
	//~ bar.ShowTimeLeft = false
	//~ bar.ShowBar = false

	//~ docker.Stdout = pw
	//~ reader := bar.NewProxyReader(pr)
	//~ session.Stdin = reader
	//~ session.Stdout = sessionLogger.WriterLevel(log.DebugLevel)
	//~ session.Stderr = sessionLogger.WriterLevel(log.DebugLevel)

	//~ log.Debug("  Starting session")
	//~ if err := session.Start(`docker load`); err != nil {
		//~ return err
	//~ }

	//~ log.Debug("  Starting docker")
	//~ if err := docker.Start(); err != nil {
		//~ return err
	//~ }

	//~ deployFailed := false

	//~ go func() error {
		//~ defer pw.Close()
		//~ log.Info("  Starting transfer..")
		//~ log.Debug("  Waiting on docker")
		//~ bar.Start()

		//~ if err := docker.Wait(); err != nil {
			//~ deployFailed = true
			//~ // TODO: we should handle this error or report to the log
			//~ log.Errorf("Failed to transfer '%s'", container.Name)
			//~ log.Warn("Was the container ever built?")
			//~ return err
		//~ }
		//~ return err
	//~ }()

	//~ session.Wait()
	//~ bar.Finish()
	//~ if !deployFailed {
		//~ log.Info("Done!")
		//~ log.Debug("Closing session")
	//~ }
	//~ session.Close()

	//~ restart := ""
	//~ mode, err := FindRunMode(container.UUID, g.Config.Onboot, g.Config.Services)
	//~ if err != nil {
		//~ log.Debug(err)
		//~ log.Errorf("Failed to find run mode for '%s:%s'", container.Name, container.UUID)
		//~ return err
	//~ } else if mode == GADGETSERVICE {
		//~ restart = "--restart=on-failure"
	//~ }

	//~ tmpString := []string{container.Net}
	//~ net := strings.Join(libgadget.PrependToStrings(tmpString[:], "--net "), " ")

	//~ tmpString = []string{container.PID}
	//~ pid := strings.Join(libgadget.PrependToStrings(tmpString[:], "--pid "), " ")

	//~ readOnly := ""
	//~ if container.Readonly {
		//~ readOnly = "--read-only"
	//~ }

	//~ binds := strings.Join(libgadget.PrependToStrings(container.Binds[:], "-v "), " ")
	//~ caps := strings.Join(libgadget.PrependToStrings(container.Capabilities[:], "--cap-add "), " ")
	//~ devs := strings.Join(libgadget.PrependToStrings(container.Devices[:], "--device "), " ")
	//~ commands := strings.Join(container.Command[:], " ")

	//~ stdout, stderr, err := libgadget.RunRemoteCommand(client, "docker create --name", container.Alias,
		//~ net, pid, readOnly, binds, caps, devs, restart, container.ImageAlias, commands)

	//~ log.Debugf("docker create --name %s %s %s %s %s %s %s %s", container.Alias,
		//~ net, pid, readOnly, binds, caps, devs, restart, container.ImageAlias, commands)

	//~ // delete image danglers
	//~ err = GadgetRmiDanglers(g)

	//~ log.WithFields(log.Fields{
		//~ "function":     "GadgetDelete",
		//~ "name":         container.Alias,
		//~ "delete-stage": "rmi (danglers)",
	//~ }).Debug(stdout)
	//~ log.WithFields(log.Fields{
		//~ "function":     "GadgetDelete",
		//~ "name":         container.Alias,
		//~ "delete-stage": "rmi (danglers)",
	//~ }).Debug(stderr)

	//~ if err != nil {

		//~ log.Errorf("Failed to set %s to always restart on Gadget", container.Alias)
		//~ return err
	//~ }

	//~ log.WithFields(log.Fields{
		//~ "function":     "DeployContainer",
		//~ "name":         container.Alias,
		//~ "deploy-stage": "create restarting",
	//~ }).Debug(stdout)
	//~ log.WithFields(log.Fields{
		//~ "function":     "DeployContainer",
		//~ "name":         container.Alias,
		//~ "deploy-stage": "create restarting",
	//~ }).Debug(stderr)

	//~ // copy the config file over for autostarts
	//~ libgadget.GadgetInstallConfig(g)

	return nil
}

// Process the build arguments and execute build
func GadgetFlash(args []string, g *libgadget.GadgetContext) error {

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
	
	// check for non-empty board definition
	board := g.Config.Rootfs.From
	image := g.Config.Rootfs.Hash
	
	if board == "" && image == "" {
		log.Errorf("Failed to find rootfs")
		log.Errorf("One or more [mis/un]configured entries:")
		log.Errorf("From: %s", board)
		log.Errorf("Hash: %s", image)
		return errors.New("Failed to flash rootfs")
	}
	
	// check that board is supported
	matchedBoard := ArtifactDef { Board: "", }
	
	for _, def := range ArtDefs {
		if board == def.Board {
			matchedBoard = def
			log.Infof("  Flashing: %s", board)
			break
		}
	}
	
	if matchedBoard.Board == "" {
		log.Errorf("%s is not a valid From:", board)
		return errors.New("Invalid board definition")
	}
	
	// test to make sure all payload files present
	for _, payloadPart := range matchedBoard.Artifacts {
		partLocation := g.WorkingDirectory + "/.images/" + payloadPart
		partExists, err := libgadget.PathExists(partLocation)
		if !partExists {
			log.Errorf("Could not locate '%s'", partLocation)
			return errors.New("Failed to locate linux config")
		}
		if err != nil {
			log.Errorf("Failed to determine if '%s' exists", partLocation)
			return err
		}
	}
	
	// flash each part
	for i, flashPart := range matchedBoard.Artifacts {
		
		partLocation := g.WorkingDirectory + "/.images/" + flashPart
		partType := matchedBoard.ArtifactType[i]
		
		err = GadgetFlashFile(client, partLocation, partType, g)
		if err != nil {
			return err
		}
	}
	
	return err
}
