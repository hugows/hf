package main

import (
	"log"
	"os"
	"os/exec"
	"strings"
)

func runCmdWithArgs(rawcmd string) {
	parts := strings.Split(rawcmd, " ")
	cmd := exec.Command(parts[0], parts[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}
