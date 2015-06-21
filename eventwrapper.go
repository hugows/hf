package main

import (
	"sync/atomic"

	"github.com/nsf/termbox-go"
)

type TermboxEventWrapper struct {
	queue chan termbox.Event
	count int32
}

func NewTermboxEventWrapper() *TermboxEventWrapper {
	t := new(TermboxEventWrapper)

	t.queue = make(chan termbox.Event, 1)

	go func() {
		for {
			t.queue <- termbox.PollEvent()
			atomic.AddInt32(&t.count, 1)
		}
	}()
	return t
}

func (t *TermboxEventWrapper) Peek() bool {
	return atomic.LoadInt32(&t.count) > 0
}

func (t *TermboxEventWrapper) Poll() termbox.Event {
	atomic.AddInt32(&t.count, -1)
	ev := <-t.queue
	return ev
}
