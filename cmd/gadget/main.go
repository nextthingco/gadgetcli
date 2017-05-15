package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

var (
	Version = "unknown"
	GitCommit = "unknown"
)

func version() {
	fmt.Println(filepath.Base(os.Args[0]))
	fmt.Printf("  version: %s\n", Version)
	fmt.Printf("  commit: %s\n", GitCommit)
	os.Exit(0)
}

func main() {
	g := GadgetContext{}

	flag.Usage = func() {
		fmt.Printf("USAGE: %s [options] COMMAND\n\n", filepath.Base(os.Args[0]))
		fmt.Printf("Commands:\n")
		fmt.Printf("  init        Initialize gadget project\n")
		fmt.Printf("  build       Build gadget config file\n")
		fmt.Printf("  deploy      Build gadget config file\n")
		fmt.Printf("  start       Build gadget config file\n")
		fmt.Printf("  stop        Build gadget config file\n")
		fmt.Printf("  status      Build gadget config file\n")
		fmt.Printf("  delete      Build gadget config file\n")
		fmt.Printf("  shell       Connect to remote device running GadgetOS\n")
		fmt.Printf("  log         Build gadget config file\n")
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
	case "build":
		build(args[1:], &g)
	case "init":
		gadgetInit(args[1:], &g)
	//	case "ssh":
	//		shell(args[1:])
	case "version":
		version()
	case "help":
		flag.Usage()
	default:
		fmt.Printf("%q is not valid command.\n\n", args[0])
		flag.Usage()
		os.Exit(1)
	}
}
