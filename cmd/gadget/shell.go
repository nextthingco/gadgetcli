package main

import (
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
	"os"
	log "github.com/sirupsen/logrus"
)

// Process the build arguments and execute build
func GadgetShell(args []string, g *GadgetContext) error {
	
	EnsureKeys()

	client, err := GadgetLogin(gadgetPrivKeyLocation)

	if err != nil {
		return err
	}

	session, err := client.NewSession()
	if err != nil {
		return err
	}

	session.Stdout = os.Stdout
	session.Stderr = os.Stderr
	session.Stdin  = os.Stdin

	modes := ssh.TerminalModes{
		ssh.ECHO:          1,     // disable echoing
		ssh.ECHONL:        1,
	}

	if err := session.RequestPty("xterm", 25, 80, modes); err != nil {
		session.Close()
		return err
	}
	
	if err := session.Shell(); err != nil {
		return err
	}
	
	log.WithFields(log.Fields{
		"function": "GadgetShell",
	}).Debug("Entering shell..")

	
	oldState, err := terminal.MakeRaw(0)
	if err != nil {
		return err
	}
	defer terminal.Restore(0, oldState)

	session.Wait()
	
	log.WithFields(log.Fields{
		"function": "GadgetShell",
	}).Debug("Closed shell.")
	return nil	
}
