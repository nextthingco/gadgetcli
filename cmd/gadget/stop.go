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


	fmt.Println("[GADGT]  Stopping:")
	fmt.Println("[GADGT]    - Onboot:")

	for _, onboot := range g.Config.Onboot {
		fmt.Printf("[GADGT]    %s ", onboot.Alias)
		
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

	fmt.Println("[GADGT]    - Services:")
	
	for _, service := range g.Config.Services {
		fmt.Printf("[GADGT]    %s ", service.Alias)
		
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
