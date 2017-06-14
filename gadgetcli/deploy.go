package main

import (
	"os/exec"
	"io"
	"errors"
	"golang.org/x/crypto/ssh"
	"gopkg.in/cheggaaa/pb.v1"
	"strings"
	"github.com/nextthingco/libgadget"
	log "github.com/sirupsen/logrus"
)

func DeployContainer( client *ssh.Client, container *libgadget.GadgetContainer,g *libgadget.GadgetContext, autostart bool) error {
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
		stdout, stderr, err := libgadget.RunRemoteCommand(client, "docker",
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
func GadgetDeploy(args []string, g *libgadget.GadgetContext) error {
	
	err := libgadget.EnsureKeys()
	if err != nil {
		log.Errorf("Failed to connect to Gadget")
		return err
	}

	client, err := libgadget.GadgetLogin(libgadget.GadgetPrivKeyLocation)
	if err != nil {
		log.Errorf("Failed to connect to Gadget")
		return err
	}

	stagedContainers, err := libgadget.FindStagedContainers(args, append(g.Config.Onboot, g.Config.Services...))
	
	deployFailed := false
	
	for _, container := range stagedContainers {
		
		// stop and delete possible older versions of image/container
		// not collecting the errors, as errors may be returned
		// when trying to delete an img/cntnr that was never deployed
		tmpName := make( []string, 1 )
		tmpName[0] = container.Name
		
		log.Infof("Stopping/deleting older '%s' if applicable", container.Name)
		
		if ! g.Verbose {
			log.SetLevel(log.PanicLevel)
		}
		
		_ = GadgetStop(tmpName, g)
		_ = GadgetDelete(tmpName, g)
		
		if ! g.Verbose {
			log.SetLevel(log.InfoLevel)
		}
		
		err = DeployContainer(client, &container, g, false)
		deployFailed = true
	}
	
	if deployFailed == true {
		err = errors.New("Failed to deploy one or more containers")
	}
	
	return err
}
