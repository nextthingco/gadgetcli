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
	"fmt"
	"github.com/nextthingco/libgadget"
	"github.com/satori/go.uuid"
	"reflect"
	"testing"
)

func TestGadgetAdd(t *testing.T) {

	// USAGE TEST PREP
	emptyContext := libgadget.GadgetContext{
		WorkingDirectory: "/tmp",
	}
	initArgs := []string{"these", "aren't", "used"}
	err := GadgetInit(initArgs, &emptyContext)
	if err != nil {
		t.Error("failed to initialize /tmp/gadget.yml for tests")
	}

	// USAGE TEST 0
	args := []string{"blorp", "bleep"}
	err = GadgetAdd(args, &emptyContext)
	if err == nil {
		t.Error("Should have failed to `gadget add bleep blorp`")
	}
	// USAGE TEST 1
	args = []string{"service", "onboot"}
	err = GadgetAdd(args, &emptyContext)
	if err != nil {
		t.Error("Should have succeeded in `gadget add service onboot`", err)
	}
	// USAGE TEST 2
	args = []string{"onboot", "service"}
	err = GadgetAdd(args, &emptyContext)
	if err != nil {
		t.Error("Should have succeeded in `gadget add onboot service`", err)
	}
	// USAGE TEST 3
	args = []string{"onboot"}
	err = GadgetAdd(args, &emptyContext)
	if err == nil {
		t.Error("Should have failed to `gadget add onboot`", err)
	}
	// USAGE TEST 4
	args = []string{"prj_name_only"}
	err = GadgetAdd(args, &emptyContext)
	if err == nil {
		t.Error("Should have failed to `gadget add prj_name_only`", err)
	}

	// TEST 0

	testContext := libgadget.GadgetContext{
		Config: libgadget.GadgetConfig{
			Name: "test",
			UUID: "1234",
			Type: "docker",
			Onboot: []libgadget.GadgetContainer{
				{
					Name:  "hello-world",
					Image: "armhf/hello-world",
				},
			},
		},
	}

	expectedContext := libgadget.GadgetContext{
		Config: libgadget.GadgetConfig{
			Name: "test",
			UUID: "1234",
			Type: "docker",
			Onboot: []libgadget.GadgetContainer{
				{
					Name:  "hello-world",
					Image: "armhf/hello-world",
				},
				{
					Name:  "newonboot",
					Image: "test/newonboot",
				},
			},
		},
	}

	args = []string{"onboot", "newonboot"}

	GadgetAdd(args, &testContext)

	testContext.Config.Onboot[0].UUID = ""
	testContext.Config.Onboot[1].UUID = ""

	if !reflect.DeepEqual(testContext, expectedContext) {
		t.Error("not deeply equal")
	}

	// TEST 1

	testContext = libgadget.GadgetContext{
		Config: libgadget.GadgetConfig{
			Name: "test",
			UUID: "1234",
			Type: "docker",
			Onboot: []libgadget.GadgetContainer{
				{
					Name:  "hello-world",
					Image: "armhf/hello-world",
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

func TestCleanConfig(t *testing.T) {

	fmt.Println("TestCleanConfig")

	initUu1 := uuid.NewV4()
	initUu2 := uuid.NewV4()

	// TEST 0

	testContext := libgadget.GadgetContext{
		Config: libgadget.TemplateConfig("thisisthename", fmt.Sprintf("%s", initUu1), fmt.Sprintf("%s", initUu2)),
	}

	args := []string{"onboot", "newonboot"}

	// do some post-load processing, should add alias and imagealias values
	GadgetAdd(args, &testContext)
	// manually add values to be sure
	testContext.Config.Onboot[0].Alias = "justtobesureit'spopulated"
	// now clean
	libgadget.CleanConfig(testContext.Config)

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

func TestWalkUp(t *testing.T) {

	// TestGadgetAdd will have created the /tmp/gadget.yml

	// should walk up nonexistant directories too
	base, err := libgadget.WalkUp("/tmp/some/fake/set/of/directories")
	if err != nil {
		t.Error(err)
	}

	// check return value
	if base != "/tmp" {
		fmt.Println(base)
		t.Error("failed to find /tmp/gadget.yml")
	}

	// test the failure case
	_, err = libgadget.WalkUp("/nonexistant")
	if err == nil {
		t.Error("Should have failed to find /nonexistant/gadget.yml")
	}

}
