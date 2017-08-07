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

package libgadget

import (
	"fmt"
	"github.com/nextthingco/libgadget"
	"gopkg.in/satori/go.uuid.v1"
	"reflect"
	"testing"
)

func TestTemplateConfig(t *testing.T) {

	fmt.Println("TestTemplateConfig")

	initUu1 := uuid.NewV4()
	initUu2 := uuid.NewV4()

	// TEST 0

	// create the default programatically
	expectedConfig := libgadget.GadgetConfig{
		Spec: Version,
		Name: "thisisthename",
		UUID: fmt.Sprintf("%s", initUu1),
		Type: "docker",
		Onboot: []libgadget.GadgetContainer{
			{
				Name:  "hello-world",
				Image: "arm32v7/hello-world",
				UUID:  fmt.Sprintf("%s", initUu2),
			},
		},
	}

	testConfig := libgadget.TemplateConfig("thisisthename", fmt.Sprintf("%s", initUu1), fmt.Sprintf("%s", initUu2))

	// check for deep equality
	if !reflect.DeepEqual(testConfig, expectedConfig) {
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

func TestLoadConfig(t *testing.T) {

	// run the init test [on pre-existing /tmp/gadget.yml]

	initContext := GadgetContext{
		WorkingDirectory: "/tmp",
	}

	// just run through the LoadConfig for now
	err := initContext.LoadConfig()
	if err != nil {
		t.Error("Failed to TestLoadConfig /tmp/gadget.yml")
	}

}

func TestFind(t *testing.T) {

	// all that Find checks is the Name entry
	testGadCont := libgadget.GadgetContainers{
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
	if !reflect.DeepEqual(returnContainer, testGadCont[0]) {
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
