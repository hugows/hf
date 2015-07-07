package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func firstNonEmpty(arr []string) string {
	for _, s := range arr {
		if len(s) > 0 {
			return s
		}
	}
	return ""
}

func isVolumeRoot(path string) bool {
	return os.IsPathSeparator(path[len(path)-1])
}

func isGitRoot(path string) bool {
	gitpath := filepath.Join(path, ".git")

	fi, err := os.Stat(gitpath)
	if err != nil {
		return false
	}

	return fi.IsDir()
}

func findGitRoot(path string) (bool, string) {
	if absRoot, err := filepath.Abs(path); err == nil {
		path = absRoot
	}
	if withoutLinks, err := filepath.EvalSymlinks(path); err == nil {
		path = withoutLinks
	}

	path = filepath.Clean(path)

	for isVolumeRoot(path) == false {
		if isGitRoot(path) {
			return true, path
		} else {
			if path == filepath.Dir(path) {
				panic("findGitRoot will loop")
			}
			path = filepath.Dir(path)
		}
	}
	return false, ""
}

const usage string = `Use happyfinder like this:

  hf [path] [command]

The default path is the current folder. 
The default command is the first valid of $GIT_EDITOR, $EDITOR, or vim (subl on Windows).

To find in your git project, use:

  hf -git [command]

 (if -git is provided, then next argument (optional) is assumed to be a command)

If the binary name is 'hfg', the -git option is assumed. 
This was done because Windows users have no easy way of creating command aliases.

Examples:
    hf
    hf -git
    hf -git vim
    hf ~/go/src/
    hfg rm
    hf . rm

When running:
    a-z0-9      Edit input string (those also work: backspace/C-a/C-e/C-k/Home/End)
    Up/down     Move cursor to next file
    Space       Toggle mark for current file and move to next
    C-t         Toggle mark for all files 
    TAB         Jump to edit command (and back)
    RET         Run command
    ESC         Quit 

`

type Options struct {
	git          bool   // user wants to search in git project
	command      string // name of binary
	rootDir      string // path to recursively search
	runCmd       string // initial command to run
	initialInput string // starting input to speed things up (not implemented)
}

func ParseArgs() (opts *Options, err error) {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, usage)
	}

	git := flag.Bool(
		"git",
		false,
		"Find in current git project instead of folder")

	flag.Parse()

	opts = &Options{
		git:     *git,
		command: os.Args[0],
		runCmd:  "vim",
	}

	// this is hacky :(
	if strings.HasSuffix(opts.command, "hfg") {
		opts.git = true
	}

	if opts.git {
		foundGit, gitFolder := findGitRoot(".")
		if !foundGit {
			err = errors.New("Git project not found")
			return
		} else {
			opts.rootDir = gitFolder
		}
	}

	var defaultEditor string
	if runtime.GOOS == "windows" {
		defaultEditor = "subl"
	} else {
		defaultEditor = "vim"
	}
	defaultCmd := firstNonEmpty([]string{os.Getenv("GIT_EDITOR"), os.Getenv("EDITOR"), defaultEditor})
	opts.runCmd = defaultCmd

	switch len(flag.Args()) {
	case 0:
		if !opts.git {
			opts.rootDir = "."
		}
	case 1:
		if opts.git {
			opts.runCmd = flag.Arg(0)
		} else {
			opts.rootDir = flag.Arg(0)
		}
	case 2:
		if opts.git {
			opts.runCmd = flag.Arg(0)
			opts.initialInput = flag.Arg(1)
		} else {
			opts.rootDir = flag.Arg(0)
			opts.runCmd = flag.Arg(1)
		}
	case 3:
		if opts.git {
			err = errors.New("Could not parse options")
			return
		} else {
			opts.rootDir = flag.Arg(0)
			opts.runCmd = flag.Arg(1)
			opts.initialInput = flag.Arg(2)
		}
	default:
		err = errors.New("Could not parse options")
	}

	return
}
