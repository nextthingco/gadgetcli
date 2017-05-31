
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


