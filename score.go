package main

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

func (m *Matcher) Hash() uint32 {
	return uint32(m.linepos)<<16 | uint32(m.inputpos)
}

func charsToMap(chars []int) map[int]bool {
	highlight := make(map[int]bool, 0)
	for _, idx := range chars {
		highlight[idx] = true
	}
	return highlight
}

func score(when int64, line string, input string) (best *BestScore) {
	best = new(BestScore)

	if len(input) > len(line) || len(input) == 0 || len(line) == 0 {
		best.score = -1
		return
	}

	matchers := make([]*Matcher, 0, 1)
	matchersmap := make(map[uint32]bool)

	x := new(Matcher)
	matchers = append(matchers, x)

Outer:
	for {
		for _, m := range matchers {
			// Skip this matcher?
			if m.ended {
				continue
			}

			if line[m.linepos] == input[m.inputpos] {
				// New matcher to find alternatives
				if (m.linepos + 1) < len(line) { // If matcher isn't starting beyond the current line...
					new_matcher := m.Clone() // Alternate matcher is like the current
					new_matcher.linepos++    // but skip current char

					// Dont add repeated matchers
					if _, seen := matchersmap[new_matcher.Hash()]; !seen {
						matchers = append(matchers, new_matcher)
						matchersmap[new_matcher.Hash()] = true
					}
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

					if global_lastkeypress > when {
						// return with first complete since our deadline is past!
						return
					}
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

	if global_lastkeypress > when {
		best.score = -1
		return
	}

	first := true
	for _, m := range matchers {
		if !m.complete_match {
			continue
		}

		score := 100 + (len(input)-m.groups)*10 - m.total_distance + m.lenlongestgroup

		if first || score > best.score {
			first = false
			best.highlight = charsToMap(m.chars)
			best.groups = m.groups
			best.distance = m.total_distance
			best.score = score
		}
	}

	return
}
