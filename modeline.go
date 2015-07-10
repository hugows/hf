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

	folder string

	// user input
	input *Editbox
}

func NewModeline(folder string) *Modeline {
	input := new(Editbox)
	input.fg = termbox.ColorDefault
	input.bg = termbox.ColorDefault
	return &Modeline{
		input:  input,
		folder: folder,
	}
}

func (m *Modeline) Summarize(results *ResultsView) string {
	sel := results.result_selected + 1
	if results.resultCount == 0 {
		sel = 0
	}

	var s string

	if m.walkFinished {
		s = strconv.Itoa(sel) + "/" + strconv.Itoa(results.resultCount)
	} else {
		s = strconv.Itoa(sel) + "/?"
		if m.paused {
			s += " paused"
		}

	}

	// if !m.walkFinished {
	// 	s += "?"
	// }
	// if m.paused {
	// 	s += " paused"
	// }
	// s += ")"
	return s
}

func (m *Modeline) Draw(x, y, w int, results *ResultsView, active bool) {
	coldef := termbox.ColorDefault
	spaceForCursor := 2

	text := m.folder + " " + m.Summarize(results) //.Summarize(m.paused)

	tbprint(w-len(text), y, termbox.ColorCyan|termbox.AttrBold, coldef, text)

	// modeline.Draw(2, , w-2, 1)
	spaceLeft := w - spaceForCursor - len(text)
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
