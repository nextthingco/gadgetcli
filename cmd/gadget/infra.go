package main

import (
	"bufio"
	"errors"
	"bytes"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/tmc/scp"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"math/rand"
	"os"
	"os/user"
	"os/exec"
	"runtime"
	"time"
	"strings"
	"path/filepath"
	log "github.com/sirupsen/logrus"
)

var (
	defaultKey = `
-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEApDqwmXssH0AIk6pec02wSYDrORY9J7VlE/45IqAACzTQuD7X
CifmEZ6lzKCqvGSSXBFpfm/nh8W0vcTn20yfehkz1xI9yP+sJSUhNavkW+5hl9Yo
uXFK91sbRKuvAUqu+ozV5wD3ZP4tFtoH39SxzldBJaiW9NQ1XdCgGy7RulREuMBP
3tFUHP1o1btiql5YO3YoHrizAwcIbfoyCgnwVlMldq249AGBO5jvEoRnhO4Toikh
TZik5isKsVx1hfMp6fkk5pdCnRWXCTNOXo+DbsUAzXKoFLoA9c5reIeOjeZG6+3L
T/G+QBEwjimTRXoPFmPYcJ/KMurp2P/o1+QVbQIDAQABAoIBAC5Pwpc1bcbONtz1
UTcwtEK2ER8DD3HQLFXL/e6usfR3C1i5l8hsYeucEmM2946yybcezeHyypa2APb2
vO9RlzNGQiEnKrcwqim7Y7cP5xCpk2nO4aMRuLMyROlDhNFXbyqGZpeC5UDckHh+
OXQ8NXvbjSqCdTdLVFVFTLD9rfTd9L85lSxAD6jUhMbiu62b+9j14xWRAvUkt/sd
1cr/5yh7x0MYcJkHfShMNai/ExU8cUSwNF/JfxhmIt1aO7sefshLhxJUuAPT1ljd
mnD0ZeGWUt52yqrjaWLMKgtTtJuSUiRsTQ9Qx4eIzP6PJWJIL6+j63M3VL25yRz5
D8SOpEECgYEA2Xc0Gtl3RPPv+AakGEJb/TdW/mz5Fa6yTVjjimB/UcVrwOFT6y5h
kpoqSH/SzBZxdkdNli6Cgyr0ajNb6oXL/rf/0R7Do+VNQYxZIvrT95ANxCH6HEy4
4UMLGT8Xz5gOEBKH27+j/sdUqqT/dsv8iLRXVk6yGQZXHfYjXBcLOPECgYEAwVSM
U4aE4JBDAxTGDejaGFjNlmKcPjSEa/Yv5Qd5/qcGyQM6wHCV8TXE1ry9DSONviCI
qmnR9BqjFb4/I6jI7zuHDneAXUJja1Kap01rTbrWCaRJTwRGDELNj1/aDC6HOaaf
zZN21dafAsg8d0vv1SXzQThJ63LwQt1qTTSPxD0CgYEAkwerbPPXVgFwH8utqtFD
DMMbyE25Y1WILA+LWIXBz3GhVvmCGaJ0SgB90iLKTT5nXEb9SCsOBs1GD3/GB5yK
vh99kNAyCmAAie7wXVwlcF4vUIqAZh3hajxABsPHv43ZBDjjLko2AQ6YSf/g0Vs9
1NfJrQrsE0tcH1/JrHvQFKECgYBEjn7Uf7dPCtk4ln1FIXV1fMgqs/1D8cujnUGO
rgAM1Z4KWiLTaxlA2BhdLcC8kAcLjO3pwGy7a1a5tyUcuBXJAAr8jlPuvkQTIs/E
1CdhAQg1kxSL+K/+WRIb7ZmdCELbpsK0W76gReNNUURf6YW6yCJi1lsgKzoX+/xe
NG1m4QKBgG2UOBk+9hF3bcq0Wo4zSDa3wTzPlTlySnUOU1m6pMlUW97qhmUzKdj/
EGjLLdEY/nQkBYT5HmV4lilHlrb+fZcM0+FegopkKXAOzEqkLTI2ibiItCT12nLB
FwRYLLbqbGByhykSn5ybp/DuSQpH4blitu/fEYOg6QX/I/6zayd+
-----END RSA PRIVATE KEY-----
`

	ip = "192.168.81.1:22"

	sshLocation            = ""
	defaultPrivKeyLocation = ""
	gadgetPrivKeyLocation  = ""
	gadgetPubKeyLocation   = ""
)

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func GenGadgetKeys() (string, string, error) {

	randSeedSource := rand.NewSource(time.Now().UnixNano())
	randSeed := rand.New(randSeedSource)

	privateKey, err := rsa.GenerateKey(randSeed, 2048)
	if err != nil {
		return "", "", err
	}

	privateKeyDer := x509.MarshalPKCS1PrivateKey(privateKey)
	privateKeyBlock := pem.Block{
		Type:    "RSA PRIVATE KEY",
		Headers: nil,
		Bytes:   privateKeyDer,
	}
	privateKeyPem := string(pem.EncodeToMemory(&privateKeyBlock))

	publicKey := privateKey.PublicKey

	pub, err := ssh.NewPublicKey(&publicKey)
	if err != nil {
		return privateKeyPem, "", err
	}

	pubString := string(ssh.MarshalAuthorizedKey(pub))

	return privateKeyPem, pubString, err
}

