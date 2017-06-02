package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
)

// Process the build arguments and execute build
func GadgetLogs(args []string, g *GadgetContext) error {
	
	EnsureKeys()

	client, err := GadgetLogin(gadgetPrivKeyLocation)

	if err != nil {
		return err
	}

	log.Info("[GADGT]  Retrieving logs:")
	
	for _, onboot := range g.Config.Onboot {
		commandFormat := `docker logs %s`
		cmd := fmt.Sprintf(commandFormat, onboot.Alias)
		RunRemoteCommand(client, cmd)
		if err != nil {
			return err
		}
	}
	return nil
}
