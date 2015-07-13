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

const usage string = `The basic command is:

  hf [path] [command]

The default path is the current folder.
The default command is the first valid of $GIT_EDITOR, $EDITOR, or vim 
(subl on Windows).

To find in your git project, use:

  hf -git [command]

 (if -git is provided, then next argument (optional) is assumed to be a command)

If the binary name is 'hfg', the -git option is assumed.
This was done because Windows users have no easy way of creating command aliases.

A -cmd=<yourcmd> option is provided to simplify aliases. For example, I defined 
iga (interactive git add) like this:

alias iga='hf -cmd="git add"'

Examples:
    hf
    hf -git
    hf -git vim
    hf ~/go/src/
    hfg rm
    hf . rm

Inside the app:
    a-z0-9      Edit input string (those also work: backspace/C-a/C-e/C-k/Home/End)
    Up/down     Move cursor to next file
    Space       Toggle mark for current file and move to next
    C-t         Toggle mark for all files 
    C-s         Toggle "run command in shell"
    TAB         Jump to edit command (and back)
    RET         Run command
    ESC         Quit

When editing the command, the string $FILES is special and will
replaced by the select (or marked) files, properly quoted.

`

type Options struct {
	debug     bool   // debug mode (print stats when closing, etc)
	fakefiles int    // generate fake filenames for testing performance
	git       bool   // user wants to search in git project
	rootDir   string // path to recursively search
	runCmd    string // initial command to run

	folderDisplay string // string to display in modeline
}

func ParseArgs() (opts *Options, err error) {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, usage)
	}

	git := flag.Bool("git", false, "Find in current git project instead of folder")
	cmd := flag.String("cmd", "", "Command to run (alternate syntax to simplify alias)")
	debug := flag.Bool("debug", false, "Print stats in the end (debug only)")
	fakefiles := flag.Int("fakefiles", 0, "Generate N fake file names for testing")

	flag.Parse()

	opts = &Options{
		debug:     *debug,
		fakefiles: *fakefiles,
		git:       *git,
		runCmd:    "",
		rootDir:   ".",
	}

	// this is hacky :(
	command := os.Args[0]
	if strings.HasSuffix(command, "hfg") || strings.HasSuffix(command, "hfg.exe") {
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
	hasCmd := (*cmd != "")

	if *cmd != "" {
		opts.runCmd = *cmd
	} else {
		opts.runCmd = defaultCmd
	}

	argLen := len(flag.Args())

	if argLen >= 3 || (argLen == 2 && (opts.git || hasCmd)) || (argLen == 1 && opts.git && hasCmd) {
		err = errors.New("Could not parse options")
	} else if opts.git && !hasCmd && argLen == 1 {
		opts.runCmd = flag.Arg(0)
	} else if !opts.git && argLen == 1 {
		opts.rootDir = flag.Arg(0)
	} else if !opts.git && !hasCmd && argLen == 2 {
		opts.rootDir = flag.Arg(0)
		opts.runCmd = flag.Arg(1)
	}

	opts.rootDir = filepath.Clean(opts.rootDir)
	if withoutLinks, err := filepath.EvalSymlinks(opts.rootDir); err == nil {
		opts.rootDir = withoutLinks
	}

	if opts.git {
		opts.folderDisplay += "git:"
	} else {
		opts.folderDisplay += "cwd:"
	}
	opts.folderDisplay += opts.rootDir

	return
}
