package main

import "fmt"

// Process the build arguments and execute build
func gadgetDelete(args []string, g *GadgetContext) {
	g.loadConfig()
	ensureKeys()

	client, err := gadgetLogin(gadgetPrivKeyLocation)

	if err != nil {
		panic(err)
	}

	fmt.Println("[GADGT]  Deleting:")

	stagedContainers,_ := findStagedContainers(args, append(g.Config.Onboot, g.Config.Services...))
	
	for _, container := range stagedContainers {
		fmt.Printf("[GADGT]    %s ", container.ImageAlias)
		err = runRemoteCommand(client, "docker", "rmi", container.ImageAlias)
		if err != nil {
			fmt.Printf("✘\n")
			panic(err)
		}
		
		fmt.Printf("✔\n")
	}
}
