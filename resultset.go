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
		rs.queue = make([]string, 100)
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

func (rs *ResultSet) Filter(userinput string, cancel chan bool) (filtered ResultSet) {
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
