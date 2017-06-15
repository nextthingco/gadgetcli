package main

import (
	"github.com/nextthingco/libgadget"
	"testing"
)

func TestGadgetInit(t *testing.T) {

	// USAGE TEST PREP
	initContext := libgadget.GadgetContext{
		WorkingDirectory: ".",
	}

	initArgs := []string{"these", "aren't", "used"}
	err := GadgetInit(initArgs, &initContext)
	if err != nil {
		t.Error("failed to initialize ./gadget.yml for tests")
	}

	err = initContext.LoadConfig()
	if err != nil {
		t.Error("failed to load ./gadget.yml for tests")
	}
	if initContext.Config.Name != "gadgetcli" {
		t.Error("./gadget.yml name should have been 'gadgetcli'")
	}
	if initContext.Config.UUID == "" {
		t.Error("./gadget.yml project uuid is empty")
	}
	if initContext.Config.Onboot[0].UUID == "" {
		t.Error("./gadget.yml container[hello-world] uuid is empty")
	}
	if initContext.Config.Onboot[0].Directory != "" {
		t.Error("./gadget.yml container[hello-world] directory is not empty")
	}

}
