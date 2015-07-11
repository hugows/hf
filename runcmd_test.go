package main

import (
	"fmt"
	"testing"
)

var FILES = []string{"a", "b", "space jam"}

func TestRuncmdSimple(t *testing.T) {
	if err := runCmdWithArgs("test", "cat $FILES", false, FILES); err != nil {
		t.FailNow()
	}
}

func TestRuncmdPiped(t *testing.T) {
	if err := runCmdWithArgs("test", "cat $FILES | wc -l", true, FILES); err != nil {
		fmt.Println(err)
		t.FailNow()
	}
}
