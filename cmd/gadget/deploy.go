package main

import (
	"os/exec"
	"io"
	"errors"
	"golang.org/x/crypto/ssh"
	"gopkg.in/cheggaaa/pb.v1"
	"strings"
	log "github.com/sirupsen/logrus"
)

func DeployContainer( client *ssh.Client, container * GadgetContainer,g *GadgetContext, autostart bool) error {
	binary, err := exec.LookPath("docker")
	if err != nil {
		return err
	}
		
	log.Infof("Deploying: '%s'", container.Name)
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
	
	bar := pb.New(0)
	bar.SetUnits(pb.U_BYTES)
	bar.ShowSpeed = true
	bar.ShowPercent = false
	bar.ShowTimeLeft = false
	bar.ShowBar = false
	
	docker.Stdout = pw
	reader := bar.NewProxyReader(pr)
	session.Stdin = reader
	session.Stdout = sessionLogger.WriterLevel(log.DebugLevel)
	session.Stderr = sessionLogger.WriterLevel(log.DebugLevel)
	
	log.Debug("  Starting session")
	if err := session.Start(`docker load`); err != nil {
		return err
	}

	log.Debug("  Starting docker")
	if err := docker.Start(); err != nil {
		return err
	}

	deployFailed := false
	
	go func() error {
		defer pw.Close()
		log.Info("  Starting transfer..")
		log.Debug("  Waiting on docker")
		bar.Start()
		
		if err := docker.Wait(); err != nil {
			deployFailed = true
			// TODO: we should handle this error or report to the log
			log.Errorf("Failed to transfer '%s'", container.Name)
			log.Warn("Was the container ever built?")
			return err
		}
		return err
	}()
	
	session.Wait()
	bar.Finish()
	if ! deployFailed {
		log.Info("Done!")
		log.Debug("Closing session")
	}
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
	
	err := EnsureKeys()
	if err != nil {
		log.Errorf("Failed to connect to Gadget")
		return err
	}

	client, err := GadgetLogin(gadgetPrivKeyLocation)
	if err != nil {
		log.Errorf("Failed to connect to Gadget")
		return err
	}

	stagedContainers, err := FindStagedContainers(args, append(g.Config.Onboot, g.Config.Services...))
	
	deployFailed := false
	
	for _, container := range stagedContainers {
		err = DeployContainer(client, &container, g, false)
		deployFailed = true
	}
	
	if deployFailed == true {
		err = errors.New("Failed to deploy one or more containers")
	}
	
	return err
}
