# happyfinder [![Build Status](https://travis-ci.org/hugows/hf.svg?branch=master)](https://travis-ci.org/hugows/hf)

hf is a command line utility to quickly find files and execute a command - something like Helm/Anything/CtrlP for the terminal. It tries to find the best match, like other fuzzy finders (Sublime, ido, Helm).

Here is it in action:

![happyfinder on osx](http://g.recordit.co/bWae8XRKMV.gif)

If you have any suggestions, please open an issue.

## Similar projects

This project wasn't implemented as an improvement to any of the projects below - I only knew about the text editor plugins (Anything, Helm, CtrlP, Sublime's native matching...). I wanted something similar to that for the command line, and a reason to use termbox-go, a great library for creating terminal-based apps.

- [fzf is a blazing fast command-line fuzzy finder written in Go](https://github.com/junegunn/fzf)
- [Selecta is a fuzzy text selector for files and anything else you need to select](https://github.com/garybernhardt/selecta/)
- [Pick is "just like Selecta, but faster"](https://robots.thoughtbot.com/announcing-pick)
- [icepick is a reimplementation of Selecta in Rust](https://github.com/felipesere/icepick)

Those projects take the elegant route - the Unix way - in the sense they don't reimplement the "find files recursively" part. I didn't consider that route because I wanted hf to work in the Windows prompt as well.

## Installation

If you have Go configured in your system, you should be able to install 
(and update) happyfinder with the command:

```
go install github.com/hugows/hf@latest
```

## Usage

```
The basic command is:

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
```