func RequiredSsh() error {

	usr, err := user.Current()
	if err != nil {
		return err
	}

	// get proper homedir locations
	sshLocation = filepath.Join(usr.HomeDir,".ssh")
	defaultPrivKeyLocation = filepath.Join(sshLocation,"gadget_default_rsa")
	gadgetPrivKeyLocation = filepath.Join(sshLocation,"gadget_rsa")
	gadgetPubKeyLocation = filepath.Join(sshLocation,"gadget_rsa.pub")

	// check OS for IP address
	if runtime.GOOS == "windows" {
		ip = "192.168.82.1:22"
	}

	// check/create ~/.ssh
	sshDirExists, err := PathExists(sshLocation)
	if err != nil {
		return err
	}

	if !sshDirExists {
		err = os.Mkdir(sshLocation, 0644)
		if err != nil {
			return err
		}
	}

	// check/create ~/.ssh/gadget_default_rsa
	defaultPrivExists, err := PathExists(defaultPrivKeyLocation)
	if err != nil {
		return err
	}

	if !defaultPrivExists {
		log.Warn("Unable to locate default gadget ssh key, generating..")

		log.WithFields(log.Fields{
			"function": "RequiredSsh",
		}).Debug("default private key: ~/.ssh/gadget_default_rsa")
		
		outBytes := []byte(defaultKey)
		err = ioutil.WriteFile(defaultPrivKeyLocation, outBytes, 0600)
		if err != nil {
			return err
		}
	}

	// check/create ~/.ssh/gadget_rsa[.pub]
	gadgetPrivExists, err := PathExists(gadgetPrivKeyLocation)
	if err != nil {
		log.WithFields(log.Fields{
			"function": "RequiredSsh",
			"keyLocation": gadgetPrivKeyLocation,
			"error": err,
		}).Error("something went wrong with gadgetPubExists")
		return err
	}
	gadgetPubExists, err := PathExists(gadgetPubKeyLocation)
	if err != nil {
		log.WithFields(log.Fields{
			"function": "RequiredSsh",
			"keyLocation": gadgetPubKeyLocation,
			"error": err,
		}).Error("something went wrong with gadgetPubExists")
		return err
	}

	if !gadgetPrivExists && !gadgetPubExists {
		log.Warn("Unable to locate personal gadget ssh keys, generating..")
		
		log.WithFields(log.Fields{
			"function": "RequiredSsh",
			"keyLocation": gadgetPubKeyLocation,
			"error": err,
		}).Debug("private key: ~/.ssh/gadget_rsa[.pub]")
		
		privkey, pubkey, err := GenGadgetKeys()
		if err != nil {
			log.WithFields(log.Fields{
				"function": "RequiredSsh",
				"error": err,
			}).Error("something went wrong with genGadgetKeys")
			return err
		}
		
		log.Info("    private key: ~/.ssh/gadget_rsa")
		outBytes := []byte(privkey)
		err = ioutil.WriteFile(gadgetPrivKeyLocation, outBytes, 0600)
		if err != nil {
			log.WithFields(log.Fields{
				"function": "RequiredSsh",
				"keyLocation": gadgetPrivKeyLocation,
				"error": err,
			}).Error("something went wrong with gadgetPrivKey")
			return err
		}
		
		log.Info("    public key: ~/.ssh/gadget_rsa.pub")
		outBytes = []byte(pubkey)
		err = ioutil.WriteFile(gadgetPubKeyLocation, outBytes, 0600)
		if err != nil {
			log.WithFields(log.Fields{
				"function": "RequiredSsh",
				"keyLocation": gadgetPrivKeyLocation,
				"error": err,
			}).Error("something went wrong with gadgetPrivKey")
			return err
		}
	}

	return nil
}

func GadgetLogin(keyLocation string) (*ssh.Client, error) {
	key, err := ioutil.ReadFile(keyLocation)
	if err != nil {
		log.WithFields(log.Fields{
			"function": "RequiredSsh",
			"keyLocation": gadgetPrivKeyLocation,
			"error": err,
		}).Error("something went wrong with gadgetPrivKey")
		return nil, err
	}

	// Create the Signer for this private key.
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return nil, err
	}

	config := &ssh.ClientConfig{
		User: "root",
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout: (time.Second * 3),
	}

	// Connect to the remote server and perform the SSH handshake.
	client, err := ssh.Dial("tcp", ip, config)
	if err != nil {
		return nil, err
	}

	return client, err
}

