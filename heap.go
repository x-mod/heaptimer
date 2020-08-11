package heaptimer

import (
	"container/heap"
	"time"
)

type Node struct {
	value interface{}
	tm    time.Time
}

type Heap []*Node

func NewHeap() *Heap {
	container := &Heap{}
	heap.Init(container)
	return container
}

func (container Heap) Len() int { return len(container) }

func (container Heap) Less(i, j int) bool {
	return container[i].tm.Before(container[j].tm)
}

func (container Heap) Swap(i, j int) {
	container[i], container[j] = container[j], container[i]
}

func (container *Heap) Push(x interface{}) {
	item := x.(*Node)
	*container = append(*container, item)
}

func (container *Heap) Pop() interface{} {
	old := *container
	n := len(old)
	x := old[n-1]
	*container = old[0 : n-1]
	return x
}

func (container *Heap) Head() *Node {
	n := len(*container)
	if n > 0 {
		return (*container)[0]
	}
	return nil
}

func (container *Heap) Tail() *Node {
	n := len(*container)
	if n > 0 {
		return (*container)[n-1]
	}
	return nil
}
