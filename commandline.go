package main

import "github.com/nsf/termbox-go"

type CommandLine struct {
	// dimensions
	x, y, w int

	// allow user to edit it
	input *Editbox

	// program to call
	originalCmd string
}

func NewCommandLine(x, y, w int, cmd string) *CommandLine {
	input := new(Editbox)
	input.fg = termbox.ColorRed
	input.bg = termbox.ColorDefault

	return &CommandLine{
		x: x, y: y, w: w,
		input:       input,
		originalCmd: cmd,
	}
}

func (cmd *CommandLine) Update(r *Result) {
	if r != nil {
		cmd.input.text = []byte(r.displayContents)
	} else {
		cmd.input.text = []byte("")
	}
}

func (cmd *CommandLine) Draw(x, y, w int, active bool) {
	// tclearcolor(x, y, w, 1, bg)
	// text := cmd.originalCmd + " ---------" + cmd.input.Contents() //+ " (RET to run)"
	// tbprint(x, y, fg, bg, text)

	cmd.input.Draw(x, y, w)

	if active {
		termbox.SetCursor(cmd.input.CursorX(), cmd.y)
	}
}
