package main

import (
	"sort"
	"strings"
)

type ResultArray []*Result

type ResultSet struct {
	results ResultArray
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
		best := score(when, entry.contents, userinput)
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
