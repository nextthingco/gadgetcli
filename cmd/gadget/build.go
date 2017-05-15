package main

import (
	//"flag"
	"fmt"
	"os"
	"os/exec"
	"bufio"
)

// Process the build arguments and execute build
func build(args []string, g *GadgetContext) {
	// find docker binary in path
	binary, lookErr := exec.LookPath("docker")
	if lookErr != nil {
		panic(lookErr)
	}

	// loop through 'onboot' config and build containers
	for _,container := range g.Config.Onboot {
		fmt.Println(" ==> Building:", container.Name)

		cmd := exec.Command(binary,
			"build",
			"--tag",
			fmt.Sprintf("%s_%s-img", container.Name, container.UUID), //"gadget-networkd_112icx9s-img",
			fmt.Sprintf("%s/%s", g.WorkingDirectory, container.From)) //
		cmd.Env = os.Environ()
		
		stdOutReader, execErr := cmd.StdoutPipe()
		stdErrReader, execErr := cmd.StderrPipe()
		outScanner := bufio.NewScanner(stdOutReader)
		errScanner := bufio.NewScanner(stdErrReader)

		// goroutine to print stdout and stderr
		go func() {
			// TODO: goroutine gets launched and never exits.
			for {
				// TODO: add a check here to only print stdout if verbose
				/*if outScanner.Scan() {
					fmt.Println(string(outScanner.Text()))
				}*/
				_ = outScanner.Scan()
				if errScanner.Scan() {
					fmt.Println(string(errScanner.Text()))
				}
			}
		}()

		execErr = cmd.Start()
		if execErr != nil {
			panic(execErr)
		}
		execErr = cmd.Wait()
		if execErr != nil {
			panic(execErr)
		}
	}
}
