
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
	
	// create the default programatically
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
	
	// check for deep equality
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
	
	// spot check some values
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
	
	// do some post-load processing, should add alias and imagealias values
	GadgetAdd(args, &testContext)
	// manually add values to be sure
	testContext.Config.Onboot[0].Alias = "justtobesureit'spopulated"
	// now clean
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
	
	// just because TestGadgetAdd creates the /tmp/gadget.yml
	TestGadgetAdd(t)
	
	// should walk up nonexistant directories too
	base, err := WalkUp("/tmp/some/fake/set/of/directories")
	if err != nil {
		t.Error(err)
	}
	
	// check return value
	if base != "/tmp" {
		fmt.Println(base)
		t.Error("failed to find /tmp/gadget.yml")
	}
	
	// test the failure case
	_, err = WalkUp("/nonexistant")
	if err == nil {
		t.Error("Should have failed to find /nonexistant/gadget.yml")
	}
	
}


func TestLoadConfig(t *testing.T){
	
	// run the init test
	TestGadgetInit(t)
	
	initContext := GadgetContext{
		WorkingDirectory: "/tmp",
	}
	
	// just run through the LoadConfig for now
	err := initContext.LoadConfig()
	if err != nil {
		t.Error("Failed to TestLoadConfig /tmp/gadget.yml")
	}
	
}


func TestFind(t *testing.T){
	
	// all that Find checks is the Name entry
	testGadCont := GadgetContainers{
		{
			Name: "test0",
		},
		{
			Name: "test1",
		},
		{
			Name: "test2",
		},
	}
	
	// this should find test0 at testGadCont[0], check for DeepEquality
	returnContainer, err := testGadCont.Find("test0")
	if err != nil {
		t.Error("Failed to testGadCont.Find(\"test0\")")
	}
	if ! reflect.DeepEqual(returnContainer, testGadCont[0]) {
		t.Error("isn't deeply equal, but should have been")
		fmt.Println("%+v", returnContainer)
		fmt.Println("%+v", testGadCont[0])
	}
	
	// test the fail case
	returnContainer, err = testGadCont.Find("shouldfail")
	if err == nil {
		t.Error("Should have failed to testGadCont.Find(\"shouldfail\"), but didn't")
	}
	
}


