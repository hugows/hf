package main

import (
	"fmt"
	"time"

	"github.com/nsf/termbox-go"
)

type CommandLine struct {
	// allow user to edit it
	input *Editbox

	// program to call
	cmd string

	// use space for showing errors too
	showingError  bool
	modelineError string

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

func (cmd *CommandLine) ShowError(redraw chan bool, err error) {
	cmd.showingError = true
	cmd.modelineError = "Error: " + err.Error()
	clearErrorTimer := time.NewTimer(1 * time.Second)
	go func() {
		<-clearErrorTimer.C
		cmd.showingError = false
		redraw <- true
	}()
}

func (cmd *CommandLine) Draw(x, y, w int, active bool) {
	if cmd.showingError {
		tclearcolor(x, y, w, 1, cmd.input.bg)
		tbprint(x, y, termbox.ColorRed, cmd.input.bg, cmd.modelineError)
		return
	}

	if active {
		cmd.input.Draw(x, y, w)
		termbox.SetCursor(cmd.input.CursorX(), y)
	} else {
		tclearcolor(x, y, w, 1, cmd.input.bg)
		tbprint(x, y, cmd.input.fg, cmd.input.bg, cmd.SummarizeCommand(w))
	}

}
