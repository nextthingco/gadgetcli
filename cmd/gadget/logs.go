package main

import (
	"fmt"
)

// Process the build arguments and execute build
func gadgetLogs(args []string, g *GadgetContext) {
	loadConfig(g)
	ensureKeys()

	client, err := gadgetLogin(gadgetPrivKeyLocation)

	if err != nil {
		panic(err)
	}

	for _, onboot := range g.Config.Onboot {
		commandFormat := `docker logs $(docker ps -aq --filter ancestor=%s)`
		cmd := fmt.Sprintf(commandFormat, onboot.ImageAlias)
		runRemoteCommand(client, cmd)
		if err != nil {
			panic(err)
		}
	}
}
