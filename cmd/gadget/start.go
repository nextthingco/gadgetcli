package main

import (
	"fmt"
	"strings"
)

// Process the build arguments and execute build
func GadgetStart(args []string, g *GadgetContext) {
	g.LoadConfig()
	EnsureKeys()

	client, err := GadgetLogin(gadgetPrivKeyLocation)

	if err != nil {
		panic(err)
	}

	fmt.Println("[GADGT]  Starting:")
	stagedContainers,_ := FindStagedContainers(args, append(g.Config.Onboot, g.Config.Services...))
	for _, container := range stagedContainers {
		
		fmt.Printf("[GADGT]    %s ", container.Alias)
		binds := strings.Join( PrependToStrings(container.Binds[:],"-v "), " ")
		commands := strings.Join(container.Command[:]," ")
		
		stdout, stderr, err := RunRemoteCommand(client, "docker create --name", container.Alias, binds, container.ImageAlias, commands)
		fmt.Println(stdout)
		fmt.Println(stderr)
		if err != nil {
			fmt.Printf("✘\n")
			panic(err)
		}

		stdout, stderr, err = RunRemoteCommand(client, "docker start", container.Alias)
		fmt.Println(stdout)
		fmt.Println(stderr)
		if err != nil {
			fmt.Printf("✘\n")
			panic(err)
		}

		fmt.Printf("✔\n")
	}
}
