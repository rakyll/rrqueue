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

package rrqueue

import (
	"errors"
	"sync"
	"time"
)

var (
	ErrQueueEmpty = errors.New("queue is empty")
)

type RRQueue struct {
	numPr  int
	queues [][]interface{}
	locks  []sync.Mutex

	TickInterval time.Duration
	Timeout      time.Duration
	Fn           func(interface{})
}

func New(n int) *RRQueue {
	return &RRQueue{
		TickInterval: time.Microsecond,
		numPr:        n,
		queues:       make([][]interface{}, n, n),
		locks:        make([]sync.Mutex, n, n),
	}
}

func (p *RRQueue) Enqueue(pr int, item interface{}) error {
	p.locks[pr].Lock()
	defer p.locks[pr].Unlock()

	p.queues[pr] = append(p.queues[pr], item)
	return nil
}

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
