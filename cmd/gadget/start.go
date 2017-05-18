package main

import (
	"fmt"
	"strings"
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

		commandFormat := `docker start %s`
		cmd := fmt.Sprintf(commandFormat, onboot.Alias)
		runRemoteCommand(client, cmd)
		if err != nil {
			panic(err)
		}

	}
	for _, onboot := range g.Config.Services {
		commandFormat := `docker create --name %s %s %s`
		cmd := fmt.Sprintf(commandFormat, onboot.Alias, onboot.ImageAlias, strings.Join(onboot.Command[:]," "))
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
