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

	fmt.Println("==> deploying:", container.ImageAlias)
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
	
	fmt.Println("Starting session")
	if err := session.Start(`docker load`); err != nil {
		panic(err)
	}

	fmt.Println("Starting docker")
	if err := docker.Start(); err != nil {
		panic(err)
	}


	go func() {
		defer pw.Close()
		fmt.Println("Waiting on docker")
		if err := docker.Wait(); err != nil {
			panic(err)
		}
	}()
	
	session.Wait()
	fmt.Println("closing session")
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

	loadConfig(g)
	ensureKeys()

	client, err := gadgetLogin(gadgetPrivKeyLocation)

	if err != nil {
		panic(err)
	}

	for _, onboot := range g.Config.Onboot {
		deployContainer(client, &onboot, true)
	}
	for _, service := range g.Config.Services {
		deployContainer(client, &service, false)
	}
}
