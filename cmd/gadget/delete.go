package main

import "fmt"

// Process the build arguments and execute build
func gadgetDelete(args []string, g *GadgetContext) {
	loadConfig(g)
	ensureKeys()

	client, err := gadgetLogin(gadgetPrivKeyLocation)

	if err != nil {
		panic(err)
	}

	fmt.Println("[GADGT]  Deleting:")
	
	for _, onboot := range g.Config.Onboot {
		fmt.Printf("[GADGT]    %s ", onboot.Name)
		runRemoteCommand(client, "docker", "rmi", onboot.ImageAlias)
		fmt.Printf("✔\n")
	}

	for _, service := range g.Config.Services {
		fmt.Printf("[GADGT]    %s ", service.Name)
		runRemoteCommand(client, "docker", "rmi", service.ImageAlias)		
		fmt.Printf("✔\n")
	}
}
