package main

import (
	"strings"
	log "github.com/sirupsen/logrus"
)

func GadgetRun(args []string, g *GadgetContext) error {
	
	EnsureKeys()

	client, err := GadgetLogin(gadgetPrivKeyLocation)
	if err != nil {
		return err
	}

	stdout, stderr, err := RunRemoteCommand(client, strings.Join(args, " "))
	
	if err != nil {
		log.Errorf("\n%s", stdout)
		log.Errorf("\n%s", stderr)
		return err
	}
	
	log.Infof("\n%s", stdout)
	log.Debugf("\n%s", stderr)
	
	return err
}
