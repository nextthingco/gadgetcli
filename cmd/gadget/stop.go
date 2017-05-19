package main

import (
	"fmt"
)

// Process the build arguments and execute build
func gadgetStop(args []string, g *GadgetContext) {
	loadConfig(g)
	ensureKeys()

	client, err := gadgetLogin(gadgetPrivKeyLocation)

	if err != nil {
		panic(err)
	}

	for _, onboot := range g.Config.Onboot {
		commandFormat := `docker stop %s`
		cmd := fmt.Sprintf(commandFormat, onboot.Alias)
		runRemoteCommand(client, cmd)
		if err != nil {
			panic(err)
		}

		commandFormat = `docker rm %s`
		cmd = fmt.Sprintf(commandFormat, onboot.Alias)
		runRemoteCommand(client, cmd)
		if err != nil {
			panic(err)
		}
	}

	for _, service := range g.Config.Services {
		commandFormat := `docker stop %s`
		cmd := fmt.Sprintf(commandFormat, service.Alias)
		runRemoteCommand(client, cmd)
		if err != nil {
			panic(err)
		}

		commandFormat = `docker rm %s`
		cmd = fmt.Sprintf(commandFormat, service.Alias)
		runRemoteCommand(client, cmd)
		if err != nil {
			panic(err)
		}
	}
}
