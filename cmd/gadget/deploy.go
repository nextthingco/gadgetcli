package main

import (
	"fmt"
	"os"
	"os/exec"
	"io"
	"golang.org/x/crypto/ssh"
	"strings"
)

func deployContainer( client *ssh.Client, container * GadgetContainer, autostart bool) {
	binary, err := exec.LookPath("docker")
	if err != nil {
		panic(err)
	}
	
	fmt.Printf("[GADGT]  Deploying: %s\n", container.ImageAlias)
	docker := exec.Command(binary, "save", container.ImageAlias)

	session, err := client.NewSession()
	if err != nil {
		client.Close()
		panic(err)
	}

	// create pipe for local -> remote file transmission
	pr, pw := io.Pipe()

	docker.Stdout =  pw
	session.Stdin = pr
	session.Stdout = os.Stdout
	session.Stderr = os.Stderr
	
	fmt.Println("[GADGT]    Starting session")
	if err := session.Start(`docker load`); err != nil {
		panic(err)
	}

	fmt.Println("[GADGT]    Starting docker")
	if err := docker.Start(); err != nil {
		panic(err)
	}


	go func() {
		defer pw.Close()
		fmt.Println("[GADGT]    Waiting on docker")
		if err := docker.Wait(); err != nil {
			panic(err)
		}
	}()
	
	session.Wait()
	fmt.Println("[GADGT]    Closing session")
	session.Close()

	if autostart {
		runRemoteCommand(client, "docker",
			"create",
			"--name", container.Alias,
			"--restart=always",
			container.ImageAlias,
			strings.Join(container.Command[:]," "))
	}
}
// Process the build arguments and execute build
func gadgetDeploy(args []string, g *GadgetContext) {

	g.loadConfig()
	ensureKeys()

	client, err := gadgetLogin(gadgetPrivKeyLocation)

	if err != nil {
		panic(err)
	}

	stagedContainers,_ := findStagedContainers(args, append(g.Config.Onboot, g.Config.Services...))

	for _, container := range stagedContainers {
		deployContainer(client, &container, false)
	}
}
