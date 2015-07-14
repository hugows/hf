package main

import "testing"

var FILES = []string{"a", "b", "space jam"}

func TestRuncmdSimple(t *testing.T) {
	cmd := "type $FILES"

	if err := runCmdWithArgs("test", cmd, true, FILES); err != nil {
		t.FailNow()
	}
}
