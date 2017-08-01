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
	"gopkg.in/satori/go.uuid.v1"
	log "gopkg.in/sirupsen/logrus.v1"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os/exec"
)

var (
	defconfigList = []string{"chippro", "chippro4gb", "chip"}
)

func addUsage() error {
	log.Info("Usage:  gadget [flags] add [type] [name]    ")
	log.Info("                *opt        *req   *req     ")
	log.Info("Type:           service | onboot | rootfs   ")
	log.Info("Name (service): friendly name for container ")
	log.Info("Name (rootfs):  valid defconfig prefix      ")

	return errors.New("Incorrect add usage")
}

func GadgetAddRootfs(board string, g *libgadget.GadgetContext) (error, string) {

	if g.Config.Rootfs.From != "" ||
		g.Config.Rootfs.Hash != "" {

		log.Error("Rootfs entry already exists")
		log.Infof("Clear values: 'from: %s' and 'hash: %s' to update the image base", g.Config.Rootfs.From, g.Config.Rootfs.Hash)
		return errors.New("Rootfs setting already exists"), ""
	}

	log.Infof("Retrieving build image for '%s'", board)

	// find docker binary in path
	binary, err := exec.LookPath("docker")
	if err != nil {
		log.Error("Failed to find local docker binary")
		log.Warn("Is docker installed?")

		log.WithFields(log.Fields{
			"function": "GadgetAddRootfs",
			"stage":    "LookPath(docker)",
		}).Debug("Couldn't find docker in the $PATH")
		return err, ""
	}

	err = libgadget.EnsureDocker(binary, g)
	if err != nil {
		log.Errorf("Failed to contact the docker daemon.")
		log.Warnf("Is it installed and running with appropriate permissions?")
		return err, ""
	}

	latestContainer := fmt.Sprintf("nextthingco/gadget-build-%s-%s:latest", board, libgadget.GitBranch)

	// pull container
	stdout, stderr, err := libgadget.RunLocalCommand(binary,
		"Downloading", g,
		"pull",
		latestContainer)

	log.WithFields(log.Fields{
		"function": "GadgetAddRootfs",
		"name":     latestContainer,
		"stage":    "docker pull",
	}).Debug(stdout)
	log.WithFields(log.Fields{
		"function": "GadgetAddRootfs",
		"name":     latestContainer,
		"stage":    "docker pull",
	}).Debug(stderr)

	if err != nil {
		log.Error("Failed to download build image")
		return err, ""
	}

	bashLine := fmt.Sprintf("%s", fmt.Sprintf("/bin/echo nextthingco/gadget-build-%s-$(cat .branch)", board))

	// get hash for container
	hash, stderr, err := libgadget.RunLocalCommand(binary,
		"", g,
		"run",
		"-i",
		"--rm",
		latestContainer,
		"/bin/bash",
		"-c",
		bashLine,
	)

	log.WithFields(log.Fields{
		"function": "GadgetAddRootfs",
		"name":     latestContainer,
		"bashLine": bashLine,
		"stream":   "stdout",
		"stage":    "get tag name",
	}).Debug(hash)
	log.WithFields(log.Fields{
		"function": "GadgetAddRootfs",
		"name":     latestContainer,
		"bashLine": bashLine,
		"stream":   "stderr",
		"stage":    "get tag name",
	}).Debug(stderr)

	if err != nil {
		log.Error("Failed to parse build image information")
		return err, hash
	}

	// tag container
	stdout, stderr, err = libgadget.RunLocalCommand(binary,
		"", g,
		"tag",
		latestContainer,
		hash)

	log.WithFields(log.Fields{
		"function": "GadgetAddRootfs",
		"name":     latestContainer,
		"stage":    "docker tag",
	}).Debug(stdout)
	log.WithFields(log.Fields{
		"function": "GadgetAddRootfs",
		"name":     latestContainer,
		"stage":    "docker tag",
	}).Debug(stderr)

	if err != nil {
		log.Error("Failed to tag build image")
		return err, ""
	}

	log.Debugf("hash: %s", hash)

	return nil, hash
}

// Process the build arguments and execute build
func GadgetAdd(args []string, g *libgadget.GadgetContext) error {

	addUu := uuid.NewV4()

	if len(args) != 2 {
		return addUsage()
	}

	log.Infof("Adding new %s: \"%s\" ", args[0], args[1])

	addGadgetContainer := libgadget.GadgetContainer{
		Name:  args[1],
		Image: fmt.Sprintf("%s/%s", g.Config.Name, args[1]),
		UUID:  fmt.Sprintf("%s", addUu),
	}

	// parse arguments
	switch args[0] {
	case "service":
		g.Config.Services = append(g.Config.Services, addGadgetContainer)
	case "onboot":
		g.Config.Onboot = append(g.Config.Onboot, addGadgetContainer)
	case "rootfs":
		matched := ""
		for _, i := range defconfigList {
			if args[1] == i {
				matched = i
				break
			}
		}
		if matched == "" {
			log.Errorf("  %q is not valid defconfig.", args[1])
			return addUsage()
		}

		garError, garHash := GadgetAddRootfs(matched, g)
		if garError != nil {
			return garError
		}

		g.Config.Rootfs.From = matched
		g.Config.Rootfs.Hash = garHash
	default:
		log.Errorf("  %q is not valid command.", args[0])
		return addUsage()
	}

	g.Config = libgadget.CleanConfig(g.Config)

	fileLocation := fmt.Sprintf("%s/gadget.yml", g.WorkingDirectory)

	outBytes, err := yaml.Marshal(g.Config)
	if err != nil {

		log.WithFields(log.Fields{
			"function":   "GadgetAdd",
			"location":   fileLocation,
			"init-stage": "parsing",
		}).Debug("The config file is probably malformed")

		log.Errorf("Failed to parse config file [%s]", fileLocation)
		log.Warn("Is this a valid gadget.yaml?")
		return err
	}

	err = ioutil.WriteFile(fmt.Sprintf("%s/gadget.yml", g.WorkingDirectory), outBytes, 0644)
	if err != nil {

		log.WithFields(log.Fields{
			"function":   "GadgetAdd",
			"location":   fileLocation,
			"init-stage": "writing file",
		}).Debug("This is likely due to a problem with permissions")

		log.Errorf("Failed to edit config file [%s]", fileLocation)
		log.Warn("Do you have permission to modify this file?")

		return err
	}

	return err
}
