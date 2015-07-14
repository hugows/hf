// +build !windows

package main

import (
	"fmt"
	"testing"
)

var FILES = []string{"a", "b", "space jam"}

func TestRuncmdSimple(t *testing.T) {
	cmd := "cat $FILES"

	// Test without shell	
	if err := runCmdWithArgs("test", cmd, false, FILES); err != nil {
		t.FailNow()
	}

	// Test with shell (cmd /c on windows, sh -c on linux/osx)
	if err := runCmdWithArgs("test", cmd, true, FILES); err != nil {
		t.FailNow()
	}
}

func TestRuncmdPiped(t *testing.T) {
	cmd := "cat $FILES | wc -l"

	if err := runCmdWithArgs("test", cmd, true, FILES); err != nil {
		t.FailNow()
	}
}
