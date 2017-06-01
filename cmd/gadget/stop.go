package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
)

// Process the build arguments and execute build
func GadgetStop(args []string, g *GadgetContext) error {
	//~ g.LoadConfig()
	EnsureKeys()

	client, err := GadgetLogin(gadgetPrivKeyLocation)

	if err != nil {
		return err
	}

	log.Info(fmt.Sprintf("[GADGT]  Stopping:"))
	stagedContainers,_ := FindStagedContainers(args, append(g.Config.Onboot, g.Config.Services...))
	
	for _, container := range stagedContainers {
		log.Info(fmt.Sprintf("[GADGT]    %s", container.Alias))
		
		stdout, stderr, err := RunRemoteCommand(client, "docker stop", container.Alias)
		fmt.Println(stdout)
		fmt.Println(stderr)
		if err != nil {
			//~ fmt.Printf("✘\n")
			return err
		}

		stdout, stderr, err = RunRemoteCommand(client, "docker rm", container.Alias)
		fmt.Println(stdout)
		fmt.Println(stderr)
		if err != nil {
			//~ fmt.Printf("✘\n")
			return err
		}

		//~ fmt.Printf("✔\n")
	}
	return nil
}
