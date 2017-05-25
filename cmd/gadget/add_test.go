
package main

import (
	"testing"
	"reflect"
	)

func TestGadgetAdd(t *testing.T){
	
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
					Image:   "test/newonbodot",
				},
			},
		},
	}
	
	args := []string{ "onboot", "newonboot" }
	
	GadgetAdd(args, &testContext)
	
	testContext.Config.Onboot[0].UUID = ""
	testContext.Config.Onboot[1].UUID = ""
	
	if ! reflect.DeepEqual(testContext, expectedContext) {
		t.Error("not deeply equal")
	}
	
}
