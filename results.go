package main

import (
	"fmt"
	"sort"

	"github.com/nsf/termbox-go"
)

type Result struct {
	contents    string
	highlighted map[int]bool
	marked      bool
	score       int
}

func (res Result) Draw(x, y, w int, selected bool) {
	const coldef = termbox.ColorDefault

	color := coldef
	if selected {
		color = coldef | termbox.AttrReverse
	}

	line := ""
	line += fmt.Sprintf("%4d ", res.score)
	if res.marked {
		line += "*"
	} else {
		line += " "
	}
	tbprint(x, y, len(line), color, color, line)
	x += len(line)

	for idx, c := range res.contents {
		fg := color
		if res.highlighted[idx] {
			fg = color | termbox.ColorGreen | termbox.AttrBold
		}
		termbox.SetCell(x, y, c, fg, color)
		x++
	}

	c := termbox.Cell{Ch: ' ', Fg: color, Bg: color}
	fill(x, y, w, 1, c)
}

type ResultCollection []*Result

func (rc ResultCollection) Len() int {
	return len(rc)
}
func (rc ResultCollection) Swap(i, j int) {
	rc[i], rc[j] = rc[j], rc[i]
}
func (rc ResultCollection) Less(i, j int) bool {
	return rc[i].score > rc[j].score
}

type Results struct {
	// Array of results to be filtered
	allresults ResultCollection
	results    ResultCollection

	// Current user input
	lastuserinput string

	// Visible result lines
	top_result    int
	bottom_result int

	// Total number of results
	result_count int

	// Index of currently selected line
	result_selected int

	// View size
	x, y, h, w int
}

func (r *Results) SelectFirst() {
	r.result_selected = 0
	r.top_result = 0

	if r.result_count > r.h {
		r.bottom_result = r.h
	} else {
		r.bottom_result = r.result_count
	}
}

func (r *Results) SelectPrevious() {
	if r.result_selected > 0 {
		r.result_selected--
	}
	if r.top_result > 0 {
		r.top_result--
		r.bottom_result--
	}
}

func (r *Results) SelectNext() {
	if r.result_selected < (r.result_count - 1) {
		r.result_selected++

		if r.result_selected >= r.bottom_result {
			r.top_result++
			r.bottom_result++
		}
	}
}

func (r *Results) Insert(s string) {
	result := new(Result)
	result.contents = s
	r.allresults = append(r.allresults, result)
	r.result_count++
}

func (r *Results) Queue(s string) {
}

func tbprint(x, y, w int, fg, bg termbox.Attribute, msg string) {
	for _, c := range msg {
		termbox.SetCell(x, y, c, fg, bg)
		x++
	}
}

func (r *Results) Draw() {

	fill(r.x, r.y, r.w, r.h, termbox.Cell{Ch: ' '})

	cy := r.y

	for cnt, res := range r.results[r.top_result:r.bottom_result] {
		is_selected := (cnt + r.top_result) == r.result_selected
		res.Draw(r.x, cy, r.w, is_selected)
		cy++
	}
}

func (r *Results) ToggleMark() {
	if r.result_count > 0 {
		r.results[r.result_selected].marked = !r.results[r.result_selected].marked
		r.SelectNext()
	}
}

func (r *Results) ToggleMarkAll() {
	for _, res := range r.results {
		res.marked = !res.marked
	}
}

func (r *Results) SetSize(x, y, w, h int) {
	r.x, r.y, r.w, r.h = x, y, w, h

	r.top_result = 0
	if r.result_count > r.h {
		r.bottom_result = r.h
	} else {
		r.bottom_result = r.result_count
	}
}

func (r *Results) CopyAll() {
	r.results = r.allresults
}

func (r *Results) Filter(userinput string) {
	if len(userinput) == 0 {
		r.results = r.allresults
		r.result_count = len(r.allresults)

		for _, res := range r.results {
			res.highlighted = nil
		}

		r.SelectFirst()
		return
	}

	initialset := r.allresults

	// Optimization
	// Now invalid because results are changing...
	// if len(r.lastuserinput) > 0 && strings.HasPrefix(userinput, r.lastuserinput) {
	// 	initialset = r.results
	// 	if len(r.results) == 0 {
	// 		r.result_count = 0
	// 		r.SelectFirst()
	// 		return
	// 	}
	// }
	// r.lastuserinput = userinput

	r.results = make([]*Result, 0, 100)
	r.result_count = 0

	// Filter
	for _, res := range initialset {
		best := score2(res.contents, userinput)
		res.score, res.highlighted = best.score, best.highlight
		if res.score > 0 {
			r.results = append(r.results, res)
			r.result_count++
		}
	}

	// Sort
	sort.Sort(r.results)

	// TODO: better cursor behaviouree
	r.SelectFirst()
}
