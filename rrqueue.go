// Copyright 2013 Google Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package rrqueue provides a container for priority queues
// and a simple round-robin scheduled consumer.
package rrqueue

import (
	"errors"
	"sync"
	"time"
)

var (
	// ErrQueueEmpty represents a failed Dequeue because
	// there are not queued items available.
	ErrQueueEmpty = errors.New("queue is empty")
)

type RRQueue struct {
	numPr  int
	queues [][]interface{}
	locks  []sync.Mutex

	// TickInterval represents how frequently
	// consumer consumer will tick.
	TickInterval time.Duration
	// Timeout represents the max duration a consumer
	// function should consume on an item.
	Timeout time.Duration
	// Fn represents a function that will consume
	// queued items.
	Fn func(interface{})
}

// New returns a new rr queue instance.
func New(n int) *RRQueue {
	return &RRQueue{
		TickInterval: time.Microsecond,
		numPr:        n,
		queues:       make([][]interface{}, n, n),
		locks:        make([]sync.Mutex, n, n),
	}
}

// Enqueue puts the given item to its priority queue
func (p *RRQueue) Enqueue(pr int, item interface{}) error {
	p.locks[pr].Lock()
	defer p.locks[pr].Unlock()

	p.queues[pr] = append(p.queues[pr], item)
	return nil
}

// Dequeue retrieves an items from the given priority
// queue. It returns ErrQueueEmpty if the priorty queue
// has no elements.
func (p *RRQueue) Dequeue(pr int) (interface{}, error) {
	p.locks[pr].Lock()
	defer p.locks[pr].Unlock()

	q := p.queues[pr]
	if q == nil || len(q) == 0 {
		return nil, ErrQueueEmpty
	}
	// TODO: benchmark
	item := q[0].(interface{})
	p.queues[pr] = q[1:]
	return item, nil
}

// Start initiates the round robin processing, runs forever
// to consume queued items. It waits TickInterval before
// starting a new pass.
func (p *RRQueue) Start() {
	done := make(chan bool, 1)
	go func() {
		for {
			p.rrtick()
			<-time.After(p.TickInterval)
		}
	}()
	<-done
}

func (p *RRQueue) Stop() {
	panic("not implemented")
}

func (p *RRQueue) rrtick() {
	done := make(chan bool, p.numPr)
	for i := 0; i < p.numPr; i++ {
		go p.process(done, i)
	}
	// TODO: handle timeout, etc
	<-done
}

func (p *RRQueue) process(done chan<- bool, pr int) {
	item, _ := p.Dequeue(pr)
	if item != nil && p.Fn != nil {
		p.Fn(item)
	}
	done <- true
}
