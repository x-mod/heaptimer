package heaptimer

import (
	"container/heap"
	"context"
	"time"

	"github.com/x-mod/event"
)

type Timer struct {
	C        chan interface{}
	heap     *Heap
	duration time.Duration
	timer    *time.Timer
	close    *event.Event
	stopped  *event.Event
}

func New(duration time.Duration) *Timer {
	return &Timer{
		C:        make(chan interface{}),
		heap:     NewHeap(),
		duration: duration,
		timer:    time.NewTimer(duration),
		close:    event.New(),
		stopped:  event.New(),
	}
}

func (tm *Timer) Pop() (interface{}, bool) {
	val, ok := <-tm.C
	return val, ok
}

func (tm *Timer) Drain() interface{} {
	if tm.heap.Len() > 0 {
		node := heap.Pop(tm.heap).(*Node)
		return node.value
	}
	return nil
}

func (tm *Timer) Push(val interface{}, duration time.Duration) {
	node := &Node{value: val, tm: time.Now().Add(duration)}
	tm.heap.Push(node)
	head := tm.heap.Head()
	if head.tm.After(time.Now()) {
		tm.timer.Reset(head.tm.Sub(time.Now()))
	}
}

func (tm *Timer) PushByTime(val interface{}, t time.Time) {
	node := &Node{value: val, tm: t}
	tm.heap.Push(node)
	head := tm.heap.Head()
	if head.tm.After(time.Now()) {
		tm.timer.Reset(head.tm.Sub(time.Now()))
	}
}

func (tm *Timer) Serve(ctx context.Context) error {
	defer tm.stopped.Fire()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-tm.close.Done():
			return nil
		case <-tm.timer.C:
			d := tm.duration
			if head := tm.heap.Head(); head != nil {
				if head.tm.Before(time.Now()) {
					node := heap.Pop(tm.heap).(*Node)
					tm.C <- node.value
				} else {
					d = head.tm.Sub(time.Now())
				}
			}
			tm.timer.Reset(d)
		}
	}
}

func (tm *Timer) Close() <-chan struct{} {
	close(tm.C)
	tm.timer.Stop()
	tm.close.Fire()
	return tm.stopped.Done()
}
