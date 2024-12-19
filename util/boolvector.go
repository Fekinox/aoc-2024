package util

type BoolVector []int64

func NewBoolVector(size int) BoolVector {
	q, r := size / 64, size % 64
	if r > 0 {
		q += 1
	}

	return make([]int64, q)
}

func (b *BoolVector) Get(i int) bool {
	q, r := i / 64, i % 64
	return (*b)[q] & (1 << r) > 0
}

func (b *BoolVector) Set(i int) {
	q, r := i / 64, i % 64
	(*b)[q] |= (1 << r)
}

func (b *BoolVector) Unset(i int) {
	q, r := i / 64, i % 64
	(*b)[q] &^= (1 << r)
}

func (b *BoolVector) Toggle(i int) {
	q, r := i / 64, i % 64
	(*b)[q] ^= (1 << r)
}
