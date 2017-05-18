package main

import (
	"fmt"
	"os"
	"os/exec"
	"io"
)

// Process the build arguments and execute build
func gadgetDeploy(args []string, g *GadgetContext) {
	binary, err := exec.LookPath("docker")
	if err != nil {
		panic(err)
	}

	loadConfig(g)
	ensureKeys()

	client, err := gadgetLogin(gadgetPrivKeyLocation)

	if err != nil {
		panic(err)
	}

	for _, onboot := range g.Config.Onboot {

		fmt.Println("==> deploying:", onboot.ImageAlias)
		docker := exec.Command(binary, "save", onboot.ImageAlias)

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
	}
}
