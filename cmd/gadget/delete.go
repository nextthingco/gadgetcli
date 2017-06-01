package main


import (
	"fmt"
	log "github.com/sirupsen/logrus"
)

// Process the build arguments and execute build
func GadgetDelete(args []string, g *GadgetContext) error {
	//~ g.LoadConfig()
	EnsureKeys()

	client, err := GadgetLogin(gadgetPrivKeyLocation)

	if err != nil {
		return err
	}

	log.Info(fmt.Sprintf("  Deleting:"))

	stagedContainers,_ := FindStagedContainers(args, append(g.Config.Onboot, g.Config.Services...))
	
	for _, container := range stagedContainers {
		log.Info(fmt.Sprintf("    %s ", container.ImageAlias))
		stdout, stderr, err := RunRemoteCommand(client, "docker", "rmi", container.ImageAlias)
		if err != nil {
			//~ fmt.Printf("✘\n")
			return err
		}
		fmt.Println(stdout)
		fmt.Println(stderr)
		
		//~ fmt.Printf("✔\n")
	}
	return nil
}
