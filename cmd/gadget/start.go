package main

import (
	"fmt"
	"os"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
	"syscall"
)


// Process the build arguments and execute build
func gadgetStart(args []string, g *GadgetContext) {
	// find docker binary in path

	config := &ssh.ClientConfig{
		User: os.Getenv("LOGNAME"),
		Auth: []ssh.AuthMethod{
			ssh.RetryableAuthMethod(ssh.PasswordCallback(
				func() (secret string, err error) {
					fmt.Print("Password:")
					password,err := terminal.ReadPassword(syscall.Stdin)
					fmt.Println("")
					return string(password),err
				}), 3),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	client, err := ssh.Dial("tcp", "localhost:22", config)
	if err != nil {
		panic(err)
	}

	session, err := client.NewSession()
	if err != nil {
		client.Close()
		panic(err)
	}

	fmt.Println("Running docker info...")
	// docker create --name ${name}_${uuid} ${name}_${uuid}-img
	out, err := session.CombinedOutput("printenv SHELL")
	if err != nil {
		panic(err)
	}

	fmt.Print( string(out) )
}
