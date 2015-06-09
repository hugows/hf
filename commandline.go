package main

import "github.com/nsf/termbox-go"

type CommandLine struct {
	text string
}

func (cmd *CommandLine) Update(s string) {
	cmd.text = s
}

func (cmd *CommandLine) Draw(x, y, w int) {
	bg := termbox.ColorDefault
	fg := termbox.ColorRed
	tclearcolor(x, y, w, 1, bg)
	text := "vim " + cmd.text //+ " (RET to run)"
	tbprint(x, y, fg, bg, text)
}
