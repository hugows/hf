// +build !windows

package main

import (
	"fmt"
	"testing"
)

var FILES = []string{"a", "b", "space jam"}

func TestRuncmdSimple(t *testing.T) {

	if err := runCmdWithArgs("test", "cat $FILES", false, FILES); err != nil {
		fmt.Println(err)
		t.FailNow()
	}

	// I'm not sure if you wanted both tests ran on linux, and only the latter
	// in windows.

	cmd := "cat $FILES"

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
