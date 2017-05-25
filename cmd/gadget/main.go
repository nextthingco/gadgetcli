package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

var (
	Version   = "unknown"
	GitCommit = "unknown"
)

func GadgetVersion() {
	fmt.Println(filepath.Base(os.Args[0]))
	fmt.Printf("  version: %s\n", Version)
	fmt.Printf("  commit: %s\n", GitCommit)
	os.Exit(0)
}

func main() {
	g := GadgetContext{}

	err := requiredSsh()
	if err != nil {
		panic(err)
	}

	flag.Usage = func() {
		fmt.Printf("USAGE: %s [options] COMMAND\n\n", filepath.Base(os.Args[0]))
		fmt.Printf("Commands:\n")
		fmt.Printf("  init        Initialize gadget project\n")
		fmt.Printf("  add         Initialize gadget project\n")
		fmt.Printf("  build       Build gadget config file\n")
		fmt.Printf("  deploy      Build gadget config file\n")
		fmt.Printf("  start       Build gadget config file\n")
		fmt.Printf("  stop        Build gadget config file\n")
		fmt.Printf("  status      Build gadget config file\n")
		fmt.Printf("  delete      Build gadget config file\n")
		fmt.Printf("  shell       Connect to remote device running GadgetOS\n")
		fmt.Printf("  logs        Build gadget config file\n")
		fmt.Printf("  version     Print version information\n")
		fmt.Printf("  help        Print this message\n")
		fmt.Printf("\n")
		fmt.Printf("Run '%s COMMAND --help' for more information on the command\n", filepath.Base(os.Args[0]))
		fmt.Printf("\n")
		fmt.Printf("Options:\n")
		flag.PrintDefaults()
	}
	flag.BoolVar(&g.Verbose, "v", false, "Verbose execution")
	flag.StringVar(&g.WorkingDirectory, "C", ".", "Run in directory")

	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		fmt.Printf("Please specify a command.\n\n")
		flag.Usage()
		os.Exit(1)
	}

	// parse arguments
	switch args[0] {
	case "init":
		GadgetInit(args[1:], &g)
	case "add":
		GadgetAdd(args[1:], &g)
	case "build":
		GadgetBuild(args[1:], &g)
	case "deploy":
		GadgetDeploy(args[1:], &g)
	case "start":
		GadgetStart(args[1:], &g)
	case "stop":
		GadgetStop(args[1:], &g)
	case "status":
		GadgetStatus(args[1:], &g)
	case "delete":
		GadgetDelete(args[1:], &g)
	case "shell":
		GadgetShell(args[1:])
	case "logs":
		GadgetLogs(args[1:], &g)
	case "version":
		GadgetVersion()
	case "help":
		flag.Usage()
	case "run":
		GadgetRun(args[1:], &g)
	default:
		fmt.Printf("%q is not valid command.\n\n", args[0])
		flag.Usage()
		os.Exit(1)
	}
}
