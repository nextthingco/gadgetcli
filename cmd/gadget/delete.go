package main

import "fmt"

// Process the build arguments and execute build
func GadgetDelete(args []string, g *GadgetContext) {
	g.LoadConfig()
	EnsureKeys()

	client, err := GadgetLogin(gadgetPrivKeyLocation)

	if err != nil {
		panic(err)
	}

	fmt.Println("[GADGT]  Deleting:")

	stagedContainers,_ := FindStagedContainers(args, append(g.Config.Onboot, g.Config.Services...))
	
	for _, container := range stagedContainers {
		fmt.Printf("[GADGT]    %s ", container.ImageAlias)
		stdout, stderr, err := RunRemoteCommand(client, "docker", "rmi", container.ImageAlias)
		if err != nil {
			fmt.Printf("✘\n")
			panic(err)
		}
		fmt.Println(stdout)
		fmt.Println(stderr)
		
		fmt.Printf("✔\n")
	}
}
