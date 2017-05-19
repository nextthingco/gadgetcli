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


	fmt.Println("[GADGT]  Starting:")
	fmt.Println("[GADGT]    Onboot:")

	for _, onboot := range g.Config.Onboot {

		fmt.Printf("[GADGT]      %s ", onboot.Alias)
		
		commandFormat := `docker start %s`
		cmd := fmt.Sprintf(commandFormat, onboot.Alias)
		runRemoteCommand(client, cmd)
		if err != nil {
			panic(err)
		}
	}
	
	fmt.Println("[GADGT]    Services:")
	
	for _, onboot := range g.Config.Services {
		
		fmt.Printf("[GADGT]      %s ", onboot.Alias)
		
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

		fmt.Printf("✔\n")
	}
}
