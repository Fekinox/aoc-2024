package util

type Grid[T any] struct {
	data   []T
	Width  int
	Height int
}

type GridItem[T any] struct {
	Value T
	X int
	Y int
}

func MakeGrid[T any](width, height int, def T) Grid[T] {
	width = max(width, 0)
	height = max(height, 0)
	data := make([]T, width*height)
	for i := 0; i < width*height; i++ {
		data[i] = def
	}

	return Grid[T]{
		data:   data,
		Width:  width,
		Height: height,
	}
}

func MakeGridWith[T any](width, height int, gen func(x, y int) T) Grid[T] {
	width = max(width, 0)
	height = max(height, 0)
	data := make([]T, width*height)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			data[y*width+x] = gen(x, y)
		}
	}

	return Grid[T]{
		data:   data,
		Width:  width,
		Height: height,
	}
}

func GridFromSlice[T any](width, height int, elems ...T) Grid[T] {
	width = max(width, 0)
	height = max(height, 0)
	data := make([]T, width*height)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			data[y*width+x] = elems[y*width+x]
		}
	}

	return Grid[T]{
		data:   data,
		Width:  width,
		Height: height,
	}
}

func GridFromSlices[T any](slices ...[]T) Grid[T] {
	height := len(slices)
	var width int
	if len(slices) != 0 {
		width = len(slices[0])
	}
	data := make([]T, width*height)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			data[y*width+x] = slices[y][x]
		}
	}

	return Grid[T]{
		data:   data,
		Width:  width,
		Height: height,
	}
}

func GridFromStrings(strings ...string) Grid[rune] {
	slices := make([][]rune, len(strings))
	for i := 0; i < len(strings); i++ {
		slices[i] = []rune(strings[i])
	}

	return GridFromSlices(slices...)
}

func (g *Grid[T]) InBounds(x, y int) bool {
	return x >= 0 && x < g.Width && y >= 0 && y < g.Height
}

func (g *Grid[T]) Get(x int, y int) (T, bool) {
	if !g.InBounds(x, y) {
		return *new(T), false
	}

	return g.data[y*g.Width+x], true
}

func (g *Grid[T]) MustGet(x, y int) T {
	if !g.InBounds(x, y) {
		panic("Out of bounds")
	}
	return g.data[y*g.Width+x]
}

func (g *Grid[T]) Set(x, y int, val T) bool {
	if !g.InBounds(x, y) {
		return false
	}

	g.data[y*g.Width+x] = val

	return true
}

func (g *Grid[T]) Resize(ox, oy int, neww, newh int, def T) Grid[T] {
	return MakeGridWith(neww, newh, func(x, y int) T {
		xx := x - ox
		yy := y - oy
		val, ok := g.Get(xx, yy)
		if ok {
			return val
		} else {
			return def
		}
	})
}

func (g *Grid[T]) ShallowClone() Grid[T] {
	return MakeGridWith(g.Width, g.Height, func(x, y int) T {
		return g.MustGet(x, y)
	})
}

func (g *Grid[T]) GetPoint(p Point) (T, bool) {
	return g.Get(int(p.X), int(p.Y))
}

func (g *Grid[T]) MustGetPoint(p Point) T {
	return g.MustGet(int(p.X), int(p.Y))
}

func (g *Grid[T]) SetPoint(p Point, val T) bool {
	return g.Set(int(p.X), int(p.Y), val)
}

func (g *Grid[T]) NeighborhoodInclusive(x, y int, radius int) []GridItem[T] {
	if radius <= 0 {
		v, ok := g.Get(x, y)
		if !ok {
			return []GridItem[T]{}
		}
		return []GridItem[T]{
			{
				Value: v,
				X: x,
				Y: y,
			},
		}
	}
	res := []GridItem[T]{}
	yLo := max(y - radius, 0)
	yHi := min(y + radius, g.Height-1)
	xLo := max(x - radius, 0)
	xHi := min(x + radius, g.Width-1)

	for yy := yLo; yy <= yHi; yy++ {
		for xx := xLo; xx <= xHi; xx++ {
			res = append(res, GridItem[T]{
				Value: g.MustGet(xx, yy),
				X: xx,
				Y: yy,
			})
		}
	}

	return res
}

func (g *Grid[T]) NeighborhoodExclusive(x, y int, radius int) []GridItem[T] {
	if radius <= 0 {
		return []GridItem[T]{}
	}

	res := []GridItem[T]{}
	yLo := max(y - radius, 0)
	yHi := min(y + radius, g.Height-1)
	xLo := max(x - radius, 0)
	xHi := min(x + radius, g.Width-1)

	for yy := yLo; yy <= yHi; yy++ {
		for xx := xLo; xx <= xHi; xx++ {
			if xx == x && yy == y {
				continue
			}
			res = append(res, GridItem[T]{
				Value: g.MustGet(xx, yy),
				X: xx,
				Y: yy,
			})
		}
	}

	return res
}

func (g *Grid[T]) Region(ox, oy int, w, h int) Grid[T] {
	oxMin := max(ox, 0)
	oxMax := min(ox + w, g.Width)
	oyMin := max(oy, 0)
	oyMax := min(oy + h, g.Height)
	return MakeGridWith(oxMax - oxMin, oyMax - oyMin, func(x, y int) T {
		return g.MustGet(x - oxMin, y - oyMin)
	})
}

func (g *Grid[T]) RegionValues(ox, oy int, w, h int) []GridItem[T] {
	res := []GridItem[T]{}

	yLo := max(oy, 0)
	yHi := min(oy + h, g.Height)
	xLo := max(ox, 0)
	xHi := min(ox + w, g.Width)

	for yy := yLo; yy < yHi; yy++ {
		for xx := xLo; xx < xHi; xx++ {
			res = append(res, GridItem[T]{
				Value: g.MustGet(xx, yy),
				X: xx,
				Y: yy,
			})
		}
	}

	return res
}

