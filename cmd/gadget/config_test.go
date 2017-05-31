
package main

import (
	"fmt"
	"testing"
	"reflect"
	"github.com/satori/go.uuid"
)

func TestTemplateConfig(t *testing.T) {
	
	fmt.Println("TestTemplateConfig")
	
	initUu1 := uuid.NewV4()
	initUu2 := uuid.NewV4()

	// TEST 0
	
	expectedConfig := GadgetConfig{
		Spec: Version,
		Name: "thisisthename",
		UUID: fmt.Sprintf("%s", initUu1),
		Type: "docker",
		Onboot: []GadgetContainer{
			{
				Name:    "hello-world",
				Image:   "armhf/hello-world",
				UUID:    fmt.Sprintf("%s", initUu2),
			},
		},
	}
		
	testConfig := TemplateConfig("thisisthename", fmt.Sprintf("%s", initUu1), fmt.Sprintf("%s", initUu2))
	
	if ! reflect.DeepEqual(testConfig, expectedConfig) {
		t.Error("isn't deeply equal, but should have been")
		fmt.Println("%+v", expectedConfig)
		fmt.Println("%+v", testConfig)
	}
	
}


func TestParseConfig(t *testing.T) {
	
	in := []byte(`spec: unknown
name: tmp
uuid: ee0212c1-6880-433e-994b-795dfc71000
type: docker
onboot:
- name: hello-world
  uuid: 1e885632-c518-475b-b405-31ca23bd001
  image: armhf/hello-world
  directory: ""
  net: ""
  pid: ""
  readonly: false
  command: []
  binds: []
  capabilities: []
services:
- name: bello-world
  uuid: 1e885632-c518-475b-b405-31ca23bd5002
  image: armhf/bello-world
  directory: ""
  net: ""
  pid: ""
  readonly: false
  command: []
  binds: []
  capabilities: []`)
	
	config, err := ParseConfig(in)
	if err != nil {
		t.Error("Parse Config failed")
	}
	
	if config.Name != "tmp" {
		t.Error("Name should have been tmp")
	}
	
	if config.Services[0].Image != "armhf/bello-world" {
		t.Error("Services[0].Image should have been armhf/bello-world")
	}
	
}

func TestCleanConfig(t *testing.T){
	
	fmt.Println("TestCleanConfig")
	
	initUu1 := uuid.NewV4()
	initUu2 := uuid.NewV4()

	// TEST 0
	
	testContext := GadgetContext {
		Config : TemplateConfig("thisisthename", fmt.Sprintf("%s", initUu1), fmt.Sprintf("%s", initUu2)),
	}
	
	args := []string{ "onboot", "newonboot" }
	
	GadgetAdd(args, &testContext)
	testContext.Config.Onboot[0].Alias = "justtobesureit'spopulated"
	CleanConfig(testContext.Config)
	
	if testContext.Config.Onboot[0].Alias != "" {
		t.Error("failed to clean config")
	}
	if testContext.Config.Onboot[0].ImageAlias != "" {
		t.Error("failed to clean config")
	}
	if testContext.Config.Onboot[1].Alias != "" {
		t.Error("failed to clean config")
	}
	if testContext.Config.Onboot[1].ImageAlias != "" {
		t.Error("failed to clean config")
	}
	
}

func TestWalkUp(t *testing.T){
	
	TestGadgetAdd(t)
	
	base, err := WalkUp("/tmp/some/fake/set/of/directories")
	if err != nil {
		t.Error(err)
	}
	
	if base != "/tmp" {
		fmt.Println(base)
		t.Error("failed to clean config")
	}
	
	_, err = WalkUp("/nonexistant")
	if err == nil {
		t.Error("Should have failed to find /nonexistant/gadget.yml")
	}
	
}


func TestLoadConfig(t *testing.T){
	
}


func TestFind(t *testing.T){
	
}


