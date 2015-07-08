package main

import (
	"fmt"
	"sync"
	"time"
)

type Stats struct {
	sync.Mutex
	start    time.Time
	counters map[string]int32
}

func NewStats() *Stats {
	return &Stats{
		counters: make(map[string]int32),
		start:    time.Now(),
	}
}

func (s *Stats) Inc(key string) {
	s.Lock()
	val, ok := s.counters[key]
	if ok {
		s.counters[key] = val + 1
	} else {
		s.counters[key] = 1
	}
	s.Unlock()
}

func (s *Stats) Print() {
	s.Lock()
	fmt.Println("*** stats *** - program ran for", time.Since(s.start))
	for k, v := range s.counters {
		fmt.Printf("%6d call(s) %s\n", v, k)
	}
	s.Unlock()
}
