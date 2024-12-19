package util

import "cmp"

type GenHeapItem[K cmp.Ordered, V any] struct {
	Key K
	Value V
}

type GenHeap[K cmp.Ordered, V any] []*GenHeapItem[K, V]

func (h *GenHeap[K, V]) siftUp(idx int) {
	for idx > 0 {
		cur, parent := (*h)[idx], (*h)[idx/2]
		if parent.Key < cur.Key {
			break
		}
		(*h)[idx], (*h)[idx/2] = parent, cur
		idx /= 2
	}
}

func (h *GenHeap[K, V]) siftDown(idx int) {
	for 2*idx <= len(*h)-1 {
		left, right, smallest := 2*idx, 2*idx+1, idx

		if left <= len(*h)-1 && (*h)[left].Key < (*h)[smallest].Key {
			smallest = left
		}

		if right <= len(*h)-1 && (*h)[right].Key < (*h)[smallest].Key {
			smallest = right
		}

		if smallest == idx {
			break
		}

		(*h)[idx], (*h)[smallest] = (*h)[smallest], (*h)[idx]
		idx = smallest
	}
}

func (h *GenHeap[K, V]) Push(key K, value V) {
	*h = append(*h, &GenHeapItem[K, V]{
		Key: key,
		Value: value,
	})
	h.siftUp(len(*h)-1)
}

func (h *GenHeap[K, V]) Pop() *GenHeapItem[K, V] {
	if len(*h) == 0 {
		panic("Attempted to pop empty heap")
	}
	res := (*h)[0]
	(*h)[0] = (*h)[len(*h)-1]
	*h = (*h)[:len(*h)-1]
	h.siftDown(0)
	return res
}

func (h *GenHeap[K, V]) Peek() (*GenHeapItem[K, V], bool) {
	if len(*h) == 0 {
		return nil, false
	}
	return (*h)[0], true
}

func (h *GenHeap[K, V]) Len() int {
	return len(*h)
}
