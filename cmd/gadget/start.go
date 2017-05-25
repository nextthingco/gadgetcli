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

	stagedContainers := findStagedContainers(args, append(g.Config.Onboot, g.Config.Services...))
	
	for _, container := range stagedContainers {
		
		fmt.Printf("[GADGT]    %s ", container.Alias)

		binds := strings.Join( prependToStrings(container.Binds[:],"-v "), " ")
		commands := strings.Join(container.Command[:]," ")
		
		err = runRemoteCommand(client, "docker create --name", container.Alias, binds, container.ImageAlias, commands)
		if err != nil {
			fmt.Printf("✘\n")
			panic(err)
		}

		err = runRemoteCommand(client, "docker start", container.Alias)
		if err != nil {
			fmt.Printf("✘\n")
			panic(err)
		}

		fmt.Printf("✔\n")
	}
}
