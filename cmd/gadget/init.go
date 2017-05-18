package main

import (
	"fmt"
	"github.com/satori/go.uuid"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
)

var (
	dockerFileContents = "FROM armhf/alpine\n\nADD blink-leds /"
	
	blinkLedsContents = `
#!/bin/sh
MODE="ascend"
SPEED=0.125
for i in $(seq 132 139)
do
	echo $i > /sys/class/gpio/export &
done

sleep $SPEED

for i in $(seq 132 139)
do
	echo out > /sys/class/gpio/gpio$i/direction &
done

if [ $MODE == "all" ]; then
	while true; do
		for i in $(seq 132 139)
		do
			echo 1 > /sys/class/gpio/gpio$i/value &
		done
		sleep $SPEED
		for i in $(seq 132 139)
		do
			echo 0 > /sys/class/gpio/gpio$i/value &
		done
		sleep $SPEED
	done
else
	while true; do
                for i in $(seq 132 139)
                do
                        echo 1 > /sys/class/gpio/gpio$i/value &
			sleep $SPEED
                done
                sleep $SPEED
                for i in $(seq 132 139)
                do
                        echo 0 > /sys/class/gpio/gpio$i/value &
			sleep $SPEED
                done
                sleep $SPEED
        done
fi
`
)

// Process the build arguments and execute build
func gadgetInit(args []string, g *GadgetContext) {

	initUu1 := uuid.NewV4()
	initUu2 := uuid.NewV4()
	initUu3 := uuid.NewV4()

	g.WorkingDirectory, _ = filepath.Abs(g.WorkingDirectory)
	initName := filepath.Base(g.WorkingDirectory)

	initConfig := templateConfig(initName, fmt.Sprintf("%s", initUu1), fmt.Sprintf("%s", initUu2), fmt.Sprintf("%s", initUu3))

	outBytes, err := yaml.Marshal(initConfig)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = ioutil.WriteFile(fmt.Sprintf("%s/gadget.yml", g.WorkingDirectory), outBytes, 0644)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	
	containerDir := fmt.Sprintf("%s/%s", g.WorkingDirectory, initConfig.Services[0].Directory)
	if _, err := os.Stat(containerDir); os.IsNotExist(err) {
		os.Mkdir(containerDir, 0755)
	} else {
		panic(err)
	}
	
	err = ioutil.WriteFile(
		fmt.Sprintf("%s/Dockerfile", containerDir), 
		[]byte(dockerFileContents), 0644)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	
	err = ioutil.WriteFile(
		fmt.Sprintf("%s/blink-leds", containerDir), 
		[]byte(blinkLedsContents), 0755)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	
}
