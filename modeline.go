package main

import (
	"strconv"
	"strings"

	"github.com/nsf/termbox-go"
)

type Modeline struct {
	// dimensions
	x, y, w int

	// feedback for user
	paused       bool
	walkFinished bool

	// user input
	input *Editbox
}

func NewModeline() *Modeline {
	input := new(Editbox)
	input.fg = termbox.ColorDefault
	input.bg = termbox.ColorDefault
	return &Modeline{
		input: input,
	}
}

func (m *Modeline) Summarize(results *ResultsView) string {
	sel := results.result_selected + 1
	if results.result_count == 0 {
		sel = 0
	}

	s := "(" + strconv.Itoa(sel) + "/" + strconv.Itoa(results.result_count)
	if !m.walkFinished {
		s += "+"
	}
	if m.paused {
		s += " paused"
	}
	s += ")"
	return s
}

func (m *Modeline) Draw(x, y, w int, results *ResultsView, active bool) {
	coldef := termbox.ColorDefault
	spaceForCursor := 2

	summary := m.Summarize(results) //.Summarize(m.paused)

	tbprint(w-len(summary), y, termbox.ColorCyan|termbox.AttrBold, coldef, summary)

	// modeline.Draw(2, , w-2, 1)
	spaceLeft := w - spaceForCursor - len(summary)
	m.input.Draw(x+spaceForCursor, y, spaceLeft)
	termbox.SetCell(0, y, '>', coldef, coldef)

	if active {
		termbox.SetCursor(spaceForCursor+m.input.CursorX(), y)
	}
}

func (m *Modeline) Contents() string {
	return strings.ToLower(string(m.input.text))
}

func (m *Modeline) FlagPause(b bool) {
	m.paused = b
}
func (m *Modeline) FlagLastFile() {
	m.walkFinished = true
}
