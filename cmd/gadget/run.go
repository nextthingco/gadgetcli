package main

import (
	"fmt"
	"strings"
)

func GadgetRun(args []string, g *GadgetContext) error {
	//g.LoadConfig()
	EnsureKeys()

	client, err := GadgetLogin(gadgetPrivKeyLocation)
	if err != nil {
		return err
	}

	stdout, stderr, err := RunRemoteCommand(client, strings.Join(args, " "))
	fmt.Println(stdout.String())
	fmt.Println(stderr.String())
	if err != nil {
		return err
	}
	return nil
}
