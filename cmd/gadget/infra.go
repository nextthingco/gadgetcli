
package main

import (
    "runtime"
	"os"
	"os/user"
	"fmt"
	"time"
	"math/rand"
	"io/ioutil"
	"crypto/rsa"
	"crypto/x509"
	"golang.org/x/crypto/ssh"
	"encoding/pem"
	//~ "encoding/base64"
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
	
	sshLocation = ""
	defaultPrivKeyLocation = ""
	gadgetPrivKeyLocation = ""
	gadgetPubKeyLocation = ""
)

func exists(path string) (bool, error) {
    _, err := os.Stat(path)
    if err == nil { return true, nil }
    if os.IsNotExist(err) { return false, nil }
    return true, err
}

func genGadgetKeys () (string, string, error) {
	
	randSeedSource := rand.NewSource(time.Now().UnixNano())
	randSeed := rand.New(randSeedSource)
	
	privateKey, err := rsa.GenerateKey(randSeed, 2014)
	if err != nil {
		fmt.Println("ERROR: something went wrong with rsa.GenerateKey: %s", err)
		os.Exit(1)
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
		fmt.Println("ERROR: something went wrong with ssh.NewPublicKey: %s", err)
		os.Exit(1)
	}
	pubString := string(ssh.MarshalAuthorizedKey(pub))
	
	fmt.Println(privateKeyPem)
	fmt.Println(pubString)
	
	
	return privateKeyPem, pubString, err
}

func requiredSsh () error {
	
	usr, err := user.Current()
	if err != nil {
		fmt.Println("ERROR: something went wrong with user.Current: %s", err)
		os.Exit(1)		
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
		fmt.Println("ERROR: something went wrong with pathExists `%s`: %s", sshLocation, err)
		os.Exit(1)
	}
	
	if !pathExists {
		err = os.Mkdir(sshLocation, 0644)
	}
	
	// check/create ~/.ssh/gadget_default_rsa
	pathExists, err = exists(defaultPrivKeyLocation)
	if err != nil {
		fmt.Println("ERROR: something went wrong with pathExists `%s`: %s", defaultPrivKeyLocation, err)
		os.Exit(1)
	}
	
	if !pathExists {
		outBytes := []byte(defaultKey)
		err = ioutil.WriteFile(defaultPrivKeyLocation, outBytes, 0600)
		if err != nil {
			fmt.Println("ERROR: something went wrong with defaultKey `%s`: %s", defaultPrivKeyLocation, err)
			os.Exit(1)
		}
	}
	
	// check/create ~/.ssh/gadget_rsa[.pub]
	gadgetPrivExists, err := exists(gadgetPrivKeyLocation)
	if err != nil {
		fmt.Println("ERROR: something went wrong with gadgetPrivExists `%s`: %s", gadgetPrivKeyLocation, err)
		os.Exit(1)
	}
	gadgetPubExists, err := exists(gadgetPubKeyLocation)
	if err != nil {
		fmt.Println("ERROR: something went wrong with gadgetPubExists `%s`: %s", gadgetPubKeyLocation, err)
		os.Exit(1)
	}
	
	var privkey, pubkey string = "", ""
	_ = privkey
	_ = pubkey
	// ^gross
	if !gadgetPrivExists && !gadgetPubExists {
		privkey, pubkey, err = genGadgetKeys()
		
		outBytes := []byte(privkey)
		err = ioutil.WriteFile(gadgetPrivKeyLocation, outBytes, 0600)
		if err != nil {
			fmt.Println("ERROR: something went wrong with gadgetPrivKey `%s`: %s", gadgetPrivKeyLocation, err)
			os.Exit(1)
		}
		
		outBytes = []byte(pubkey)
		err = ioutil.WriteFile(gadgetPubKeyLocation, outBytes, 0600)
		if err != nil {
			fmt.Println("ERROR: something went wrong with gadgetPrivKey `%s`: %s", gadgetPubKeyLocation, err)
			os.Exit(1)
		}
	}
	
	return nil
}
