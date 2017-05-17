package main

import (
	"fmt"
	"os"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
	"syscall"
	"sync"
)


// Process the build arguments and execute build
func gadgetStart(args []string, g *GadgetContext) {
	loadConfig(g)
	
	ensureKeys()
	
	sshConfig := &ssh.ClientConfig{
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

	client, err := ssh.Dial("tcp", "localhost:22", sshConfig)
	if err != nil {
		panic(err)
	}


	outputChannel := make(chan string)

	wg := sync.WaitGroup{}
	for _,container := range g.Config.Onboot {
		wg.Add(1)
		go func(container GadgetContainer) {
			defer wg.Done()
			session, err := client.NewSession()
			if err != nil {
				client.Close()
				panic(err)
			}
			// docker create --name ${name}_${uuid} ${name}_${uuid}-img
			out, err := session.CombinedOutput("printenv SHELL")
			if err != nil {
				panic(err)
			}
			outputString := fmt.Sprintf("%s: %s", container.Name, string(out))
			outputChannel <- outputString
		}(container)
	}
	go func() {
		for data := range outputChannel {
			fmt.Print(data)
		}
	}()
	wg.Wait()
}
