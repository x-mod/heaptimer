package heaptimer

import (
	"container/heap"
	"context"
	"sync"
	"time"

	"github.com/x-mod/event"
)

type HeapTimer struct {
	heap    *Heap
	mu      sync.Mutex
	serving *event.Event
	stopped *event.Event
	putting *event.Event
	ticker  *time.Ticker
	close   chan struct{}
	V       chan interface{}
}

func New() *HeapTimer {
	return &HeapTimer{
		heap:    NewHeap(),
		mu:      sync.Mutex{},
		serving: event.New(),
		stopped: event.New(),
		putting: event.New(),
		ticker:  nil,
		close:   make(chan struct{}),
		V:       make(chan interface{}),
	}
}

func (ht *HeapTimer) Push(val interface{}, tm time.Time) {
	ht.push(&Node{Value: val, ExpireAt: tm})
	ht.putting.Fire()
}

func (ht *HeapTimer) PushWithDuration(val interface{}, d time.Duration) {
	ht.Push(val, time.Now().Add(d))
}

func (ht *HeapTimer) push(n *Node) {
	ht.mu.Lock()
	defer ht.mu.Unlock()
	heap.Push(ht.heap, n)
}

func (ht *HeapTimer) pop() *Node {
	ht.mu.Lock()
	defer ht.mu.Unlock()
	if ht.heap.Len() > 0 {
		n := heap.Pop(ht.heap)
		return n.(*Node)
	}
	return nil
}

func (ht *HeapTimer) Pop() (interface{}, bool) {
	<-ht.serving.Done()
	val, ok := <-ht.V
	return val, ok
}

func (ht *HeapTimer) reschedule() {
	ht.mu.Lock()
	defer ht.mu.Unlock()
	ht.putting = event.New()
	if ht.ticker == nil {
		ht.ticker = time.NewTicker(time.Second * 10)
	}
	if n := ht.heap.Head(); n != nil {
		if time.Until(n.ExpireAt) > 0 {
			ht.ticker.Reset(time.Until(n.ExpireAt))
		} else {
			ht.ticker.Reset(time.Millisecond * 10)
		}
	}
}
func (ht *HeapTimer) Serve(ctx context.Context) error {
	defer ht.stopped.Fire()
	defer close(ht.V)

	ht.serving.Fire()
	for {
		ht.reschedule()
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ht.close:
			return nil
		case <-ht.putting.Done():
			continue
		case <-ht.ticker.C:
			if n := ht.pop(); n != nil {
				ht.V <- n.Value
			}
		}
	}
}

func (ht *HeapTimer) Serving() <-chan struct{} {
	return ht.serving.Done()
}

func (ht *HeapTimer) Close() <-chan struct{} {
	if ht.serving.HasFired() {
		close(ht.close)
		return ht.stopped.Done()
	}
	return event.Done()
}
