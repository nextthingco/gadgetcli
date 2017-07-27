/*
This file is part of the Gadget command-line tools.
Copyright (C) 2017 Next Thing Co.

Gadget is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 2 of the License, or
(at your option) any later version.

Gadget is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with Gadget.  If not, see <http://www.gnu.org/licenses/>.
*/

package main

import (
	"errors"
	//~ "fmt"
	"github.com/nextthingco/libgadget"
	//~ "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
	//~ "gopkg.in/yaml.v2"
	//~ "io/ioutil"
	//~ "github.com/kr/pty"
	//~ "bufio"
	"os/exec"
	//~ "io"
	//~ "bytes"
	//~ "strings"
	//~ "os"
)

var (
)

func editUsage() error {
	log.Info("Usage:  gadget [flags] edit [type] [value]     ")
	log.Info("                *opt         *req   *opt       ")
	log.Info("Type:           service | onboot | rootfs      ")
	log.Info("Value (containers): not yet implemented        ")
	log.Info("Value (rootfs): kernel <more to be added soon> ")

	return errors.New("Incorrect edit usage")
}

func GadgetEditKernel(g *libgadget.GadgetContext) error {
	
	log.Info("Edit Kernel")
	
	//~ c := exec.Command("docker", "run", "-it", "--rm", "c98", "make", "")
	//~ f, err := pty.Start(c)
	//~ if err != nil {
		//~ panic(err)
	//~ }

	//~ go func() {
		//~ f.Write([]byte("foo\n"))
		//~ f.Write([]byte("bar\n"))
		//~ f.Write([]byte("baz\n"))
		//~ f.Write([]byte{4}) // EOT
	//~ }()
	//~ io.Copy(os.Stdout, f)
	
	//~ log.Debugf("Calling %s %s", binary, arguments)

	//~ cmd := exec.Command("docker", "run", "-it", "--rm", "c98", "make", "menuconfig")
	//~ f, err := pty.Start(c)
	//~ cmd.Env = os.Environ()

	//~ os.Stdout = cmd.StdoutPipe()
	//~ os.Stderr = cmd.StderrPipe()
	//~ cmd.StdinPipe() = os.Stdin
	
	

	//~ stdOutReader, execErr := cmd.StdoutPipe()
	//~ if execErr != nil {
		//~ log.Debugf("Couldn't connect to cmd.StdoutPipe()")
	//~ }
	
	//~ stdErrReader, execErr := cmd.StderrPipe()
	//~ if execErr != nil {
		//~ log.Debugf("Couldn't connect to cmd.StderrPipe()")
	//~ }
	
	//~ stdInWriter, execErr := cmd.StdinPipe()
	//~ if execErr != nil {
		//~ log.Debugf("Couldn't connect to cmd.StderrPipe()")
	//~ }
	
	//~ outScanner := bufio.NewScanner(stdOutReader)
	//~ errScanner := bufio.NewScanner(stdErrReader)
	//~ inScanner := bufio.NewScanner(os.Stdin)
	
	//~ var outBuffer bytes.Buffer
	//~ var errBuffer bytes.Buffer
	
	//~ // goroutines to print stdout and stderr [doesn't quite work]
	//~ go func(){
		//~ if g.Verbose {
			//~ for outScanner.Scan(){
				//~ log.Debugf(string(outScanner.Text()))
				//~ outBuffer.WriteString(string(outScanner.Text()))
			//~ }
		//~ } else {
			//~ for outScanner.Scan(){
				//~ if strings.Contains(outScanner.Text(), "Step "){
					//~ log.Infof("    %s",string(outScanner.Text()))
				//~ }
				//~ outBuffer.WriteString(string(outScanner.Text()))
			//~ }
		//~ }
	//~ }()
	
	//~ printedStderr := false
	
	//~ go func(){
		//~ for errScanner.Scan(){
			//~ log.Warnf(string(errScanner.Text()))
			//~ printedStderr = true
			//~ errBuffer.WriteString(string(errScanner.Text()))
		//~ }
	//~ }()
	
	//~ go func(){
		//~ for inScanner.Scan(){
			
			//~ log.Warnf(string(errScanner.Bytes()))
			//~ errBuffer.WriteString(string(errScanner.Text()))
		//~ }
	//~ }()
	
	return nil
}

// Process the build arguments and execute build
func GadgetEdit(args []string, g *libgadget.GadgetContext) error {
	
	log.Info("Edit")
	log.Debugf("args %s", args)
	
	// find docker binary in path
	binary, err := exec.LookPath("docker")
	if err != nil {
		log.Error("Failed to find local docker binary")
		log.Warn("Is docker installed?")

		log.WithFields(log.Fields{
			"function": "GadgetBuild",
			"stage":    "LookPath(docker)",
		}).Debug("Couldn't find docker in the $PATH")
		return err
	}
	
	err = libgadget.EnsureDocker(binary, g)
	if err != nil {
		log.Errorf("Failed to contact the docker daemon.")
		log.Warnf("Is it installed and running with appropriate permissions?")
		return err
	}
	
	if len(args) != 2 {
		log.Error("Invalid arguments for `gadget edit`")
		return editUsage();
	}
	
	// parse arguments
	switch args[0] {
	case "rootfs":
		// parse edit rootfs argument
		switch args[1] {
		case "kernel":
			err = GadgetEditKernel(g)
			if err != nil {
				log.Errorf("Failed to edit the kernel config.")
				return err
			}
		default:
			log.Errorf("  %q is not valid argument or is not yet supported.", args[1])
			return editUsage()
		}
	default:
		log.Errorf("  %q is not valid argument or is not yet supported.", args[0])
		return editUsage()
	}

	return nil
}
