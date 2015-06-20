package main

import (
	"sort"
	"strings"

	"github.com/nsf/termbox-go"
)

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
	initialset ResultCollection
	currentset ResultCollection

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

func (r *Results) SelectPrevious() *Result {
	if r.result_selected > 0 {
		r.result_selected--
	}
	if r.top_result > 0 {
		r.top_result--
		r.bottom_result--
	}

	return r.currentset[r.result_selected]
}

func (r *Results) SelectNext() *Result {
	if r.result_selected < (r.result_count - 1) {
		r.result_selected++

		if r.result_selected >= r.bottom_result {
			r.top_result++
			r.bottom_result++
		}
	}

	return r.currentset[r.result_selected]
}

func (r *Results) Insert(s string) {
	result := new(Result)
	result.contents = strings.ToLower(s)
	result.displayContents = s
	r.initialset = append(r.initialset, result)
	r.result_count++
}

func (r *Results) Queue(s string) {
}

func tbprint(x, y int, fg, bg termbox.Attribute, msg string) {
	for _, c := range msg {
		termbox.SetCell(x, y, c, fg, bg)
		x++
	}
}

func (r *Results) Draw() {

	tclear(r.x, r.y, r.w, r.h)

	cy := r.y

	for cnt, res := range r.currentset[r.top_result:r.bottom_result] {
		is_selected := (cnt + r.top_result) == r.result_selected
		res.Draw(r.x, cy, r.w, is_selected)
		cy++
	}
}

func (r *Results) ToggleMark() {
	if r.result_count > 0 {
		r.currentset[r.result_selected].marked = !r.currentset[r.result_selected].marked
		r.SelectNext()
	}
}

func (r *Results) ToggleMarkAll() {
	for _, res := range r.currentset {
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
	r.currentset = r.initialset
}

func (r *Results) Filter(userinput string, keypressed chan bool) {
	if len(userinput) == 0 {
		r.currentset = r.initialset
		r.result_count = len(r.initialset)

		for _, res := range r.currentset {
			res.highlighted = nil
		}

		r.SelectFirst()
		return
	}

	// Optimization
	// Now invalid because results are changing...
	// if len(r.lastuserinput) > 0 && strings.HasPrefix(userinput, r.lastuserinput) {
	// 	initialset = r.currentset
	// 	if len(r.currentset) == 0 {
	// 		r.result_count = 0
	// 		r.SelectFirst()
	// 		return
	// 	}
	// }
	// r.lastuserinput = userinput

	r.currentset = make([]*Result, 0, 100)
	r.result_count = 0

	// Filter
	rchan := make(chan *Result)
	quit := make(chan bool)

	go func() {
		for _, entry := range r.initialset {
			best := score2(entry.contents, userinput)
			entry.score, entry.highlighted = best.score, best.highlight
			rchan <- entry
		}
		quit <- true
	}()

	// Cancellable
Loop:
	for {
		select {
		case res := <-rchan:
			if res.score > 0 {
				r.currentset = append(r.currentset, res)
				r.result_count++
			}
		case <-quit:
			break Loop
		case <-keypressed:
			return
		}
	}

	// Sort
	sort.Sort(r.currentset)

	// TODO: better cursor behaviouree
	r.SelectFirst()

}

func (r *Results) GetSelected() *Result {
	return r.currentset[r.result_selected]
}
