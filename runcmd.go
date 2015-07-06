package main

import (
	"log"
	"os"
	"os/exec"
	"strings"
)

// built-in windows:
// cmd := exec.Command("cmd","/c","del","a")
func runCmdWithArgs(rundir, rawcmd string) {
	parts := strings.Split(rawcmd, " ")
	cmd := exec.Command(parts[0], parts[1:]...)
	cmd.Dir = rundir
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}
