package util

import (
	"cmp"
	"container/heap"
)

type PriorityQueueItem[K cmp.Ordered, V any] struct {
	Value V
	Priority K
	Index int
}

type PriorityQueue[K cmp.Ordered, V any] []*PriorityQueueItem[K, V]

var (
	_ heap.Interface = &PriorityQueue[int, int]{}
)

func (pq PriorityQueue[K, V]) Len() int {
	return len(pq)
}

func (pq PriorityQueue[K, V]) Less(i, j int) bool {
	return cmp.Compare(pq[i].Priority, pq[j].Priority) < 0
}

func (pq PriorityQueue[K, V]) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].Index = i
	pq[j].Index = j
}

func (pq *PriorityQueue[K, V]) Push(x any) {
	n := len(*pq)
	item := x.(*PriorityQueueItem[K, V])
	item.Index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue[K, V]) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	item.Index = -1
	*pq = old[0:n-1]
	return item
}

func (pq *PriorityQueue[K, V]) Update(item *PriorityQueueItem[K, V], value V, priority K) {
	item.Value = value
	item.Priority = priority
	heap.Fix(pq, item.Index)
}
