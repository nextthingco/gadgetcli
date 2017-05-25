package main

import (
	"fmt"
)

// Process the build arguments and execute build
func GadgetStop(args []string, g *GadgetContext) {
	g.LoadConfig()
	EnsureKeys()

	client, err := GadgetLogin(gadgetPrivKeyLocation)

	if err != nil {
		panic(err)
	}

	fmt.Println("[GADGT]  Stopping:")
	stagedContainers,_ := FindStagedContainers(args, append(g.Config.Onboot, g.Config.Services...))
	
	for _, container := range stagedContainers {
		fmt.Printf("[GADGT]    %s ", container.Alias)
		
		stdout, stderr, err := RunRemoteCommand(client, "docker stop", container.Alias)
		fmt.Println(stdout)
		fmt.Println(stderr)
		if err != nil {
			fmt.Printf("✘\n")
			panic(err)
		}

		stdout, stderr, err = RunRemoteCommand(client, "docker rm", container.Alias)
		fmt.Println(stdout)
		fmt.Println(stderr)
		if err != nil {
			fmt.Printf("✘\n")
			panic(err)
		}

		fmt.Printf("✔\n")
	}
}
