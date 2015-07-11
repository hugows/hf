package main

import (
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

func runCmdInternal(dir string, cmd []string) error {
	exec := exec.Command(cmd[0], cmd[1:]...)
	exec.Dir = dir
	exec.Stdin = os.Stdin
	exec.Stdout = os.Stdout
	exec.Stderr = os.Stderr
	err := exec.Run()
	return err
}

// receives
//   ["git add $FILES", "$FILES", ["a","b"] ]
// returns
//   ["git", "add", "a", "b"]
func expandInArray(arr []string, when string, with []string) []string {
	expanded := make([]string, 0, len(arr))

	for _, e := range arr {
		if e == when {
			for _, arg := range with {
				expanded = append(expanded, arg)
			}
		} else {
			expanded = append(expanded, e)
		}
	}
	return expanded
}

// On OSX, both shell=true and shell=false work w/ GUI apps. Both work with text vim too.
// On Windows, I think GUI apps will lock the command prompt if shell=false (CHECK).
//
// Are there any disadvantadges in always using shell=true?
func runCmdWithArgs(dir string, userCommand string, shell bool, files []string) error {
	var cmd []string

	if shell {
		if runtime.GOOS == "windows" {
			cmd = []string{"cmd", "/c"}
		} else {
			cmd = []string{"sh", "-c"}
		}
		quotedFiles := make([]string, len(files))
		for i, f := range files {
			quotedFiles[i] = strconv.Quote(f)
		}
		filesString := strings.Join(quotedFiles, " ")
		cmdReplaced := strings.Replace(userCommand, "$FILES", filesString, -1)
		cmd = append(cmd, cmdReplaced)
	} else {
		cmd = strings.Split(userCommand, " ")
		cmd = expandInArray(cmd, "$FILES", files)
	}

	return runCmdInternal(dir, cmd)
}
