package main

import (
	"strings"
	"github.com/nextthingco/libgadget"
	log "github.com/sirupsen/logrus"
)

func GadgetRun(args []string, g *libgadget.GadgetContext) error {
	
	libgadget.EnsureKeys()

	client, err := libgadget.GadgetLogin(libgadget.GadgetPrivKeyLocation)
	if err != nil {
		return err
	}

	stdout, stderr, err := libgadget.RunRemoteCommand(client, strings.Join(args, " "))
	
	if err != nil {
		log.Errorf("\n%s", stdout)
		log.Errorf("\n%s", stderr)
		return err
	}
	
	log.Infof("\n%s", stdout)
	log.Debugf("\n%s", stderr)
	
	return err
}
