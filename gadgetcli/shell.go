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
	"github.com/nextthingco/libgadget"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
	"os"
	"runtime"
)

// Process the build arguments and execute build
func GadgetShell(args []string, g *libgadget.GadgetContext) error {

	libgadget.EnsureKeys()

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
		ssh.ECHO:   1, // disable echoing
		ssh.ECHONL: 1,
	}
	if runtime.GOOS == "windows" {
		modes = ssh.TerminalModes{
			ssh.ECHO:   0, // disable echoing
			ssh.ECHONL: 0,
			ssh.IGNCR: 1,
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
