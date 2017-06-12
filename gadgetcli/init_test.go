
package main

import (
	"testing"
	//~ "reflect"
)

func TestGadgetInit(t *testing.T){
	
	// USAGE TEST PREP
	initContext := GadgetContext{
		WorkingDirectory: "/tmp",
	}
	
	initArgs := []string { "these", "aren't", "used" }
	err := GadgetInit(initArgs, &initContext)
	if err != nil {
		t.Error("failed to initialize /tmp/gadget.yml for tests")
	}
	
	err = initContext.LoadConfig()
	if err != nil {
		t.Error("failed to load /tmp/gadget.yml for tests")
	}
	if initContext.Config.Name != "tmp" {
		t.Error("/tmp/gadget.yml name should have been 'tmp'")
	}
	if initContext.Config.UUID == "" {
		t.Error("/tmp/gadget.yml project uuid is empty")
	}
	if initContext.Config.Onboot[0].UUID == "" {
		t.Error("/tmp/gadget.yml container[hello-world] uuid is empty")
	}
	if initContext.Config.Onboot[0].Directory != "" {
		t.Error("/tmp/gadget.yml container[hello-world] directory is not empty")
	}
	
}
