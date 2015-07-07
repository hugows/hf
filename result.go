package main

import (
	"strings"

	"github.com/nsf/termbox-go"
)

type Result struct {
	id              int          // unique id to identify a specific result between "sorts"
	contents        string       // user input is matched against lowercase version
	displayContents string       // original filename (or line) to display
	highlighted     map[int]bool // hashmap of the characters to be highlighted
	marked          bool         // true when the current line is selected
	score           int          // what is the score for this particular result?
}

// const CHECKMARK_CHAR = " ✓"
const CHECKMARK_CHAR = " * "

func (res Result) Draw(x, y, w int, selected bool) {
	CHECKMARK_PAD := strings.Repeat(" ", len(CHECKMARK_CHAR))
	const coldef = termbox.ColorDefault

	color := coldef

	if selected {
		color = color | termbox.AttrReverse
	}

	line := ""
	// User doesn't need to the score.... just for debugging.
	//line += fmt.Sprintf("%4d ", res.score)

	if res.marked {
		line += CHECKMARK_CHAR
	} else {
		line += CHECKMARK_PAD
	}
	tclearcolor(x, y, w, 1, color)
	tbprint(x, y, color, color, line)
	x += len(line)

	for idx, c := range res.displayContents {
		fg := color
		bg := color
		if res.highlighted[idx] {
			fg = termbox.ColorGreen | termbox.AttrBold
		}
		termbox.SetCell(x, y, c, fg, bg)
		x++
	}
}
