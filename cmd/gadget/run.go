package main

import (
	"fmt"
	"strings"
)

func GadgetRun(args []string, g *GadgetContext) {
	g.loadConfig()
	ensureKeys()

	client, err := gadgetLogin(gadgetPrivKeyLocation)
	if err != nil {
		panic(err)
	}

	stdout, stderr, err := runRemoteCommand(client, strings.Join(args, " "))
	fmt.Println(stdout.String())
	fmt.Println(stderr.String())
	if err != nil {
		panic(err)
	}
}
