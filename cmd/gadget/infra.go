package main

import (
	"bufio"
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

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

func genGadgetKeys() (string, string, error) {

	randSeedSource := rand.NewSource(time.Now().UnixNano())
	randSeed := rand.New(randSeedSource)

	privateKey, err := rsa.GenerateKey(randSeed, 2014)
	if err != nil {
		panic(err)
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
		panic(err)
	}

	pubString := string(ssh.MarshalAuthorizedKey(pub))

	return privateKeyPem, pubString, err
}

func requiredSsh() error {

	usr, err := user.Current()
	if err != nil {
		panic(err)
	}

	// get proper homedir locations
	sshLocation = fmt.Sprintf("%s/.ssh", usr.HomeDir)
	defaultPrivKeyLocation = fmt.Sprintf("%s/.ssh/gadget_default_rsa", usr.HomeDir)
	gadgetPrivKeyLocation = fmt.Sprintf("%s/.ssh/gadget_rsa", usr.HomeDir)
	gadgetPubKeyLocation = fmt.Sprintf("%s/.ssh/gadget_rsa.pub", usr.HomeDir)

	// check OS for IP address
	if runtime.GOOS == "windows" {
		ip = "192.168.82.1:22"
	}

	// check/create ~/.ssh
	pathExists, err := exists(sshLocation)
	if err != nil {
		panic(err)
	}

	if !pathExists {
		err = os.Mkdir(sshLocation, 0644)
		if err != nil {
			panic(err)
		}
	}

	// check/create ~/.ssh/gadget_default_rsa
	pathExists, err = exists(defaultPrivKeyLocation)
	if err != nil {
		panic(err)
	}

	if !pathExists {
		fmt.Println("[SETUP]  Unable to locate default gadget ssh key, generating..")
		fmt.Printf("[SETUP]    default private key: ~/.ssh/gadget_default_rsa  ")
		outBytes := []byte(defaultKey)
		err = ioutil.WriteFile(defaultPrivKeyLocation, outBytes, 0600)
		if err != nil {
			panic(err)
		}
		fmt.Printf("✔\n")
	}

	// check/create ~/.ssh/gadget_rsa[.pub]
	gadgetPrivExists, err := exists(gadgetPrivKeyLocation)
	if err != nil {
		panic(err)
	}
	gadgetPubExists, err := exists(gadgetPubKeyLocation)
	if err != nil {
		fmt.Println("ERROR: something went wrong with gadgetPubExists `%s`: %s", gadgetPubKeyLocation, err)
		os.Exit(1)
	}

	if !gadgetPrivExists && !gadgetPubExists {
		fmt.Println("[SETUP]  Unable to locate personal gadget ssh keys, generating..")
		
		privkey, pubkey, err := genGadgetKeys()
		if err != nil {
			fmt.Println("ERROR: something went wrong with genGadgetKeys: %s", err)
			os.Exit(1)
		}
		
		fmt.Printf("[SETUP]    private key: ~/.ssh/gadget_rsa  ")
		outBytes := []byte(privkey)
		err = ioutil.WriteFile(gadgetPrivKeyLocation, outBytes, 0600)
		if err != nil {
			fmt.Println("ERROR: something went wrong with gadgetPrivKey `%s`: %s", gadgetPrivKeyLocation, err)
			os.Exit(1)
		}
		fmt.Printf("✔\n")
		
		fmt.Printf("[SETUP]    public key: ~/.ssh/gadget_rsa.pub  ")
		outBytes = []byte(pubkey)
		err = ioutil.WriteFile(gadgetPubKeyLocation, outBytes, 0600)
		if err != nil {
			fmt.Println("ERROR: something went wrong with gadgetPrivKey `%s`: %s", gadgetPubKeyLocation, err)
			os.Exit(1)
		}
		fmt.Printf("✔\n")
	}

	return nil
}

func gadgetLogin(keyLocation string) (*ssh.Client, error) {
	key, err := ioutil.ReadFile(keyLocation)
	if err != nil {
		panic(err)
	}

	// Create the Signer for this private key.
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		panic(err)
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

	//defer client.Close()

	return client, err
}

func gadgetInstallKeys() {
	key, err := ioutil.ReadFile(defaultPrivKeyLocation)
	if err != nil {
		panic(err)
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		panic(err)
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
		panic(err)
	}

	session, err := client.NewSession()
	if err != nil {
		client.Close()
		panic(err)
	}
	
	dest := "/root/.ssh/authorized_keys"
	fmt.Println("[COMMS]  Installing personal gadget ssh key..")
	fmt.Printf("[COMMS]    gadget:%s ", dest)
	err = scp.CopyPath(gadgetPubKeyLocation, dest, session)
	if _, err := os.Stat(gadgetPubKeyLocation); os.IsNotExist(err) {
		fmt.Printf("[COMMS]    ERROR: no such file or directory: %s", gadgetPubKeyLocation)
	} else {
		fmt.Printf("✔\n")
	}

	defer client.Close()
}

func ensureKeys() error {

	_, err := gadgetLogin(gadgetPrivKeyLocation)
	if err != nil {
		fmt.Println("[COMMS]  Private key login failed, trying default key")
		_, err = gadgetLogin(defaultPrivKeyLocation)
		if err != nil {
			fmt.Println("[COMMS]  Default key login also failed, did you leave your keys at home?")
			return err
		} else {
			fmt.Println("[COMMS]  Default key login success")
			gadgetInstallKeys()
			if err != nil {
				panic(err)
			}
		}
	} else {
		fmt.Println("[COMMS]  Private key login success")
	}

	return err
}

func runRemoteCommand(client *ssh.Client, cmd ...string) error {
	session, err := client.NewSession()
	if err != nil {
		client.Close()
		panic(err)
	}

	err = session.Start(strings.Join(cmd[:], " "))
	if err != nil {
		panic(err)
	}
	session.Stdout = os.Stdout
	session.Stderr = os.Stderr
	session.Wait()
	return nil
}

func runLocalCommand(binary string, arguments ...string) {
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
			/*if outScanner.Scan() {
				fmt.Println(string(outScanner.Text()))
			}*/
			_ = outScanner.Scan()
			if errScanner.Scan() {
				fmt.Println(string(errScanner.Text()))
			}
		}
	}()

	execErr = cmd.Start()
	if execErr != nil {
		panic(execErr)
	}
	execErr = cmd.Wait()
	if execErr != nil {
		panic(execErr)
	}
}