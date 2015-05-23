package main

import (
	"fmt"
	"strings"
	"unicode"
)

// "brazil", "bra" 					-> [[0] [1] [2]]
// "hodor", "od"   					-> [[1 3] [2]]
// "/usr/sbin/unsetpassword", "usr" -> [[1 10] [2 5 12 17 18] [3 21]]
func occur(s, chars string) (all [][]int, any bool) {
	all = make([][]int, 0, len(s))
	any = true

	if len(chars) > 0 {
		for _, m := range chars {
			charoccur := make([]int, 0)
			for i, c := range s {
				if c == m || unicode.ToLower(c) == unicode.ToLower(m) {
					charoccur = append(charoccur, i)
				}
			}
			if len(charoccur) == 0 {
				// Some letter had 0 matches... abort
				any = false
				return
			}
			all = append(all, charoccur)
		}
	}

	return
}

// [[a b] [c d]] -> [[a c] [a d] [b c] [b d]]
func combinations(path []int, sets [][]int, acc *[][]int) {
	if len(sets) == 0 {
		*acc = append(*acc, path)
	} else {
		for _, el := range sets[0] {
			path := append(path, el)
			np := make([]int, len(path))
			copy(np, path)
			combinations(np, sets[1:], acc)
		}
	}
}

func singlescore(comb []int, beat int) int {
	score := 0
	lastchpos := comb[0]
	for _, chpos := range comb[1:] {
		dist := chpos - lastchpos

		// Early return optimization
		if dist <= 0 {
			return 1000
		}
		score += dist + 1

		// Early return optimization
		if beat > 0 && score > beat {
			return 1000
		}
		lastchpos = chpos
	}
	return score
}

func score(against string, userinput string) (bestscore int, mhighlight map[int]bool) {
	mhighlight = make(map[int]bool, 0)

	if len(userinput) > len(against) || len(userinput) == 0 || len(against) == 0 {
		bestscore = -1
		return
	}

	all := make([][]int, 0i)
	oc, any := occur(against, userinput)

	// No matches found
	if !any {
		bestscore = -1
		return
	}

	combinations([]int{}, oc, &all)

	bestscore = singlescore(all[0], 0)
	highlight := all[0]

	for _, comb := range all[1:] {
		s := singlescore(comb, bestscore)
		if s < bestscore && s > 0 {
			bestscore = s
			highlight = comb
		}
	}

	for _, idx := range highlight {
		mhighlight[idx] = true
	}

	return bestscore, mhighlight

}

type Matcher struct {
	linepos, inputpos     int
	total_distance        int
	groups                int
	ended, complete_match bool
	chars                 []int

	// Helpers
	lastmatchpos    int
	lenlongestgroup int
	curgrouplen     int
}

type BestScore struct {
	score     int
	groups    int
	distance  int
	longest   int
	highlight map[int]bool
}

func displayChars(line []byte, chars []int) {
	for _, pos := range chars {
		line[pos] = strings.ToUpper(string(line[pos]))[0]
	}
}

func (m *Matcher) ToString(line, input string) string {
	nl := []byte(line)
	displayChars(nl[:], m.chars)
	return fmt.Sprintf("linepos=%2d dist=%2d inputpos=%2d lmp=%2d llg=%2d full=%v ended=%v groups=%d chars=%v (%v)",
		m.linepos, m.total_distance, m.inputpos, m.lastmatchpos, m.lenlongestgroup, m.complete_match, m.ended, m.groups, string(nl), m.chars)
}

// Return new pointer to  Matcher with relevant variables copied from the
// original Matcher.
func (m *Matcher) Clone() *Matcher {
	var chars []int

	if len(m.chars) > 0 {
		chars = make([]int, len(m.chars))
		copy(chars, m.chars[:len(m.chars)])
	}

	return &Matcher{
		linepos:        m.linepos,
		inputpos:       m.inputpos,
		groups:         m.groups,
		total_distance: m.total_distance,
		chars:          chars,
		lastmatchpos:   m.lastmatchpos,
	}
}

func charsToMap(chars []int) map[int]bool {
	highlight := make(map[int]bool, 0)
	for _, idx := range chars {
		highlight[idx] = true
	}
	return highlight
}

func score2(line string, input string) (finalscore int, finalh map[int]bool) { //(best *BestScore) {
	// fmt.Println(line, input)
	matchers := make([]*Matcher, 0, 1)
	x := new(Matcher)
	matchers = append(matchers, x)

	if len(input) > len(line) || len(input) == 0 || len(line) == 0 {
		finalscore = -1
		return
	}

Outer:
	for {
		for _, m := range matchers {
			// Skip this matcher?
			if m.ended {
				continue
			}

			// Brasil Bra
			if line[m.linepos] == input[m.inputpos] {
				// VERY SLOW
				// if strings.EqualFold(string(line[m.linepos]), string(input[m.inputpos])) {

				// New matcher to find alternatives
				if (m.linepos + 1) < len(line) { // If matcher isn't starting beyond the current line...
					new_matcher := m.Clone() // Alternate matcher is like the current
					new_matcher.linepos++    // but skip current char

					matchers = append(matchers, new_matcher)
				}

				m.chars = append(m.chars, m.linepos)

				// Advance...
				if m.groups == 0 {
					m.groups = 1
					m.curgrouplen = 1
				} else {
					distlast := m.linepos - m.lastmatchpos
					if distlast > 1 {
						m.groups++
						m.curgrouplen = 1
					} else {
						m.curgrouplen++
						if m.lenlongestgroup < m.curgrouplen {
							m.lenlongestgroup = m.curgrouplen
						}
					}
					m.total_distance += (distlast - 1)
				}

				m.lastmatchpos = m.linepos
				m.inputpos++
				m.linepos++

				// Input has matched fully?
				if m.inputpos == len(input) {
					m.ended = true
					m.complete_match = true
				}

				if m.linepos == len(line) {
					m.ended = true
				}

			} else {
				// No match for current char
				m.linepos++

				if m.linepos == len(line) {
					m.ended = true
				}
			}
		}

		all_ended := true
		for _, m := range matchers {
			all_ended = all_ended && m.ended
		}
		if all_ended {
			break Outer
		}
	}

	best := new(BestScore)

	first := true
	for _, m := range matchers {
		if !m.complete_match {
			continue
		}

		score := 100 + (len(input)-m.groups)*10 - m.total_distance + m.lenlongestgroup

		if first || score >= best.score { //m.groups < best.groups || (m.groups == best.groups && m.total_distance < best.distance) {
			first = false
			best.highlight = charsToMap(m.chars)
			best.groups = m.groups
			best.distance = m.total_distance
			best.score = score

			finalscore = score
			finalh = best.highlight
		}
	}

	return
}
