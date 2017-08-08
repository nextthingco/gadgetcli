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
	log "gopkg.in/sirupsen/logrus.v1"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"runtime"
)

func GadgetBuildRootfs(g *libgadget.GadgetContext) error {
	log.Infof("  '%s' rootfs", g.Config.Rootfs.From)

	// find docker binary in path
	binary, err := exec.LookPath("docker")
	if err != nil {
		log.Error("Failed to find local docker binary")
		log.Warn("Is docker installed?")

		log.WithFields(log.Fields{
			"function": "GadgetAddRootfs",
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

	image := g.Config.Rootfs.Hash
	board := g.Config.Rootfs.From

	linuxConfig := fmt.Sprintf("%s/%s-linux.config", g.WorkingDirectory, board)
	configExists, err := libgadget.PathExists(linuxConfig)
	if !configExists {
		log.Errorf("Could not locate '%s'", linuxConfig)
		return errors.New("Failed to locate linux config")
	}
	if err != nil {
		log.Errorf("Failed to determine if '%s' exists", linuxConfig)
		return err
	}

	// check/create ./.build/$board/output
	imagesDir := filepath.Join(g.WorkingDirectory, "/.images/")
	imagesDirExists, err := libgadget.PathExists(imagesDir)
	if err != nil {
		log.WithFields(log.Fields{
			"function": "GadgetBuildRootfs",
			"error":    err,
		}).Error("Couldn't determine if the ./.build/$board/images directory exists.")
		return err
	}

	if !imagesDirExists {
		err = os.Mkdir(imagesDir, 0755)
		if err != nil {
			log.WithFields(log.Fields{
				"function": "GadgetBuildRootfs",
				"error":    err,
			}).Error("Couldn't create ./.build/$board/images directory")
			return err
		}
	}

	linuxConfigBinds := fmt.Sprintf("%s/%s-linux.config:/opt/gadget-os-proto/gadget/board/nextthing/%s/configs/linux.config", g.WorkingDirectory, board, board)
	imagesBinds := fmt.Sprintf("%s:/opt/output/images", imagesDir)
	cmd := exec.Command("docker", "run", "-it", "--rm", "-e", "BOARD="+board, "-e", "no_docker=1", "-v", imagesBinds, "-v", linuxConfigBinds, image, "make", "gadget_build")

	cmd.Env = os.Environ()

	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr

	if err := cmd.Start(); err != nil {
		log.Errorf("An error occured: ", err)
		return err
	}

	cmd.Wait()

	// chown kernelconfig
	if runtime.GOOS != "windows" {

		whois, err := user.Current()
		if err != nil {
			log.Error("Failed to retrieve UID/GID")
			return err
		}

		chownAs := whois.Uid + ":" + whois.Gid
		imagesBinds := fmt.Sprintf("%s:/chown", imagesDir)
		stdout, stderr, err := libgadget.RunLocalCommand(binary,
			"", g,
			"run", "--rm", "-v", imagesBinds,
			image,
			"/bin/chown", "-R", chownAs, "/chown")

		log.WithFields(log.Fields{
			"function": "GadgetAddRootfs",
			"name":     image,
			"stage":    "docker tag",
		}).Debug(stdout)
		log.WithFields(log.Fields{
			"function": "GadgetAddRootfs",
			"name":     image,
			"stage":    "docker tag",
		}).Debug(stderr)

		if err != nil {
			log.Error("Failed to chown linux config")
			return err
		}

	}

	return nil
}

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
				"Step ", g,
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
				"Download", g,
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
				"", g,
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
		err = errors.New("Failed to build one or more artifacts")
	}

	if len(args) < 1 && g.Config.Rootfs.From != "" && g.Config.Rootfs.Hash != "" {
		err = GadgetBuildRootfs(g)
		if err != nil {
			log.Errorf("Failed to build rootfs")
			return err
		}
	}

	return err
}
