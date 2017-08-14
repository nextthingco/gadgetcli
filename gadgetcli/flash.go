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
	"github.com/nextthingco/libgadget"
	"golang.org/x/crypto/ssh"
	//~ "gopkg.in/cheggaaa/pb.v1"
	log "gopkg.in/sirupsen/logrus.v1"
	"io"
	//~ "io/ioutil"
	//~ "bufio"
	"os"
	"fmt"
	"crypto/sha256"
	//~ "os/exec"
	//~ "strings"
)

type ArtifactDef struct {
	Board        string
	Artifacts    []string
	ArtifactType []string
}

var (
	ArtDefs = []ArtifactDef{
		ArtifactDef { 
			Board:        "chippro",
			Artifacts:    []string {"zImage", "ntc-gr8-crumb.dtb", "rootfs.ubifs"},
			ArtifactType: []string {"kernel", "fdt", "rootfs"},
		},
		//~ ArtifactDef { 
			//~ Board: "chippro4gb",
			//~ Artifacts: []string {"zImage", "ntc-gr8-crumb.dtb", "rootfs.ubifs"},
		//~ },
		//~ ArtifactDef {
			//~ Board: "chip",
			//~ Artifacts: []string {"zImage", "ntc-r8-chip.dtb", "rootfs.ubifs"},
		//~ },
	}
)

func GadgetFlashFile(client *ssh.Client, artifactLocation string, artifactType string, g *libgadget.GadgetContext) error {
	
	log.Debugf("artLoc: %s", artifactLocation)
	log.Debugf("artTyp: %s", artifactType)
	
	// open the file
	checkFile, err := os.Open(artifactLocation)
	if err != nil || checkFile == nil {
		log.Errorf("Failed to open file %s", artifactLocation)
		return err
	}
	defer checkFile.Close()
	//~ artReader := bufio.NewReader(artFile)
	
	
	session, err := client.NewSession()
	if err != nil {
		client.Close()
		return err
	}
	
	// get file size
	fi, err := checkFile.Stat()
	if err != nil {
		log.Errorf("Failed to stat file")
		return err
	}
	size := int(fi.Size())

	// get hash
	checksum := sha256.New()
	if _, err := io.Copy(checksum, checkFile); err != nil {
		log.Error("Failed to get checksum")
		return err
	}
	log.Debugf("checksum: %x", checksum.Sum(nil))

	// create pipe for local -> remote file transmission
	//~ pr, pw := io.Pipe()
	sessionLogger := log.New()
	if g.Verbose {
		sessionLogger.Level = log.DebugLevel
	}
	
	stdin, err := session.StdinPipe()
	if err != nil {
		return fmt.Errorf("Unable to setup stdin for session: %v", err)
	}
	
	// open the file
	artFile, err := os.Open(artifactLocation)
	if err != nil || checkFile == nil {
		log.Errorf("Failed to open file %s", artifactLocation)
		return err
	}
	defer artFile.Close()
	
	go func(){
		if _, err := io.Copy(stdin, artFile); err != nil {
			log.Errorf("Failed to copy file to stdin")
			return
		}
		stdin.Close()
	}()
	
	session.Stdout = sessionLogger.WriterLevel(log.DebugLevel)
	session.Stderr = sessionLogger.WriterLevel(log.DebugLevel)

	log.Debug("  Starting session")
	
	//~ bar.Start()
	
	//~ contents_bytes, err := ioutil.ReadAll(artFile)
	//~ if err != nil {
		//~ log.Errorf("Failed to read file")
		//~ return err
	//~ }
	
	// set reader command
	sessionCmd := fmt.Sprintf("update_volume %s %d %x", artifactType, size, checksum.Sum(nil))
	//~ sessionCmd = fmt.Sprintf("cat > %s", artifactType)
	//~ sessionCmd = "md5sum"
	//~ sessionCmd := fmt.Sprintf("/bin/sh -x /sbin/dumbcat %s %d %x %s", artifactType, size, checksum.Sum(nil), contents_bytes)
	log.Debugf("sessionCmd: %s", sessionCmd)
	
	
	
	//~ go func() {
		//~ w, _ := session.StdinPipe()
		//~ defer w.Close()
		//~ fmt.Fprintln(w, contents_bytes)
	//~ }()
	
	
	if err := session.Start(sessionCmd); err != nil {
		return err
	}
	
	
	
	session.Wait()
	//~ bar.Finish()

	session.Close()

	return nil
}

// Process the build arguments and execute build
func GadgetFlash(args []string, g *libgadget.GadgetContext) error {

	err := libgadget.EnsureKeys()
	if err != nil {
		log.Errorf("Failed to connect to Gadget")
		return err
	}

	client, err := libgadget.GadgetLogin(libgadget.GadgetPrivKeyLocation)
	if err != nil {
		log.Errorf("Failed to connect to Gadget")
		return err
	}
	
	// check for non-empty board definition
	board := g.Config.Rootfs.From
	image := g.Config.Rootfs.Hash
	
	if board == "" && image == "" {
		log.Errorf("Failed to find rootfs")
		log.Errorf("One or more [mis/un]configured entries:")
		log.Errorf("From: %s", board)
		log.Errorf("Hash: %s", image)
		return errors.New("Failed to flash rootfs")
	}
	
	// check that board is supported
	matchedBoard := ArtifactDef { Board: "", }
	
	for _, def := range ArtDefs {
		if board == def.Board {
			matchedBoard = def
			log.Infof("  Flashing: %s", board)
			break
		}
	}
	
	if matchedBoard.Board == "" {
		log.Errorf("%s is not a valid From:", board)
		return errors.New("Invalid board definition")
	}
	
	// test to make sure all payload files present
	for _, payloadPart := range matchedBoard.Artifacts {
		partLocation := g.WorkingDirectory + "/.images/" + payloadPart
		partExists, err := libgadget.PathExists(partLocation)
		if !partExists {
			log.Errorf("Could not locate '%s'", partLocation)
			return errors.New("Failed to locate linux config")
		}
		if err != nil {
			log.Errorf("Failed to determine if '%s' exists", partLocation)
			return err
		}
	}
	
	// flash each part
	for i, flashPart := range matchedBoard.Artifacts {
		
		partLocation := g.WorkingDirectory + "/.images/" + flashPart
		partType := matchedBoard.ArtifactType[i]
		
		err = GadgetFlashFile(client, partLocation, partType, g)
		if err != nil {
			return err
		}
	}
	
	return err
}
