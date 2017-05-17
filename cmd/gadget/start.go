package main

import (
	"fmt"
)

// Process the build arguments and execute build
func gadgetStart(args []string, g *GadgetContext) {
	loadConfig(g)
	ensureKeys()

	client, err := gadgetLogin(gadgetPrivKeyLocation)

	if err != nil {
		panic(err)
	}

	for _, onboot := range g.Config.Onboot {
		commandFormat := `docker create --name %s %s`
		cmd := fmt.Sprintf(commandFormat, onboot.Alias, onboot.ImageAlias)
		runRemoteCommand(client, cmd)
		if err != nil {
			panic(err)
		}

		commandFormat = `docker start %s`
		cmd = fmt.Sprintf(commandFormat, onboot.Alias)
		runRemoteCommand(client, cmd)
		if err != nil {
			panic(err)
		}

	}
}
