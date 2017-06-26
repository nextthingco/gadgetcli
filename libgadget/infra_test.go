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
	"testing"
)

func TestPathExists(t *testing.T) {

	// USAGE TEST PREP
	peBool, err := PathExists("/tmp")
	if err != nil {
		t.Error("Something went wrong while looking for /tmp")
	}

	if peBool == false {
		t.Error("I think we can test against /tmp as long as you're not on windows, right?")
	}

	peBool, err = PathExists("/somefakedir")
	if err != nil {
		t.Error("Something went wrong while looking for /somefakedir")
	}
	if peBool != false {
		t.Error("Do you actually have a /somefakedir?")
	}

}
