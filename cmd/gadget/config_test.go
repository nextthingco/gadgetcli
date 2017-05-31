
package main

import (
	"fmt"
	"testing"
	"reflect"
	"github.com/satori/go.uuid"
)

func TestTemplateConfig(t *testing.T){

	initUu1 := uuid.NewV4()
	initUu2 := uuid.NewV4()

	// TEST 0
	
	expectedContext := GadgetContext {
		Config: GadgetConfig{
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
		},
	}
		
	testContext := TemplateConfig("thisisthename", fmt.Sprintf("%s", initUu1), fmt.Sprintf("%s", initUu2))
	
	if reflect.DeepEqual(testContext, expectedContext) {
		t.Error("is deeply equal, but shouldn't have been, is there a gadget.yml above me?")
	}
	
}
