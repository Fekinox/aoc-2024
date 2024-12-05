package util

type Point struct {
	X int64
	Y int64
}

var (
	N  = Point{X: 0, Y: -1}
	NE = Point{X: 1, Y: -1}
	E  = Point{X: 1, Y: 0}
	SE = Point{X: 1, Y: 1}
	S  = Point{X: 0, Y: 1}
	SW = Point{X: -1, Y: 1}
	W  = Point{X: -1, Y: 0}
	NW = Point{X: -1, Y: -1}

	CardinalDirections = []Point{N, E, S, W}
	AllDirections = []Point{N, NE, E, SE, S, SW, W, NW}
)

func (p Point) Add(other Point) Point {
	return Point{
		X: p.X + other.X,
		Y: p.Y + other.Y,
	}
}

func (p Point) Sub(other Point) Point {
	return Point{
		X: p.X - other.X,
		Y: p.Y - other.Y,
	}
}
