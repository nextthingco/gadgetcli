package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"io/ioutil"
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
	flagQuiet := flag.Bool("q", false, "Quiet execution")
	flagVerbose := flag.Bool("v", false, "Verbose execution")
	workingDirectory := flag.String("C", ".", "Run in directory")

	flag.Parse()
	if *flagQuiet && *flagVerbose {
		fmt.Printf("Can't set quiet and verbose flag at the same time\n")
		os.Exit(1)
	}

	args := flag.Args()
	if len(args) < 1 {
		fmt.Printf("Please specify a command.\n\n")
		flag.Usage()
		os.Exit(1)
	}

	fmt.Printf("Running in directory: %s\n", *workingDirectory)

	// read gadget.yml
	// TODO: this should probably get moved into the NewConfig function
	// TODO: look for gadget.yml in workingDirectory
	// TODO: look for gadget.yml in ancestors of workingDirectory
	config, err := ioutil.ReadFile(fmt.Sprintf("%s/gadget.yml", *workingDirectory))
	if err != nil {
		fmt.Printf("Cannot open config file: %v\n", err)
	}

	// create new config class from gadget.yml output
	// TODO: add error checking here.
	g, err := NewConfig(config)

	// parse arguments
	switch args[0] {
	case "build":
		build(args[1:], g, workingDirectory)
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
