package main

import (
	//"flag"
	"bufio"
	"fmt"
	"os"
	"os/exec"
)

// Process the build arguments and execute build
func build(args []string, g *GadgetContext) {

	loadConfig(g)

	// find docker binary in path
	binary, err := exec.LookPath("docker")
	if err != nil {
		panic(err)
	}

	// loop through 'onboot' config and build containers
	for _, onboot := range g.Config.Onboot {
		fmt.Println(" ==> Building:", onboot.Name)

		containerDirectory := fmt.Sprintf("%s/%s", g.WorkingDirectory, onboot.Directory)
		cmd := exec.Command(binary,
			"build",
			"--tag",
			onboot.ImageAlias,
			containerDirectory)
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
