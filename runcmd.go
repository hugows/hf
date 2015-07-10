package main

import (
	"os"
	"os/exec"
	"strings"
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

func runCmdWithArgs(rundir, rawcmd string, files []string) {
	words := strings.Split(rawcmd, " ")

	cmd := words[0]
	args := make([]string, 0, len(words))
	for _, w := range words[1:] {
		if w == "$FILES" {
			for _, f := range files {
				args = append(args, f)
			}
		} else {
			args = append(args, w)
		}
	}

	// fmt.Println("run internal", rawcmd, cmd, words, args)
	// err := runCmdInternal(rundir, cmd, args)
	runCmdInternal(rundir, cmd, args)

	// FIX THIS.
	// if err != nil {
	// 	var newcmd string
	// 	newargs := make([]string, len(args))
	// 	if runtime.GOOS == "windows" {
	// 		newcmd = "cmd"
	// 		newargs = append(newargs, "/c")
	// 	} else {
	// 		newcmd = "sh"
	// 		newargs = append(newargs, "-c")
	// 	}

	// 	newargs = append(newargs, cmd)
	// 	for _, a := range args {
	// 		newargs = append(newargs, a)
	// 	}

	// 	err = runCmdInternal(rundir, newcmd, newargs)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// }

}
