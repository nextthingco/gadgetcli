package main

import (
	"fmt"
	"os/exec"
	"io"
	"golang.org/x/crypto/ssh"
	"strings"
	log "github.com/sirupsen/logrus"
)

func DeployContainer( client *ssh.Client, container * GadgetContainer,g *GadgetContext, autostart bool) error {
	binary, err := exec.LookPath("docker")
	if err != nil {
		return err
	}
	
	log.Info(fmt.Sprintf("[GADGT]  Deploying:"))
	log.Info(fmt.Sprintf("[GADGT]    %s", container.ImageAlias))
	docker := exec.Command(binary, "save", container.ImageAlias)

	session, err := client.NewSession()
	if err != nil {
		client.Close()
		return err
	}

	// create pipe for local -> remote file transmission
	pr, pw := io.Pipe()
	sessionLogger := log.New()
	if g.Verbose { sessionLogger.Level = log.DebugLevel }
	
	docker.Stdout =  pw
	session.Stdin = pr
	session.Stdout = sessionLogger.WriterLevel(log.DebugLevel)
	session.Stderr = sessionLogger.WriterLevel(log.DebugLevel)
	
	log.Debug(fmt.Sprintf("[GADGT]    Starting session"))
	if err := session.Start(`docker load`); err != nil {
		return err
	}

	log.Debug(fmt.Sprintf("[GADGT]    Starting docker"))
	if err := docker.Start(); err != nil {
		return err
	}


	go func() error {
		defer pw.Close()
		log.Info(fmt.Sprintf("[GADGT]    Starting transfer.."))
		log.Debug(fmt.Sprintf("[GADGT]    Waiting on docker"))
		if err := docker.Wait(); err != nil {
			// TODO: we should handle this error or report to the log
			return err
		}
		return err
	}()
	
	session.Wait()
	log.Info(fmt.Sprintf("[GADGT]    Done!"))
	log.Debug(fmt.Sprintf("[GADGT]    Closing session"))
	session.Close()
		
	if autostart {
		RunRemoteCommand(client, "docker",
			"create",
			"--name", container.Alias,
			"--restart=always",
			container.ImageAlias,
			strings.Join(container.Command[:]," "))
	}
	
	return err
}

// Process the build arguments and execute build
func GadgetDeploy(args []string, g *GadgetContext) error {

	//~ g.LoadConfig()
	EnsureKeys()

	client, err := GadgetLogin(gadgetPrivKeyLocation)

	if err != nil {
		return err
	}

	stagedContainers,_ := FindStagedContainers(args, append(g.Config.Onboot, g.Config.Services...))

	for _, container := range stagedContainers {
		DeployContainer(client, &container, g, false)
	}
	
	return err
}
