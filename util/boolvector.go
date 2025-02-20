package util

type BoolVector struct {
	data []int64
	len int
}

func NewBoolVector(size int) *BoolVector {
	q, r := size / 64, size % 64
	if r > 0 {
		q += 1
	}

	return &BoolVector{
		data: make([]int64, q),
		len: size,
	}
}

func (b *BoolVector) Get(i int) bool {
	q, r := i / 64, i % 64
	return b.data[q] & (1 << r) > 0
}

func (b *BoolVector) Set(i int) {
	q, r := i / 64, i % 64
	b.data[q] |= (1 << r)
}

func (b *BoolVector) Unset(i int) {
	q, r := i / 64, i % 64
	b.data[q] &^= (1 << r)
}

func (b *BoolVector) Toggle(i int) {
	q, r := i / 64, i % 64
	b.data[q] ^= (1 << r)
}

func (b *BoolVector) Len() int {
	return b.len
}

func (b *BoolVector) Clone() *BoolVector {
	res := NewBoolVector(b.len)
	for i := range res.data {
		res.data[i] = b.data[i]
	}
	return res
}

func (b *BoolVector) Union(other *BoolVector) {
	if other == nil {
		return
	}
	if b.len != other.len {
		return
	}

	for i := range b.data {
		b.data[i] |= other.data[i]
	}
}

func (b *BoolVector) Intersection(other *BoolVector) {
	if other == nil {
		for i := range b.data {
			b.data[i] = 0 
		}
		return
	}
	if b.len != other.len {
		return
	}

	for i := range b.data {
		b.data[i] &= other.data[i]
	}
}

func (b *BoolVector) Difference(other *BoolVector) {
	if other == nil {
		return
	}
	if b.len != other.len {
		return
	}

	for i := range b.data {
		b.data[i] &^= other.data[i]
	}
}

func (b *BoolVector) Xor(other *BoolVector) {
	if b.len != other.len {
		return
	}

	for i := range b.data {
		b.data[i] ^= other.data[i]
	}
}

func (b *BoolVector) Invert() {
	for i := range b.data {
		b.data[i] = ^b.data[i]
	}
}
