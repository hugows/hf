# happyfinder [![Build Status](https://travis-ci.org/hugows/hf.svg?branch=master)](https://travis-ci.org/hugows/hf)

hf is a command line utility to quickly find files and execute a command - something like Helm/Anything/CtrlP for the terminal.

Here is it in action:

![happyfinder on osx](http://g.recordit.co/bWae8XRKMV.gif)

If you have any suggestions, please open an issue!

## Installation

You should be able to install (and update) /happyfinder/ with the command:

```
go get -u github.com/hugows/hf
```

## Usage

```
Use happyfinder like this:

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
```
