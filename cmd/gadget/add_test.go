
package main

import (
	"testing"
	"reflect"
)

func TestGadgetAdd(t *testing.T){
	
	// USAGE TEST PREP
	emptyContext := GadgetContext{
		WorkingDirectory: "/tmp",
	}
	initArgs := []string { "these", "aren't", "used" }
	err := GadgetInit(initArgs, &emptyContext)
	if err != nil {
		t.Error("failed to initialize /tmp/gadget.yml for tests")
	}
	
	// USAGE TEST 0
	args := []string{ "blorp", "bleep" }
	err = GadgetAdd(args, &emptyContext)
	if err == nil {
		t.Error("Should have failed to `gadget add bleep blorp`")
	}
	// USAGE TEST 1
	args = []string{ "service", "onboot" }
	err = GadgetAdd(args, &emptyContext)
	if err != nil {
		t.Error("Should have succeeded in `gadget add service onboot`", err)
	}
	// USAGE TEST 2
	args = []string{ "onboot", "service" }
	err = GadgetAdd(args, &emptyContext)
	if err != nil {
		t.Error("Should have succeeded in `gadget add onboot service`", err)
	}
	// USAGE TEST 3
	args = []string{ "onboot" }
	err = GadgetAdd(args, &emptyContext)
	if err == nil {
		t.Error("Should have failed to `gadget add onboot`", err)
	}
	// USAGE TEST 4
	args = []string{ "prj_name_only" }
	err = GadgetAdd(args, &emptyContext)
	if err == nil {
		t.Error("Should have failed to `gadget add prj_name_only`", err)
	}
	
	// TEST 0
	
	testContext := GadgetContext {
		Config: GadgetConfig{
			Name: "test",
			UUID: "1234",
			Type: "docker",
			Onboot: []GadgetContainer{
				{
					Name:    "hello-world",
					Image:   "armhf/hello-world",
				},
			},
		},
	}
	
	expectedContext := GadgetContext {
		Config: GadgetConfig{
			Name: "test",
			UUID: "1234",
			Type: "docker",
			Onboot: []GadgetContainer{
				{
					Name:    "hello-world",
					Image:   "armhf/hello-world",
				},
				{
					Name:    "newonboot",
					Image:   "test/newonboot",
				},
			},
		},
	}
	
	args = []string{ "onboot", "newonboot" }
	
	GadgetAdd(args, &testContext)
	
	testContext.Config.Onboot[0].UUID = ""
	testContext.Config.Onboot[1].UUID = ""
	
	if ! reflect.DeepEqual(testContext, expectedContext) {
		t.Error("not deeply equal")
	}
	
	// TEST 1
	
	testContext = GadgetContext {
		Config: GadgetConfig{
			Name: "test",
			UUID: "1234",
			Type: "docker",
			Onboot: []GadgetContainer{
				{
					Name:    "hello-world",
					Image:   "armhf/hello-world",
				},
			},
		},
	}
	
	expectedContext.Config.Onboot[1].Image = "gadgetcli/newonboot"
	
	GadgetAdd(args, &testContext)
	
	testContext.Config.Onboot[0].UUID = ""
	testContext.Config.Onboot[1].UUID = ""
	
	if reflect.DeepEqual(testContext, expectedContext) {
		t.Error("is deeply equal, but shouldn't have been, is there a gadget.yml above me?")
	}
	
}
