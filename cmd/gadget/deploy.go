package main

import (
	"fmt"
	"os"
	"os/exec"
	"io"
	"golang.org/x/crypto/ssh"
	"strings"
)

func DeployContainer( client *ssh.Client, container * GadgetContainer, autostart bool) error {
	binary, err := exec.LookPath("docker")
	if err != nil {
		return err
	}
	
	fmt.Printf("[GADGT]  Deploying: %s\n", container.ImageAlias)
	docker := exec.Command(binary, "save", container.ImageAlias)

	session, err := client.NewSession()
	if err != nil {
		client.Close()
		return err
	}

	// create pipe for local -> remote file transmission
	pr, pw := io.Pipe()

	docker.Stdout =  pw
	session.Stdin = pr
	session.Stdout = os.Stdout
	session.Stderr = os.Stderr
	
	fmt.Println("[GADGT]    Starting session")
	if err := session.Start(`docker load`); err != nil {
		return err
	}

	fmt.Println("[GADGT]    Starting docker")
	if err := docker.Start(); err != nil {
		return err
	}


	go func() error {
		defer pw.Close()
		fmt.Println("[GADGT]    Waiting on docker")
		if err := docker.Wait(); err != nil {
			// TODO: we should handle this error or report to the log
			return err
		}
		return nil
	}()
	
	session.Wait()
	fmt.Println("[GADGT]    Closing session")
	session.Close()

	if autostart {
		RunRemoteCommand(client, "docker",
			"create",
			"--name", container.Alias,
			"--restart=always",
			container.ImageAlias,
			strings.Join(container.Command[:]," "))
	}
	return nil
}
// Process the build arguments and execute build
func GadgetDeploy(args []string, g *GadgetContext) error {

	g.LoadConfig()
	EnsureKeys()

	client, err := GadgetLogin(gadgetPrivKeyLocation)

	if err != nil {
		return err
	}

	stagedContainers,_ := FindStagedContainers(args, append(g.Config.Onboot, g.Config.Services...))

	for _, container := range stagedContainers {
		DeployContainer(client, &container, false)
	}
	return nil
}
