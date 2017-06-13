package libgadget

import (
	"testing"
)

func TestPathExists(t *testing.T){
	
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
