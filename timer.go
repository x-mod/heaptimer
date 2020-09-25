package heaptimer

import (
	"container/heap"
	"context"
	"fmt"
	"sync"
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
	mu       sync.Mutex
}

type Opt func(*Timer)

func Duration(d time.Duration) Opt {
	return func(tm *Timer) {
		tm.duration = d
	}
}

func New(opts ...Opt) *Timer {
	tm := &Timer{
		C:        make(chan interface{}),
		heap:     NewHeap(),
		duration: time.Millisecond * 500,
		close:    event.New(),
		stopped:  event.New(),
	}
	for _, opt := range opts {
		opt(tm)
	}
	tm.timer = time.NewTimer(tm.duration)
	return tm
}

func (tm *Timer) Pop() (interface{}, bool) {
	val, ok := <-tm.C
	return val, ok
}

func (tm *Timer) Drain() interface{} {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	if tm.heap.Len() > 0 {
		node := heap.Pop(tm.heap).(*Node)
		return node.value
	}
	return nil
}

func (tm *Timer) Push(val interface{}, t time.Time) {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	node := &Node{value: val, tm: t}
	heap.Push(tm.heap, node)
	if head := tm.heap.Head(); head.tm.After(time.Now()) {
		tm.timer.Reset(head.tm.Sub(time.Now()))
	}
}

func (tm *Timer) PushWithDuration(val interface{}, duration time.Duration) {
	tm.Push(val, time.Now().Add(duration))
}

func (tm *Timer) Serve(ctx context.Context) (err error) {
	defer func() {
		tm.stopped.Fire()
		if rc := recover(); rc != nil {
			switch rv := rc.(type) {
			case error:
				err = rv.(error)
				return
			default:
				err = fmt.Errorf("%v", rv)
				return
			}
		}
	}()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-tm.close.Done():
			return nil
		case _, ok := <-tm.timer.C:
			if !ok { // timer closed
				return nil
			}
			d := tm.duration
			if nxt, ok := tm.next(); ok {
				d = nxt
			}
			tm.timer.Reset(d)
		}
	}
}

func (tm *Timer) next() (time.Duration, bool) {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	head := tm.heap.Head()
	for head != nil && head.tm.Before(time.Now()) {
		node := heap.Pop(tm.heap).(*Node)
		//panic when tm.C closed
		tm.C <- node.value
		head = tm.heap.Head()
	}
	if head != nil {
		return head.tm.Sub(time.Now()), true
	}
	return 0, false
}

func (tm *Timer) Close() <-chan struct{} {
	close(tm.C)
	tm.timer.Stop()
	tm.close.Fire()
	return tm.stopped.Done()
}
