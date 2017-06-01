package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"errors"
	log "github.com/sirupsen/logrus"
)

var (
	Version   = "unknown"
	GitCommit = "unknown"
)

type GadgetCommandFunc func([]string, *GadgetContext) error

type GadgetCommand struct {
	Name        string
	Function    GadgetCommandFunc
	NeedsConfig bool
}

var Commands = []GadgetCommand {
	{ Name: "init",    Function: GadgetInit,    NeedsConfig: false },
	{ Name: "add",     Function: GadgetAdd,     NeedsConfig: true  },
	{ Name: "build",   Function: GadgetBuild,   NeedsConfig: true  },
	{ Name: "deploy",  Function: GadgetDeploy,  NeedsConfig: true  },
	{ Name: "start",   Function: GadgetStart,   NeedsConfig: true  },
	{ Name: "stop",    Function: GadgetStop,    NeedsConfig: true  },
	{ Name: "status",  Function: GadgetStatus,  NeedsConfig: true  },
	{ Name: "delete",  Function: GadgetDelete,  NeedsConfig: true  },
	{ Name: "shell",   Function: GadgetShell,   NeedsConfig: false },
	{ Name: "logs",    Function: GadgetLogs,    NeedsConfig: true  },
	{ Name: "run",     Function: GadgetRun,     NeedsConfig: false },
	{ Name: "version", Function: GadgetVersion, NeedsConfig: false },
	{ Name: "help",    Function: GadgetHelp,    NeedsConfig: false },
}

func GadgetVersion(args []string, g *GadgetContext) error {
	fmt.Println(filepath.Base(os.Args[0]))
	fmt.Printf("  version: %s\n", Version)
	fmt.Printf("  commit: %s\n", GitCommit)
	return nil
}

func GadgetHelp(args []string, g *GadgetContext) error {
	flag.Usage()
	return nil
}

func findCommand(name string) (*GadgetCommand, error) {
	for _,cmd := range Commands {
		if cmd.Name == name {
			return &cmd,nil
		}
	}
	return nil, errors.New("ERROR: failed to find command")
}

func main() {
	// Hey, Listen! 
	// Everything that outputs needs to come after g.Verbose check!
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
		fmt.Printf("\n")
	}

	g := GadgetContext{}
	
	flag.BoolVar(&g.Verbose, "v", false, "Verbose execution")
	flag.StringVar(&g.WorkingDirectory, "C", ".", "Run in directory")
	flag.Parse()

	if g.Verbose {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}
	
	// Hey, Listen! 
	// Everything that outputs needs to come after g.Verbose check!
	

	err := RequiredSsh()
	if err != nil {
		fmt.Printf("ERROR: failed at RequiredSsh in main()")
		os.Exit(1)
	}

	args := flag.Args()
	if len(args) < 1 {
		flag.Usage()
		log.Error("No Command Specified")
		os.Exit(1)
	}
		
	// file command
	cmd,err := findCommand(args[0])
	if err != nil {
		flag.Usage()
		log.WithFields(log.Fields{
			"command": strings.Join(args[0:], " "),
		}).Error("Command is not valid")
		os.Exit(1)
	}

	// if command needs to use the config file, load it
	if cmd.NeedsConfig {
		err = g.LoadConfig()
		if err != nil {
			log.Error("Failed to load config")
			os.Exit(1)
		}
	}

	err = cmd.Function(args[1:], &g)
	if err != nil {
		os.Exit(1)
	}
	os.Exit(0)
}
