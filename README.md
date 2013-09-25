# rrqueue

`rrqueue` is a priority queue implementation with round-robin like scheduling to consume the enqueued items. It retrieves a single item from each priority queue on each tick and process the item with the provided func. Equal time slices are allocated for each queue on each tick, you're expected to put less items to high priority queues to avoid starvation. 

    import (
        "github.com/rakyll/rrqueue"
    )

    // Create a new rrqueue with 5 priority queues
    q := rrqueue.New(5)
    
    // Optionally set how often consumer should tick
    q.TickInterval = time.Millisecond
    
    // Set a function to process queued items
    q.Fn = func(item interface{}) {
      log.Println(item)
    }

    // Enqueue some items
    q.Enqueue(P0, "some p0 item")
    q.Enqueue(P1, "a p1 item")
    q.Enqueue(P1, "another p1 item")

    // Start to consume
    q.Start()
    
## License
    Copyright 2013 Google Inc. All Rights Reserved.
    
    Licensed under the Apache License, Version 2.0 (the "License");
    you may not use this file except in compliance with the License.
    You may obtain a copy of the License at
    
         http://www.apache.org/licenses/LICENSE-2.0
    
    Unless required by applicable law or agreed to in writing, software
    distributed under the License is distributed on an "AS IS" BASIS,
    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
    See the License for the specific language governing permissions and
    limitations under the License.
