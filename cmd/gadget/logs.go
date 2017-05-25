package main

import (
	"fmt"
)

// Process the build arguments and execute build
func GadgetLogs(args []string, g *GadgetContext) {
	g.LoadConfig()
	EnsureKeys()

	client, err := GadgetLogin(gadgetPrivKeyLocation)

	if err != nil {
		panic(err)
	}

	fmt.Println("[GADGT]  Retrieving logs:")
	
	for _, onboot := range g.Config.Onboot {
		commandFormat := `docker logs %s`
		cmd := fmt.Sprintf(commandFormat, onboot.Alias)
		RunRemoteCommand(client, cmd)
		if err != nil {
			panic(err)
		}
	}
}
