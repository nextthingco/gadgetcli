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
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
	log "gopkg.in/sirupsen/logrus.v1"
	"os"
	"runtime"
)

func GadgetShellContainer(args []string, g *libgadget.GadgetContext) error {

	stagedContainer, err := libgadget.FindStagedContainers(args, append(g.Config.Onboot, g.Config.Services...))

	log.Infof("Attempting to connect to '%s'", args[0])

	client, err := libgadget.GadgetLogin(libgadget.GadgetPrivKeyLocation)
	if err != nil {
		log.WithFields(log.Fields{
			"function":    "GadgetShell",
			"shell-stage": "gadget login",
		}).Debugf("%v", err)
		return err
	}

	session, err := client.NewSession()
	if err != nil {
		log.WithFields(log.Fields{
			"function":    "GadgetShell",
			"shell-stage": "new session",
		}).Debugf("%v", err)
		return err
	}

	session.Stdout = os.Stdout
	session.Stderr = os.Stderr
	session.Stdin = os.Stdin

	modes := ssh.TerminalModes{
		ssh.ECHO:   1, // enable echoing
		ssh.ECHONL: 1,
	}
	if runtime.GOOS == "windows" {
		modes = ssh.TerminalModes{
			ssh.ECHO:   0, // disable echoing
			ssh.ECHONL: 0,
			ssh.IGNCR:  1,
		}
	}

	if err := session.RequestPty("xterm", 25, 80, modes); err != nil {
		session.Close()
		log.WithFields(log.Fields{
			"function":    "GadgetShell",
			"shell-stage": "request pty",
		}).Debugf("%v", err)
		return err
	}

	if err := session.Start(fmt.Sprintf(`docker attach --detach-keys ctrl-d %s`, stagedContainer[0].Alias)); err != nil {
		log.WithFields(log.Fields{
			"function":    "GadgetShell",
			"shell-stage": "session.Shell",
		}).Debugf("%v", err)
		return err
	}

	log.WithFields(log.Fields{
		"function": "GadgetShell",
	}).Debug("Entering container shell..")
	log.Info("Stop and enter host shell: CTRL+C")

	if terminal.IsTerminal(0) {
		oldState, err := terminal.MakeRaw(0)
		if err != nil {
			log.WithFields(log.Fields{
				"function":    "GadgetShell",
				"shell-stage": "terminal.MakeRaw",
			}).Debugf("%v", err)
			return err
		}
		defer terminal.Restore(0, oldState)
	} else {
		log.Warn("This doesn't look like a real terminal. The shell may exhibit some strange behaviour.")
	}

	session.Wait()

	log.WithFields(log.Fields{
		"function": "GadgetShell",
	}).Debug("Closed shell.")

	return err
}

// Process the build arguments and execute build
func GadgetShell(args []string, g *libgadget.GadgetContext) error {

	err := libgadget.EnsureKeys()
	if err != nil {
		log.Errorf("Failed to connect to Gadget")
		return err
	}

	// shell into a specific container
	if len(args) == 1 {
		err := GadgetShellContainer(args, g)
		if err != nil {
			log.Errorf("Failed to connect to %s", args[0])
			return err
		}
	} else if len(args) > 1 {
		log.Errorf("'gadget shell' can either take no arguments, or one argument")
		log.Warnf("'gadget shell' will ssh into the host Gadget OS")
		log.Warnf("'gadget shell <name>' will attach to the specified container")
		return errors.New("Too many arguments specified for 'gadget shell'")
	}

	client, err := libgadget.GadgetLogin(libgadget.GadgetPrivKeyLocation)
	if err != nil {
		log.WithFields(log.Fields{
			"function":    "GadgetShell",
			"shell-stage": "gadget login",
		}).Debugf("%v", err)
		return err
	}

	session, err := client.NewSession()
	if err != nil {
		log.WithFields(log.Fields{
			"function":    "GadgetShell",
			"shell-stage": "new session",
		}).Debugf("%v", err)
		return err
	}

	session.Stdout = os.Stdout
	session.Stderr = os.Stderr
	session.Stdin = os.Stdin

	modes := ssh.TerminalModes{
		ssh.ECHO:   1, // enable echoing
		ssh.ECHONL: 1,
	}
	if runtime.GOOS == "windows" {
		modes = ssh.TerminalModes{
			ssh.ECHO:   0, // disable echoing
			ssh.ECHONL: 0,
			ssh.IGNCR:  1,
		}
	}

	if err := session.RequestPty("xterm", 25, 80, modes); err != nil {
		session.Close()
		log.WithFields(log.Fields{
			"function":    "GadgetShell",
			"shell-stage": "request pty",
		}).Debugf("%v", err)
		return err
	}

	if err := session.Shell(); err != nil {
		log.WithFields(log.Fields{
			"function":    "GadgetShell",
			"shell-stage": "session.Shell",
		}).Debugf("%v", err)
		return err
	}

	log.WithFields(log.Fields{
		"function": "GadgetShell",
	}).Debug("Entering shell..")

	if terminal.IsTerminal(0) {
		oldState, err := terminal.MakeRaw(0)
		if err != nil {
			log.WithFields(log.Fields{
				"function":    "GadgetShell",
				"shell-stage": "terminal.MakeRaw",
			}).Debugf("%v", err)
			return err
		}
		defer terminal.Restore(0, oldState)
	} else {
		log.Warn("This doesn't look like a real terminal. The shell may exhibit some strange behaviour.")
	}

	session.Wait()

	log.WithFields(log.Fields{
		"function": "GadgetShell",
	}).Debug("Closed shell.")
	return nil
}
