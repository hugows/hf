package main

import "github.com/nsf/termbox-go"

type ResultsView struct {
	// Array of results to be filtered
	// initialset ResultArray
	results ResultArray

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

func (r *ResultsView) SelectFirst() {
	r.result_selected = 0
	r.top_result = 0

	if r.result_count > r.h {
		r.bottom_result = r.h
	} else {
		r.bottom_result = r.result_count
	}
}

func (r *ResultsView) SelectPrevious() *Result {
	if r.result_selected > 0 {
		r.result_selected--
	}
	if r.top_result > 0 {
		r.top_result--
		r.bottom_result--
	}

	if len(r.results) > 0 && r.result_selected < len(r.results) {
		return r.results[r.result_selected]
	}
	return nil
}

func (r *ResultsView) SelectNext() *Result {
	if r.result_selected < (r.result_count - 1) {
		r.result_selected++

		if r.result_selected >= r.bottom_result {
			r.top_result++
			r.bottom_result++
		}
	}

	if len(r.results) > 0 && r.result_selected < len(r.results) {
		return r.results[r.result_selected]
	}
	return nil
}

func tbprint(x, y int, fg, bg termbox.Attribute, msg string) {
	for _, c := range msg {
		termbox.SetCell(x, y, c, fg, bg)
		x++
	}
}

func (r *ResultsView) Draw() {

	tclear(r.x, r.y, r.w, r.h)

	cy := r.y

	for cnt, res := range r.results[r.top_result:r.bottom_result] {
		is_selected := (cnt + r.top_result) == r.result_selected
		res.Draw(r.x, cy, r.w, is_selected)
		cy++
	}
}

func (r *ResultsView) ToggleMark() {
	if r.result_count > 0 {
		r.results[r.result_selected].marked = !r.results[r.result_selected].marked
		r.SelectNext()
	}
}

func (r *ResultsView) ToggleMarkAll() {
	for _, res := range r.results {
		res.marked = !res.marked
	}
}

func (r *ResultsView) SetSize(x, y, w, h int) {
	r.x, r.y, r.w, r.h = x, y, w, h

	r.top_result = 0
	if r.result_count > r.h {
		r.bottom_result = r.h
	} else {
		r.bottom_result = r.result_count
	}
}

func (r *ResultsView) Update(results ResultArray) {
	r.results = results
	r.result_count = len(results)
	r.SetSize(r.x, r.y, r.w, r.h)
}

func (r *ResultsView) GetSelected() *Result {
	if len(r.results) > 0 && r.result_selected < len(r.results) {
		return r.results[r.result_selected]
	}
	return nil
}
