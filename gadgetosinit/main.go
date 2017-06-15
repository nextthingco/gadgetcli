package main

import (
	"errors"
	"flag"
	"github.com/nextthingco/libgadget"
	gadgetFormatter "github.com/nextthingco/logrus-gadget-formatter"
	log "github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"strings"
)

type GadgetCommandFunc func([]string, *libgadget.GadgetContext) error

type GadgetCommand struct {
	Name        string
	Function    GadgetCommandFunc
	NeedsConfig bool
}

var Commands = []GadgetCommand{
	{Name: "init", Function: GadgetOsInit, NeedsConfig: true},
	//~ { Name: "add",     Function: GadgetAdd,     NeedsConfig: true  },
	//~ { Name: "build",   Function: GadgetBuild,   NeedsConfig: true  },
	//~ { Name: "deploy",  Function: GadgetDeploy,  NeedsConfig: true  },
	//~ { Name: "start",   Function: GadgetStart,   NeedsConfig: true  },
	//~ { Name: "stop",    Function: GadgetStop,    NeedsConfig: true  },
	//~ { Name: "status",  Function: GadgetStatus,  NeedsConfig: true  },
	//~ { Name: "delete",  Function: GadgetDelete,  NeedsConfig: true  },
	//~ { Name: "shell",   Function: GadgetShell,   NeedsConfig: false },
	//~ { Name: "logs",    Function: GadgetLogs,    NeedsConfig: true  },
	//~ { Name: "run",     Function: GadgetRun,     NeedsConfig: false },
	{Name: "version", Function: GadgetOsVersion, NeedsConfig: false},
	{Name: "help", Function: GadgetOsHelp, NeedsConfig: false},
}

func GadgetOsVersion(args []string, g *libgadget.GadgetContext) error {
	log.Infoln(filepath.Base(os.Args[0]))
	log.Infof("version: %s", libgadget.Version)
	log.Infof("built:   %s", libgadget.BuildDate)
	log.Infof("commit:  %s", libgadget.GitCommit)
	return nil
}

func GadgetOsHelp(args []string, g *libgadget.GadgetContext) error {
	flag.Usage()
	return nil
}

func FindCommand(name string) (*GadgetCommand, error) {
	for _, cmd := range Commands {
		if cmd.Name == name {
			return &cmd, nil
		}
	}
	return nil, errors.New("Failed to find command")
}

func main() {
	// Hey, Listen!
	// Everything that outputs needs to come after g.Verbose check!
	flag.Usage = func() {
		log.Info("")
		log.Infof("USAGE: %s [options] COMMAND", filepath.Base(os.Args[0]))
		log.Info("")
		log.Info("Commands:")
		log.Info("  init        Initialize gadget project")
		//~ log.Info ("  add         Initialize gadget project")
		//~ log.Info ("  build       Build gadget config file")
		//~ log.Info ("  deploy      Build gadget config file")
		//~ log.Info ("  start       Build gadget config file")
		//~ log.Info ("  stop        Build gadget config file")
		//~ log.Info ("  status      Build gadget config file")
		//~ log.Info ("  delete      Build gadget config file")
		//~ log.Info ("  shell       Connect to remote device running GadgetOS")
		//~ log.Info ("  logs        Build gadget config file")
		log.Info("  version     Print version information")
		log.Info("  help        Print this message")
		log.Info("")
		log.Infof("Run '%s COMMAND --help' for more information on the command", filepath.Base(os.Args[0]))
		//~ log.Info ("")
		//~ log.Infof("Options:")
		//~ log.Info ("  -C <path>                            ")
		//~ log.Info ("    	Run in directory (default \".\")  ")
		//~ log.Info ("  -v	Verbose execution                 ")
		log.Info("")
	}

	g := libgadget.GadgetContext{}

	g.WorkingDirectory = "/etc/"

	//~ flag.BoolVar(&g.Verbose, "v", false, "Verbose execution")
	//~ flag.StringVar(&g.WorkingDirectory, "C", ".", "Run in directory")
	flag.Parse()

	var gFormatter *gadgetFormatter.TextFormatter

	//~ if g.Verbose {
	gFormatter = new(gadgetFormatter.TextFormatter)
	gFormatter.DisableColors = true

	log.SetLevel(log.DebugLevel)
	//~ } else {
	//~ gFormatter = new(gadgetFormatter.TextFormatter)
	//~ gFormatter.DisableColors = true
	//~ gFormatter.DisableTimestamp = true
	//~ gFormatter.DisableSorting = true
	//~ gFormatter.EntryString.InfoLevelString = "I:"
	//~ gFormatter.EntryString.WarnLevelString = "W:"
	//~ gFormatter.EntryString.ErrorLevelString = "E:"

	//~ log.SetLevel(log.InfoLevel)
	//~ }

	log.SetFormatter(gFormatter)

	// Hey, Listen!
	// Everything that outputs needs to come after g.Verbose check!

	//~ err := libgadget.RequiredSsh()
	//~ if err != nil {
	//~ log.Error("Failed to verify ssh requirements")
	//~ os.Exit(1)
	//~ }

	args := flag.Args()
	if len(args) < 1 {
		flag.Usage()
		log.Error("No Command Specified")
		os.Exit(1)
	}

	// file command
	cmd, err := FindCommand(args[0])
	if err != nil {
		flag.Usage()
		log.WithFields(log.Fields{
			"command": strings.Join(args[0:], " "),
		}).Debug("Command is not valid")
		log.Errorf("Command %s is not valid", args[0:])
		os.Exit(1)
	}

	// if command needs to use the config file, load it
	if cmd.NeedsConfig {
		err = g.LoadConfig()
		if err != nil {
			log.Error("Failed to load config")
			log.Warn("Be sure to run gadget in the same directory as 'gadget.yml'")
			log.Warn("Or specify a directory e.g. 'gadget -C ../projects/gpio/ [command]'")
			os.Exit(1)
		}
	}

	err = cmd.Function(args[1:], &g)
	if err != nil {
		os.Exit(1)
	}
	os.Exit(0)
}
