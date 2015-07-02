package main

import (
	"fmt"
	"testing"
)

func TestOccur(t *testing.T) {
	cases := []struct {
		against, userinput, expected string
	}{
		{"brasil", "bra", "[[0] [1] [2]]"},
		{"brasil", "bRa", "[[0] [1] [2]]"},
		{"brasil", "sil", "[[3] [4] [5]]"},
		{"abba", "ab", "[[0 3] [1 2]]"},
		{"a", "a", "[[0]]"},
		{"abc", "", "[]"},
		{"aaa", "a", "[[0 1 2]]"},
		{"aaa", "A", "[[0 1 2]]"},
		{"abc", "x", "[]"},
		{"", "", "[]"},
	}
	for _, c := range cases {
		oc, _ := occur(c.against, c.userinput)
		soc := fmt.Sprint(oc)
		if soc != c.expected {
			t.Errorf("Occur(%q,%q) == %q, want %q", c.against, c.userinput, soc, c.expected)
		}
	}
}

func TestCombinations(t *testing.T) {
	cases := []struct {
		against, userinput, expected string
	}{
		{"abba", "ab", "[[0 1] [0 2] [3 1] [3 2]]"},
		{"ABCC", "abc", "[[0 1 2] [0 1 3]]"},
		{"aabbbcc", "abc", "[[0 2 5] [0 2 6] [0 3 5] [0 3 6] [0 4 5] [0 4 6] [1 2 5] [1 2 6] [1 3 5] [1 3 6] [1 4 5] [1 4 6]]"},
	}
	for _, c := range cases {
		all := make([][]int, 0)
		oc, _ := occur(c.against, c.userinput)
		combinations([]int{}, oc, &all)
		sall := fmt.Sprint(all)
		if sall != c.expected {
			t.Errorf("Combinations(%q,%q) == %q, want %q", c.against, c.userinput, sall, c.expected)
		}
	}
}

func TestScore(t *testing.T) {
	cases := []string{
		"zip",
		"azipa",
		"azxipo",
		"zixp",
		"zxixp",
		"zkkkkkkkkip",
		"zkkkkiidjasijdsap",
	}

	last := 0
	for _, c := range cases {
		best := score(c, "zip")
		s := best.score
		if s < 0 || s > 100 || (s < last) {
			t.Errorf("Score(%q,%q) == %d (previous %d)", c, "zip", s, last)
		}
		last = s
	}
}
