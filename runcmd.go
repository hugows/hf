package main

import (
	"log"
	"os"
	"os/exec"
	"strings"
	"runtime"
)

func runCmdInternal(rundir, rawcmd string) error {
	parts := strings.Split(rawcmd, " ")
	cmd := exec.Command(parts[0], parts[1:]...)
	cmd.Dir = rundir
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return err
}

// built-in windows:
// cmd := exec.Command("cmd","/c","del","a")
func runCmdWithArgs(rundir, rawcmd string) {
	err := runCmdInternal(rundir, rawcmd)

	if err != nil {
		if runtime.GOOS == "windows" {
			rawcmd = "cmd /k" + rawcmd
		} else {
			rawcmd = "sh -c" + rawcmd
		}
		err = runCmdInternal(rundir, rawcmd)
		if err != nil {
			log.Fatal(err)
		}
	}
}