func GadgetInstallKeys() error {
	key, err := ioutil.ReadFile(defaultPrivKeyLocation)
	if err != nil {
		log.WithFields(log.Fields{
			"function": "GadgetInstallKeys",
			"file": defaultPrivKeyLocation,
		}).Error("Failed to read private key")
		return err
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		log.WithFields(log.Fields{
			"function": "GadgetInstallKeys",
			"file": defaultPrivKeyLocation,
		}).Error("Failed parse private key")
		return err
	}

	config := &ssh.ClientConfig{
		User: "root",
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout: (time.Second * 3),
	}

	client, err := ssh.Dial("tcp", ip, config)
	if err != nil {
		log.WithFields(log.Fields{
			"function": "GadgetInstallKeys",
			"tcp": ip,
		}).Error("Failed to dial ssh")
		return err
	}

	session, err := client.NewSession()
	if err != nil {
		log.WithFields(log.Fields{
			"function": "GadgetInstallKeys",
		}).Error("Failed to create new ssh client session")
		client.Close()
		return err
	}
	
	dest := "/root/.ssh/authorized_keys"
	log.WithFields(log.Fields{
		"function": "RequiredSsh",
		"gadget": dest,
	}).Debug("Installing personal gadget ssh key..")

	err = scp.CopyPath(gadgetPubKeyLocation, dest, session)
	if _, err := os.Stat(gadgetPubKeyLocation); os.IsNotExist(err) {
		log.WithFields(log.Fields{
			"function": "RequiredSsh",
			"gadget": dest,
			"gadgetPubKeyLocation": gadgetPubKeyLocation,
		}).Error("Public key file does not exist")
	}
	
	defer client.Close()
	return nil
}

func EnsureKeys() error {

	_, err := GadgetLogin(gadgetPrivKeyLocation)
	if err != nil {
		log.Warn("  Private key login failed, trying default key")
		
		log.WithFields(log.Fields{
			"function": "EnsureKeys",
			"gadgetPrivKeyLocation": gadgetPrivKeyLocation,
			"defaultPrivKeyLocation": defaultPrivKeyLocation,
		}).Debug("  Private key login failed, trying default key")

		_, err = GadgetLogin(defaultPrivKeyLocation)
		if err != nil {
			log.Error("  Default key login also failed")
			log.Warn("  Is the gadget connected and powered on?")
			log.Warn("  Was the gadget first used on another computer/account?")
			return err
		} else {
			log.WithFields(log.Fields{
				"function": "EnsureKeys",
			}).Debug("  Default key login success")

			log.WithFields(log.Fields{
				"function": "EnsureKeys",
				"gadgetPubKeyLocation": gadgetPubKeyLocation,
			}).Debug("  Public key file does not exist")

			err = GadgetInstallKeys()
			if err != nil {
				return err
			}
		}
	} else {
		log.WithFields(log.Fields{
			"function": "EnsureKeys",
		}).Debug("  Private key login success")
	}

	return err
}

func RunRemoteCommand(client *ssh.Client, cmd ...string) (*bytes.Buffer, *bytes.Buffer, error) {
	session, err := client.NewSession()
	if err != nil {
		client.Close()
		return nil, nil, err
	}

	err = session.Start(strings.Join(cmd[:], " "))
	if err != nil {
		return nil, nil, err
	}
	var outBuffer bytes.Buffer
	var errBuffer bytes.Buffer
	session.Stdout = &outBuffer
	session.Stderr = &errBuffer
	err = session.Wait()
	
	return &outBuffer, &errBuffer, err
}

func RunLocalCommand(binary string, g *GadgetContext, arguments ...string) (string, string, error) {
	cmd := exec.Command(binary, arguments...)
	
	cmd.Env = os.Environ()

	stdOutReader, execErr := cmd.StdoutPipe()
	stdErrReader, execErr := cmd.StderrPipe()
	outScanner := bufio.NewScanner(stdOutReader)
	errScanner := bufio.NewScanner(stdErrReader)

	// goroutine to print stdout and stderr
	go func() {
		// TODO: goroutine gets launched and never exits.
		for {
			// TODO: add a check here to only print stdout if verbose
			if g.Verbose && outScanner.Scan() {
				log.Debugf(string(outScanner.Text()))
			}
			_ = outScanner.Scan()
			if errScanner.Scan() {
				log.Errorf(string(errScanner.Text()))
			}
		}
	}()

	execErr = cmd.Start()
	if execErr != nil {
		return outScanner.Text(), errScanner.Text(), execErr
	}

	execErr = cmd.Wait()
	
	return outScanner.Text(), errScanner.Text(), execErr
}


func PrependToStrings(stringArray []string, prefix string) []string {
	for key,value := range stringArray {
		s := []string{prefix, value}
		stringArray[key] = strings.Join(s,"")
	}
	return stringArray //strings.Join(stringArray, " ")
}

func FindStagedContainers(args []string, containers GadgetContainers) (GadgetContainers, error) {
	var stagedContainers GadgetContainers
	var unavailableContainers []string

	// if we have any arguments, we're specifying containers to build
	if len(args) > 0 {
		for _,arg := range args {
			c,err := containers.Find(arg)
			if err != nil {
				unavailableContainers = append(unavailableContainers, arg)
			} else {
				stagedContainers = append(stagedContainers, c)
			}
		}
	}

	if len(stagedContainers) == 0 {
		stagedContainers = containers
	}
	err := errors.New(fmt.Sprintf("  Could not find containers: %s", strings.Join(unavailableContainers, ", ")))
	return stagedContainers, err
}
