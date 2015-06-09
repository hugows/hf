package main

import (
	"fmt"

	"github.com/nsf/termbox-go"
)

type Modeline struct {
	// dimensions
	x, y, w int

	// user input
	input *Editbox
}

func NewModeline(x, y, w int) *Modeline {
	e := new(Editbox)
	return &Modeline{x, y, w, e}
}

func (results *Results) Summarize() string {
	sel := results.result_selected + 1
	if results.result_count == 0 {
		sel = 0
	}
	return fmt.Sprintf("(%d/%d)", sel, results.result_count)
}

func (m *Modeline) Clear() {

}

func (m *Modeline) Draw(results *Results) {
	coldef := termbox.ColorDefault
	spaceForCursor := 2
	summary := results.Summarize()

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
