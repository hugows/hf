package main

import (
	"fmt"

	"github.com/nsf/termbox-go"
)

type Modeline struct {
	// dimensions
	x, y, w int

	// feedback for user
	paused bool

	// user input
	input *Editbox
}

func NewModeline(x, y, w int) *Modeline {
	e := new(Editbox)
	return &Modeline{x, y, w, false, e}
}

func (results *ResultsView) Summarize(paused bool) string {
	sel := results.result_selected + 1
	if results.result_count == 0 {
		sel = 0
	}

	s := fmt.Sprintf("(%d/%d", sel, results.result_count)
	if paused {
		s += " paused"
	}
	s += ")"
	return s
}

func (m *Modeline) Draw(results *ResultsView) {
	coldef := termbox.ColorDefault
	spaceForCursor := 2
	summary := results.Summarize(m.paused)

	tbprint(m.w-len(summary), m.y, termbox.ColorCyan|termbox.AttrBold, coldef, summary)

	// modeline.Draw(2, , w-2, 1)
	spaceLeft := m.w - spaceForCursor - len(summary)
	m.input.Draw(m.x+spaceForCursor, m.y, spaceLeft)
	termbox.SetCell(0, m.y, '>', coldef, coldef)
	termbox.SetCursor(spaceForCursor+m.input.CursorX(), m.y)
}

func (m *Modeline) Contents() string {
	return string(m.input.text)
}

func (m *Modeline) Pause() {
	m.paused = true
}
func (m *Modeline) Unpause() {
	m.paused = false
}
