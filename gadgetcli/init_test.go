/*
This file is part of the Gadget command-line tools.
Copyright (C) 2017 Next Thing Co.

Gadget is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 2 of the License, or
(at your option) any later version.

Gadget is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with Gadget.  If not, see <http://www.gnu.org/licenses/>.
*/

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
