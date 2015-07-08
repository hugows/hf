package main

import (
	"fmt"
	"sync"
	"time"
)

type Stats struct {
	sync.RWMutex
	start    time.Time
	counters map[string]int
}

func NewStats() *Stats {
	return &Stats{
		counters: make(map[string]int),
		start:    time.Now(),
	}
}

func (s *Stats) Inc(key string) {
	s.Lock()
	s.counters[key]++
	s.Unlock()
}

func (s *Stats) Print() {
	s.RLock()
	fmt.Println("*** stats - elapsed time was", time.Since(s.start),"***")
	for k, v := range s.counters {
		fmt.Printf("%5d call(s) %s\n", v, k)
	}
	s.RUnlock()
}
