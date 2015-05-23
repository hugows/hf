package main

import (
	"fmt"

	"github.com/nsf/termbox-go"
)

type Statusline struct {
}

func (m *Statusline) Draw(x, y, w int, r *Results) {
	sel := r.result_selected + 1
	if r.result_count == 0 {
		sel = 0
	}
	line := fmt.Sprintf("(%d/%d)", sel, r.result_count)
	tbprint(x, y, w, termbox.ColorDefault|termbox.AttrReverse, termbox.ColorDefault, line)
}
