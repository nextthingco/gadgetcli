package main

import (
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
	"os"
	"fmt"
)

// Process the build arguments and execute build
func GadgetShell(args []string) {
	
	EnsureKeys()

	client, err := GadgetLogin(gadgetPrivKeyLocation)

	if err != nil {
		panic(err)
	}

	session, err := client.NewSession()
	if err != nil {
		panic(err)
	}

	session.Stdout = os.Stdout
	session.Stderr = os.Stderr
	session.Stdin  = os.Stdin

	modes := ssh.TerminalModes{
		ssh.ECHO:          1,     // disable echoing
		ssh.ECHONL:        1,
	}

	if err := session.RequestPty("xterm", 25, 80, modes); err != nil {
		session.Close()
		panic(err)
	}
	
	if err := session.Shell(); err != nil {
		panic(err)
	}
	
	fmt.Println("[COMMS]  Entering shell..")
	
	oldState, err := terminal.MakeRaw(0)
	if err != nil {
	        panic(err)
	}
	defer terminal.Restore(0, oldState)

	session.Wait()
	
	fmt.Println("[COMMS]  Closed shell.")
	
}
