package main

import "sync/atomic"

type RoundRobin struct {
	counter uint64
	size    uint64
}

func NewRoundRobin(minimum, maximum uint64) *RoundRobin {
	return &RoundRobin{
		counter: minimum,
		size:    maximum,
	}
}

func (rr *RoundRobin) Next() uint64 {
	return atomic.AddUint64(&rr.counter, 1) % (rr.size + 1)
}
