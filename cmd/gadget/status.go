package main

import (
	"fmt"
)

// Process the build arguments and execute build
func gadgetStatus(args []string, g *GadgetContext) {
	loadConfig(g)
	ensureKeys()

	client, err := gadgetLogin(gadgetPrivKeyLocation)

	if err != nil {
		panic(err)
	}

	for _, onboot := range g.Config.Onboot {
		commandFormat := `docker ps -a --filter=ancestor=%s --format "{{.Image}} {{.Command}} {{.Status}}"`
		cmd := fmt.Sprintf(commandFormat, onboot.ImageAlias)
		runRemoteCommand(client, cmd)
		if err != nil {
			panic(err)
		}
	}
}
