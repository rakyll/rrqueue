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
	"testing"
)

func TestNumPriorityQueues(t *testing.T) {
	q := New(4)
	if len(q.queues) != 4 {
		t.Error("New: can't create correct number of priority queues")
	}
}

func TestEnqueue(t *testing.T) {
	q := New(5)
	q.Enqueue(3, "item")

	if q.queues[3][0] != "item" {
		t.Error("Enqueue: can't put the item to the right queue")
	}
}

func TestDequeue(t *testing.T) {
	q := New(5)
	q.Enqueue(0, "P0 item")
	q.Enqueue(1, "P1 item")
	q.Enqueue(1, "Another P1 item")

	item, err := q.Dequeue(1)
	if err != nil || item != "P1 item" {
		t.Error("Dequeue: can't dequeue the right element")
	}
}
