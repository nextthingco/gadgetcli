package main

import (
	"fmt"
)

// Process the build arguments and execute build
func gadgetStop(args []string, g *GadgetContext) {
	g.loadConfig()
	ensureKeys()

	client, err := gadgetLogin(gadgetPrivKeyLocation)

	if err != nil {
		panic(err)
	}

	fmt.Println("[GADGT]  Stopping:")
	stagedContainers,_ := findStagedContainers(args, append(g.Config.Onboot, g.Config.Services...))
	
	for _, container := range stagedContainers {
		fmt.Printf("[GADGT]    %s ", container.Alias)
		
		err = runRemoteCommand(client, "docker stop", container.Alias)
		if err != nil {
			fmt.Printf("✘\n")
			panic(err)
		}

		err = runRemoteCommand(client, "docker rm", container.Alias)
		if err != nil {
			fmt.Printf("✘\n")
			panic(err)
		}

		fmt.Printf("✔\n")
	}
}
