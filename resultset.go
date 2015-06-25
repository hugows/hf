package main

import (
	"sort"
	"strings"
)

type ResultArray []*Result

type ResultSet struct {
	results ResultArray
	queue   []string
	count   int
}

// Just for sorting
func (rc ResultArray) Len() int {
	return len(rc)
}
func (rc ResultArray) Swap(i, j int) {
	rc[i], rc[j] = rc[j], rc[i]
}
func (rc ResultArray) Less(i, j int) bool {
	return rc[i].score > rc[j].score
}

func (rs *ResultSet) Insert(entry string) {
	result := new(Result)
	result.contents = strings.ToLower(entry)
	result.displayContents = entry
	rs.results = append(rs.results, result)
	rs.count++
}

func (rs *ResultSet) Queue(entry string) {
	if rs.queue == nil {
		rs.queue = make([]string, 0, 100)
	}
	rs.queue = append(rs.queue, entry)
}

func (rs *ResultSet) FlushQueue() {
	var entry string
	for len(rs.queue) > 0 {
		entry, rs.queue = rs.queue[len(rs.queue)-1], rs.queue[:len(rs.queue)-1]
		rs.Insert(entry)
	}
	rs.queue = nil
}

// Sync, blocking

func (rs *ResultSet) Filter(when int64, userinput string) (filtered ResultSet) {

	if len(userinput) == 0 {
		filtered.results = rs.results
		filtered.count = len(rs.results)

		for _, res := range filtered.results {
			res.highlighted = nil
		}

		return
	}

	filtered.results = make(ResultArray, 0, 100)
	filtered.count = 0

	// Filter
	for _, entry := range rs.results {
		if global_lastkeypress > when {
			break // partial
		}
		best := score2(when, entry.contents, userinput)
		entry.score, entry.highlighted = best.score, best.highlight
		if entry.score > 0 {
			filtered.results = append(filtered.results, entry)
			filtered.count++
		}
	}

	// Sort
	sort.Sort(filtered.results)

	// TODO: better cursor behaviouree
	// r.SelectFirst()
	return
}

// func (rs *ResultSet) Filter(when int64, userinput string) (filtered ResultSet, cancelled bool) {
// 	var a ResultSet
// 	var b bool
// 	metricsFilter.Time(func() {
// 		a, b = rs.FilterInternal(when, userinput)
// 	})
// 	// t.Update(47)
// 	return a, b
// }

// func (rs *ResultSet) AsyncFilter(userinput string, resultCh chan ResultArray) {
// 	filtered := new(ResultSet)

// 	if len(userinput) == 0 {
// 		filtered.results = rs.results
// 		filtered.count = len(rs.results)

// 		for _, res := range filtered.results {
// 			res.highlighted = nil
// 		}

// 		// r.SelectFirst() // TODO
// 		return
// 	}

// 	filtered.results = make(ResultArray, 0, 100)
// 	filtered.count = 0

// 	// Filter
// 	for _, entry := range rs.results {
// 		if termkey.Peek() {
// 			return
// 		}
// 		best := score2(entry.contents, userinput)
// 		entry.score, entry.highlighted = best.score, best.highlight
// 		if entry.score > 0 {
// 			filtered.results = append(filtered.results, entry)
// 			filtered.count++
// 		}
// 	}

// 	// Sort
// 	sort.Sort(filtered.results)

// 	resultCh <- filtered.results

// 	// TODO: better cursor behaviouree
// 	// r.SelectFirst()
// 	return
// }

func (rs *ResultSet) FilterCancel(userinput string, cancel chan bool) (filtered ResultSet) {
	if len(userinput) == 0 {
		filtered.results = rs.results
		filtered.count = len(rs.results)

		for _, res := range filtered.results {
			res.highlighted = nil
		}

		// r.SelectFirst() // TODO
		return
	}

	filtered.results = make(ResultArray, 0, 100)
	filtered.count = 0

	// Filter
	rchan := make(chan *Result)
	quit := make(chan bool)

	go func() {
		for _, entry := range rs.results {
			best := score2(0, entry.contents, userinput)
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
				filtered.results = append(filtered.results, res)
				filtered.count++
			}
		case <-quit:
			break Loop
		case <-cancel:
			return
		}
	}

	// Sort
	sort.Sort(filtered.results)

	// TODO: better cursor behaviouree
	// r.SelectFirst()
	return
}

// func (rs *ResultSet) FilterManager(inputCh <-chan string, resultCh chan<- ResultSet) {
// 	cancel := make(chan bool)
// 	tmpResults := make(chan ResultSet, 1000)

// 	filter := func(u string, cancel chan bool) {
// 		c := make(chan ResultSet)
// 		filtercancel := make(chan bool)
// 		go func() { c <- rs.FilterCancel(u, filtercancel) }()

// 		select {
// 		case filtered := <-c:
// 			tmpResults <- filtered
// 		case <-cancel:
// 			// select {
// 			// case filtercancel <- true:
// 			// default:
// 			// }
// 			break
// 		}

// 	}

// 	for {
// 		select {
// 		case input := <-inputCh:
// 			select {
// 			case cancel <- true:
// 			default:
// 			}

// 			go filter(input, cancel)
// 		case res := <-tmpResults:
// 			resultCh <- res
// 		}
// 	}
// }

// func (rs *ResultSet) AsyncFilter(userinput string, resultCh chan<- ResultSet, cancel <-chan bool) {
// 	temp := make(chan ResultSet)

// 	go func() {
// 		temp <- rs.Filter(userinput)
// 	}()

// 	go func() {
// 		select {
// 		case <-cancel:
// 			break
// 		case r := <-temp:
// 			resultCh <- r
// 		}
// 	}()
// }
