package main

import (
	"fmt"
)

// Process the build arguments and execute build
func gadgetDelete(args []string, g *GadgetContext) {
	loadConfig(g)
	ensureKeys()

	client, err := gadgetLogin(gadgetPrivKeyLocation)

	if err != nil {
		panic(err)
	}

	for _, onboot := range g.Config.Onboot {
		commandFormat := `docker rmi %s`
		cmd := fmt.Sprintf(commandFormat, onboot.ImageAlias)
		runRemoteCommand(client, cmd)
		if err != nil {
			panic(err)
		}
	}
}
