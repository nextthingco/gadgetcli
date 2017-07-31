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
	"os/exec"
	"os"
)

var (
)

func editUsage() error {
	log.Info("Usage:  gadget [flags] edit [type] [value]     ")
	log.Info("                *opt         *req   *opt       ")
	log.Info("Type:           service | onboot | rootfs      ")
	log.Info("Value (containers): not yet implemented        ")
	log.Info("Value (rootfs): kernel <more to be added soon> ")

	return errors.New("Incorrect edit usage")
}

func GadgetEditKernel(g *libgadget.GadgetContext) error {
	
	cmd := exec.Command("docker", "run", "-it", "--rm", g.Config.Rootfs.Hash, "make", "linux-menuconfig")

	cmd.Env = os.Environ()

	cmd.Stdin , cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	
    if err := cmd.Start(); err != nil {
        log.Errorf("An error occured: ", err)
        return err
    }
    
    cmd.Wait()	
	
	return nil
}

func GadgetEditUserspace(g *libgadget.GadgetContext) error {
	
	cmd := exec.Command("docker", "run", "-it", "--rm", g.Config.Rootfs.Hash, "make", "menuconfig")

	cmd.Env = os.Environ()

	cmd.Stdin , cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	
    if err := cmd.Start(); err != nil {
        log.Errorf("An error occured: ", err)
        return err
    }
    
    cmd.Wait()	
	
	return nil
}

func GadgetEditUboot(g *libgadget.GadgetContext) error {
	
	cmd := exec.Command("docker", "run", "-it", "--rm", g.Config.Rootfs.Hash, "make", "uboot-menuconfig")

	cmd.Env = os.Environ()

	cmd.Stdin , cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	
    if err := cmd.Start(); err != nil {
        log.Errorf("An error occured: ", err)
        return err
    }
    
    cmd.Wait()	
	
	return nil
}

// Process the build arguments and execute build
func GadgetEdit(args []string, g *libgadget.GadgetContext) error {
	
	log.Info("Edit")
	log.Debugf("args %s", args)
	
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
	
	if len(args) != 2 {
		log.Error("Invalid arguments for `gadget edit`")
		return editUsage();
	}
	
	// parse arguments
	switch args[0] {
	case "rootfs":
		// parse edit rootfs argument
		switch args[1] {
		case "kernel":
			err = GadgetEditKernel(g)
			if err != nil {
				log.Errorf("Failed to edit the kernel config.")
				return err
			}
		case "userspace":
			err = GadgetEditUserspace(g)
			if err != nil {
				log.Errorf("Failed to edit the userspace config.")
				return err
			}
		case "uboot":
			err = GadgetEditUboot(g)
			if err != nil {
				log.Errorf("Failed to edit the uboot config.")
				return err
			}
		default:
			log.Errorf("  %q is not valid argument or is not yet supported.", args[1])
			return editUsage()
		}
	default:
		log.Errorf("  %q is not valid argument or is not yet supported.", args[0])
		return editUsage()
	}

	return nil
}
