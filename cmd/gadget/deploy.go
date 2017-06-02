package main

import (
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
	
	log.Info("[GADGT]  Deploying:")
	log.Infof("[GADGT]    %s", container.ImageAlias)
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
	
	log.Debug("[GADGT]    Starting session")
	if err := session.Start(`docker load`); err != nil {
		return err
	}

	log.Debug("[GADGT]    Starting docker")
	if err := docker.Start(); err != nil {
		return err
	}


	go func() error {
		defer pw.Close()
		log.Info("[GADGT]    Starting transfer..")
		log.Debug("[GADGT]    Waiting on docker")
		if err := docker.Wait(); err != nil {
			// TODO: we should handle this error or report to the log
			return err
		}
		return err
	}()
	
	session.Wait()
	log.Info("[GADGT]    Done!")
	log.Debug("[GADGT]    Closing session")
	session.Close()
		
	if autostart {
		stdout, stderr, err := RunRemoteCommand(client, "docker",
			"create",
			"--name", container.Alias,
			"--restart=always",
			container.ImageAlias,
			strings.Join(container.Command[:]," "))
		
		if err != nil {
			log.Errorf("Failed to set %s to always restart on Gadget", container.Alias)
			return err
		}
		
		log.WithFields(log.Fields{
			"function": "DeployContainer",
			"name": container.Alias,
			"deploy-stage": "create restarting",
		}).Debug(stdout)
		log.WithFields(log.Fields{
			"function": "DeployContainer",
			"name": container.Alias,
			"deploy-stage": "create restarting",
		}).Debug(stderr)
		
	}
	
	return err
}

// Process the build arguments and execute build
func GadgetDeploy(args []string, g *GadgetContext) error {
	
	EnsureKeys()

	client, err := GadgetLogin(gadgetPrivKeyLocation)

	if err != nil {
		return err
	}

	stagedContainers, err := FindStagedContainers(args, append(g.Config.Onboot, g.Config.Services...))

	for _, container := range stagedContainers {
		DeployContainer(client, &container, g, false)
	}
	
	return err
}
