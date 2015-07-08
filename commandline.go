package main

import (
	"fmt"

	"github.com/nsf/termbox-go"
)

type CommandLine struct {
	// allow user to edit it
	input *Editbox

	// program to call
	cmd string

	// cached
	fullCmdline       string
	summarizedCmdline string
}

func NewCommandLine(cmd string) *CommandLine {
	input := new(Editbox)
	input.fg = termbox.ColorRed
	input.bg = termbox.ColorDefault

	return &CommandLine{
		input: input,
		cmd:   cmd,
	}
}

func (cmd *CommandLine) Update(results ResultArray) {
	text := cmd.cmd

	for _, res := range results {
		text = text + " " + res.displayContents
	}
	cmd.input.text = []byte(text)
	cmd.input.MoveCursorToBeginningOfTheLine()
	cmd.fullCmdline = text
	cmd.summarizedCmdline = fmt.Sprintf("%s <%d files...>", cmd.cmd, len(results))
}

func (cmd *CommandLine) SummarizeCommand(maxlen int) string {
	if len(cmd.fullCmdline) > maxlen {
		return cmd.summarizedCmdline
	} else {
		return cmd.fullCmdline
	}
}

func (cmd *CommandLine) Draw(x, y, w int, active bool) {
	if active {
		cmd.input.Draw(x, y, w)
		termbox.SetCursor(cmd.input.CursorX(), y)
	} else {
		tclearcolor(x, y, w, 1, cmd.input.bg)
		tbprint(x, y, cmd.input.fg, cmd.input.bg, cmd.SummarizeCommand(w))
	}

}
