package main

import (
	"log"
	"os"
	"os/exec"
	"runtime"
)

func runCmdInternal(rundir, cmd string, args []string) error {
	exec := exec.Command(cmd, args...)
	exec.Dir = rundir
	exec.Stdin = os.Stdin
	exec.Stdout = os.Stdout
	exec.Stderr = os.Stderr
	err := exec.Run()
	return err
}

func runCmdWithArgs(rundir, cmd string, args []string) {
	err := runCmdInternal(rundir, cmd, args)

	if err != nil {
		var newcmd string
		newargs := make([]string, len(args))
		if runtime.GOOS == "windows" {
			newcmd = "cmd"
			newargs = append(newargs, "/c")
		} else {
			newcmd = "sh"
			newargs = append(newargs, "-c")
		}

		newargs = append(newargs, cmd)
		for _, a := range args {
			newargs = append(newargs, a)
		}

		err = runCmdInternal(rundir, newcmd, newargs)
		if err != nil {
			log.Fatal(err)
		}
	}
}
