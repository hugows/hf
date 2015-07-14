package main

import (
	"fmt"
	"testing"
    "runtime"
)


var FILES = []string{"a", "b", "space jam"}

func TestRuncmdSimple(t *testing.T) {

	if runtime.GOOS != "windows" {
        if err := runCmdWithArgs("test", "cat $FILES", false, FILES); err != nil {
            fmt.Println(err)
    		t.FailNow()
    	}
    }
    var cmd string
    if runtime.GOOS == "windows" {
        cmd = "type $FILES"
    } else {
        cmd = "cat $FILES"
    }
    if err := runCmdWithArgs("test", cmd, true, FILES); err != nil {
        t.FailNow()
    }
}

func TestRuncmdPiped(t *testing.T) {
    var cmd string
    cmd = "cat $FILES | wc -l"

    if runtime.GOOS != "windows" {
    	if err := runCmdWithArgs("test", cmd, true, FILES); err != nil {
    		t.FailNow()
    	}
    }
}
