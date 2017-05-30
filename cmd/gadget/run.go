package main

import (
	//"fmt"
	"strings"
	log "github.com/sirupsen/logrus"
)

func GadgetRun(args []string, g *GadgetContext) error {
	//g.LoadConfig()
	EnsureKeys()

	client, err := GadgetLogin(gadgetPrivKeyLocation)
	if err != nil {
		return err
	}

	stdout, stderr, err := RunRemoteCommand(client, strings.Join(args, " "))

	log.WithFields(log.Fields{
		"command": strings.Join(args, " "),
		"stdout": stdout.String(),
		"stderr": stderr.String(),
	}).Info("Ran remote command")

	if err != nil {
		return err
	}
	return nil
}
