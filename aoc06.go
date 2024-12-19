package main

import (
	"fmt"
	"io"
	"math"
	"sync"
	"sync/atomic"

	"github.com/Fekinox/aoc-2024/util"
)

func Problem06Simple(in io.Reader, out io.Writer) {
	var silver, gold int64

	var sx, sy int
	g := util.GridFromReader(in)
	for y := range g.Height {
		for x := range g.Width {
			if g.MustGet(x, y) == '^' {
				sx = x
				sy = y
			}
		}
	}

	type pos struct {
		X int
		Y int
	}

	type rec struct {
		X  int
		Y  int
		DX int
		DY int
	}

	loopsIfBlocked := func(x, y int, d util.Point) bool {
		visited := make(map[rec]struct{})
		ox, oy := x+int(d.X), y+int(d.Y)

		for {
			if _, ok := visited[rec{X: x, Y: y, DX: int(d.X), DY: int(d.Y)}]; ok {
				return true
			}

			visited[rec{X: x, Y: y, DX: int(d.X), DY: int(d.Y)}] = struct{}{}

			xx, yy := x, y
			for {
				xxx, yyy := xx+int(d.X), yy+int(d.Y)
				p, ok := g.Get(xxx, yyy)
				if !ok {
					return false
				}
				if p == '#' || (xxx == ox && yyy == oy) {
					break
				}
				xx = xxx
				yy = yyy
			}

			x = xx
			y = yy
			d = util.Point{
				X: -d.Y,
				Y: d.X,
			}
		}
	}

	curX, curY, curD := sx, sy, util.N
	for {
		var nx, ny int
		nd := curD

		alreadyVisited := g.MustGet(curX, curY) == 'X'
		next, ok := g.Get(curX+int(curD.X), curY+int(curD.Y))

		if !alreadyVisited {
			g.Set(curX, curY, 'X')
			silver++
		}

		if !ok {
			break
		}

		if next != 'X' && next != '#' && loopsIfBlocked(curX, curY, curD) {
			gold++
		}

		if next == '#' {
			nd = util.Point{
				X: -curD.Y,
				Y: curD.X,
			}
			nx, ny = curX, curY
		} else {
			nx, ny = curX+int(curD.X), curY+int(curD.Y)
		}

		curX, curY, curD = nx, ny, nd
	}

	fmt.Fprintln(out, "silver:", silver)
	fmt.Fprintln(out, "gold:", gold)
}

func Problem06Complex(in io.Reader, out io.Writer) {
	var silver, gold int64
	var cacheMiss, cacheHit int64

	var sx, sy int
	g := util.GridFromReader(in)
	for y := range g.Height {
		for x := range g.Width {
			if g.MustGet(x, y) == '^' {
				sx = x
				sy = y
			}
		}
	}

	type pos struct {
		X int
		Y int
	}

	type rec struct {
		X  int
		Y  int
		DX int
		DY int
	}

	jumpTable := make(map[rec]pos)

	loopsIfBlocked := func(x, y int, d util.Point) bool {
		visited := make(map[rec]struct{})
		ox, oy := x+int(d.X), y+int(d.Y)

		for {
			if _, ok := visited[rec{X: x, Y: y, DX: int(d.X), DY: int(d.Y)}]; ok {
				return true
			}

			visited[rec{X: x, Y: y, DX: int(d.X), DY: int(d.Y)}] = struct{}{}

			xx, yy := x, y
			var tp pos

			if p, ok := jumpTable[rec{X: x, Y: y, DX: int(d.X), DY: int(d.Y)}]; ok {
				cacheHit++
				tp = p
			} else {
				cacheMiss++
				for {
					xxx, yyy := xx+int(d.X), yy+int(d.Y)
					p, ok := g.Get(xxx, yyy)
					if !ok {
						if xxx < 0 {
							tp = pos{
								X: math.MinInt,
								Y: yyy,
							}
						} else if xxx >= g.Width {
							tp = pos{
								X: math.MaxInt,
								Y: yyy,
							}
						} else if yyy < 0 {
							tp = pos{
								X: xxx,
								Y: math.MinInt,
							}
						} else {
							tp = pos{
								X: xxx,
								Y: math.MaxInt,
							}
						}
						break
					}
					if p == '#' {
						tp = pos{
							X: xx,
							Y: yy,
						}
						break
					}
					xx = xxx
					yy = yyy
				}
				jumpTable[rec{X: x, Y: y, DX: int(d.X), DY: int(d.Y)}] = tp
				// bx, by := xx, yy
				// for bx != x || by != y {
				// 	jumpTable[rec{X: bx, Y: by, DX: int(d.X), DY: int(d.Y)}] = tp
				// 	bx -= int(d.X)
				// 	by -= int(d.Y)
				// }
			}

			xx, yy = tp.X, tp.Y
			if d.X == 0 && ox == x {
				if d.Y == -1 && (yy <= oy && oy < y) {
					// up
					yy = oy + 1
				} else if d.Y == 1 && (yy >= oy && oy > y) {
					// down
					yy = oy - 1
				}
			}
			if d.Y == 0 && oy == y {
				if d.X == -1 && (xx <= ox && ox < x) {
					// left
					xx = ox + 1
				} else if d.X == 1 && (xx >= ox && ox > x) {
					// right
					xx = ox - 1
				}
			}

			if !g.InBounds(xx, yy) {
				return false
			}

			x = xx
			y = yy
			d = util.Point{
				X: -d.Y,
				Y: d.X,
			}
		}
	}

	curX, curY, curD := sx, sy, util.N
	for {
		var nx, ny int
		nd := curD

		alreadyVisited := g.MustGet(curX, curY) == 'X'
		next, ok := g.Get(curX+int(curD.X), curY+int(curD.Y))

		if !alreadyVisited {
			g.Set(curX, curY, 'X')
			silver++
		}

		if !ok {
			break
		}

		if next != 'X' && next != '#' && loopsIfBlocked(curX, curY, curD) {
			gold++
		}

		if next == '#' {
			nd = util.Point{
				X: -curD.Y,
				Y: curD.X,
			}
			nx, ny = curX, curY
		} else {
			nx, ny = curX+int(curD.X), curY+int(curD.Y)
		}

		curX, curY, curD = nx, ny, nd
	}

	fmt.Fprintln(out, "silver:", silver)
	fmt.Fprintln(out, "gold:", gold)
	fmt.Fprintf(out, "cache efficiency: %.2f%%\n", 100*float64(cacheHit)/float64(cacheHit+cacheMiss))
}

func Problem06ComplexTH(in io.Reader, out io.Writer) {
	var silver, gold int64

	var sx, sy int
	g := util.GridFromReader(in)
	for y := range g.Height {
		for x := range g.Width {
			if g.MustGet(x, y) == '^' {
				sx = x
				sy = y
			}
		}
	}

	type pos struct {
		X int
		Y int
	}

	type rec struct {
		X  int
		Y  int
		DX int
		DY int
	}

	jumpTable := make(map[rec]pos)

	loopsIfBlocked := func(x, y int, d util.Point) bool {
		tortoise := rec{X: x, Y: y, DX: int(d.X), DY: int(d.Y)}
		hare := rec{X: x, Y: y, DX: int(d.X), DY: int(d.Y)}
		ox, oy := x+int(d.X), y+int(d.Y)

		move := func(r rec) (rec, bool) {
			xx, yy := r.X, r.Y
			var tp pos

			if p, ok := jumpTable[r]; ok {
				tp = p
			} else {
				for {
					xxx, yyy := xx+r.DX, yy+r.DY
					p, ok := g.Get(xxx, yyy)
					if !ok {
						if xxx < 0 {
							tp = pos{
								X: math.MinInt,
								Y: yyy,
							}
						} else if xxx >= g.Width {
							tp = pos{
								X: math.MaxInt,
								Y: yyy,
							}
						} else if yyy < 0 {
							tp = pos{
								X: xxx,
								Y: math.MinInt,
							}
						} else {
							tp = pos{
								X: xxx,
								Y: math.MaxInt,
							}
						}
						break
					}
					if p == '#' {
						tp = pos{
							X: xx,
							Y: yy,
						}
						break
					}
					xx = xxx
					yy = yyy
				}
				jumpTable[r] = tp
			}

			xx, yy = tp.X, tp.Y
			if r.DX == 0 && ox == r.X {
				if r.DY == -1 && (yy <= oy && oy < r.Y) {
					// up
					yy = oy + 1
				} else if r.DY == 1 && (yy >= oy && oy > r.Y) {
					// down
					yy = oy - 1
				}
			}
			if r.DY == 0 && oy == r.Y {
				if r.DX == -1 && (xx <= ox && ox < r.X) {
					// left
					xx = ox + 1
				} else if r.DX == 1 && (xx >= ox && ox > r.X) {
					// right
					xx = ox - 1
				}
			}

			if !g.InBounds(xx, yy) {
				return rec{}, false
			}

			return rec{
				X:  xx,
				Y:  yy,
				DX: -r.DY,
				DY: r.DX,
			}, true
		}

		moveTortoise := false

		for {
			hNext, ok := move(hare)
			if !ok {
				return false
			}
			hare = hNext

			if moveTortoise {
				tortoise, ok = move(tortoise)
				if !ok {
					return false
				}
			}
			moveTortoise = !moveTortoise

			if tortoise == hare {
				return true
			}
		}
	}

	curX, curY, curD := sx, sy, util.N
	for {
		var nx, ny int
		nd := curD

		alreadyVisited := g.MustGet(curX, curY) == 'X'
		next, ok := g.Get(curX+int(curD.X), curY+int(curD.Y))

		if !alreadyVisited {
			g.Set(curX, curY, 'X')
			silver++
		}

		if !ok {
			break
		}

		if next != 'X' && next != '#' && loopsIfBlocked(curX, curY, curD) {
			gold++
		}

		if next == '#' {
			nd = util.Point{
				X: -curD.Y,
				Y: curD.X,
			}
			nx, ny = curX, curY
		} else {
			nx, ny = curX+int(curD.X), curY+int(curD.Y)
		}

		curX, curY, curD = nx, ny, nd
	}

	fmt.Fprintln(out, "silver:", silver)
	fmt.Fprintln(out, "gold:", gold)
}

func Problem06ComplexB(in io.Reader, out io.Writer) {
	var silver, gold int64

	var sx, sy int
	g := util.GridFromReader(in)
	for y := range g.Height {
		for x := range g.Width {
			if g.MustGet(x, y) == '^' {
				sx = x
				sy = y
			}
		}
	}

	type pos struct {
		X int
		Y int
	}

	type rec struct {
		X  int
		Y  int
		DX int
		DY int
	}

	jumpTable := make(map[rec]pos)

	loopsIfBlocked := func(x, y int, d util.Point) bool {
		ox, oy := x+int(d.X), y+int(d.Y)
		move := func(r rec) (rec, bool) {
			xx, yy := r.X, r.Y
			var tp pos

			if p, ok := jumpTable[r]; ok {
				tp = p
			} else {
				for {
					xxx, yyy := xx+r.DX, yy+r.DY
					p, ok := g.Get(xxx, yyy)
					if !ok {
						if xxx < 0 {
							tp = pos{
								X: math.MinInt,
								Y: yyy,
							}
						} else if xxx >= g.Width {
							tp = pos{
								X: math.MaxInt,
								Y: yyy,
							}
						} else if yyy < 0 {
							tp = pos{
								X: xxx,
								Y: math.MinInt,
							}
						} else {
							tp = pos{
								X: xxx,
								Y: math.MaxInt,
							}
						}
						break
					}
					if p == '#' {
						tp = pos{
							X: xx,
							Y: yy,
						}
						break
					}
					xx = xxx
					yy = yyy
				}
				jumpTable[r] = tp
			}

			xx, yy = tp.X, tp.Y
			if r.DX == 0 && ox == r.X {
				if r.DY == -1 && (yy <= oy && oy < r.Y) {
					// up
					yy = oy + 1
				} else if r.DY == 1 && (yy >= oy && oy > r.Y) {
					// down
					yy = oy - 1
				}
			}
			if r.DY == 0 && oy == r.Y {
				if r.DX == -1 && (xx <= ox && ox < r.X) {
					// left
					xx = ox + 1
				} else if r.DX == 1 && (xx >= ox && ox > r.X) {
					// right
					xx = ox - 1
				}
			}

			if !g.InBounds(xx, yy) {
				return rec{}, false
			}

			return rec{
				X:  xx,
				Y:  yy,
				DX: -r.DY,
				DY: r.DX,
			}, true
		}

		// Exponentially search for a cycle with Brent's algorithm
		// i.e. find i and j where f(2^i) = f(j)
		tortoise := rec{X: x, Y: y, DX: int(d.X), DY: int(d.Y)}
		hare, ok := move(rec{X: x, Y: y, DX: int(d.X), DY: int(d.Y)})
		if !ok {
			return false
		}

		power, lam := 1, 1
		for tortoise != hare {
			if power == lam {
				tortoise = hare
				power *= 2
				lam = 0
			}
			hare, ok = move(hare)
			if !ok {
				return false
			}
			lam++
		}

		return true
	}

	curX, curY, curD := sx, sy, util.N
	for {
		var nx, ny int
		nd := curD

		alreadyVisited := g.MustGet(curX, curY) == 'X'
		next, ok := g.Get(curX+int(curD.X), curY+int(curD.Y))

		if !alreadyVisited {
			g.Set(curX, curY, 'X')
			silver++
		}

		if !ok {
			break
		}

		if next != 'X' && next != '#' && loopsIfBlocked(curX, curY, curD) {
			gold++
		}

		if next == '#' {
			nd = util.Point{
				X: -curD.Y,
				Y: curD.X,
			}
			nx, ny = curX, curY
		} else {
			nx, ny = curX+int(curD.X), curY+int(curD.Y)
		}

		curX, curY, curD = nx, ny, nd
	}

	fmt.Fprintln(out, "silver:", silver)
	fmt.Fprintln(out, "gold:", gold)
}

func Problem06ComplexG(in io.Reader, out io.Writer) {
	var silver, gold int64

	var sx, sy int
	g := util.GridFromReader(in)
	for y := range g.Height {
		for x := range g.Width {
			if g.MustGet(x, y) == '^' {
				sx = x
				sy = y
			}
		}
	}

	type pos struct {
		X int
		Y int
	}

	type rec struct {
		X  int
		Y  int
		DX int
		DY int
	}

	jumpTable := make(map[rec]pos)
	var mu sync.Mutex

	loopsIfBlocked := func(x, y int, d util.Point) bool {
		mu.Lock()
		defer mu.Unlock()
		ox, oy := x+int(d.X), y+int(d.Y)
		move := func(r rec) (rec, bool) {
			xx, yy := r.X, r.Y
			var tp pos

			if p, ok := jumpTable[r]; ok {
				tp = p
			} else {
				for {
					xxx, yyy := xx+r.DX, yy+r.DY
					p, ok := g.Get(xxx, yyy)
					if !ok {
						if xxx < 0 {
							tp = pos{
								X: math.MinInt,
								Y: yyy,
							}
						} else if xxx >= g.Width {
							tp = pos{
								X: math.MaxInt,
								Y: yyy,
							}
						} else if yyy < 0 {
							tp = pos{
								X: xxx,
								Y: math.MinInt,
							}
						} else {
							tp = pos{
								X: xxx,
								Y: math.MaxInt,
							}
						}
						break
					}
					if p == '#' {
						tp = pos{
							X: xx,
							Y: yy,
						}
						break
					}
					xx = xxx
					yy = yyy
				}
				jumpTable[r] = tp
			}

			xx, yy = tp.X, tp.Y
			if r.DX == 0 && ox == r.X {
				if r.DY == -1 && (yy <= oy && oy < r.Y) {
					// up
					yy = oy + 1
				} else if r.DY == 1 && (yy >= oy && oy > r.Y) {
					// down
					yy = oy - 1
				}
			}
			if r.DY == 0 && oy == r.Y {
				if r.DX == -1 && (xx <= ox && ox < r.X) {
					// left
					xx = ox + 1
				} else if r.DX == 1 && (xx >= ox && ox > r.X) {
					// right
					xx = ox - 1
				}
			}

			if !g.InBounds(xx, yy) {
				return rec{}, false
			}

			return rec{
				X:  xx,
				Y:  yy,
				DX: -r.DY,
				DY: r.DX,
			}, true
		}

		// Exponentially search for a cycle with Brent's algorithm
		// i.e. find i and j where f(2^i) = f(j)
		hare := rec{X: x, Y: y, DX: int(d.X), DY: int(d.Y)}
		tortoises := make([]rec, 31)
		tortoises[0] = rec{X: x, Y: y, DX: int(d.X), DY: int(d.Y)}
		i := 1
		var ok bool
		for {
			hare, ok = move(hare)
			if !ok {
				return false
			}
			k := 0
			pk := 1
			for k < len(tortoises) && pk <= i {
				if hare == tortoises[k] {
					return true
				}
				pk *= 2
				k += 1
			}
			tz := 0
			for k&1 == 0 {
				k = k >> 1
				tz++
			}
			tortoises[tz] = hare
			i++
		}
	}

	var wg sync.WaitGroup
	jobs := make(chan rec, 256)
	worker := func(jobs <-chan rec) {
		defer wg.Done()
		var g int64
		for j := range jobs {
			if loopsIfBlocked(j.X, j.Y, util.Point{X: int64(j.DX), Y: int64(j.DY)}) {
				g++
			}
		}
		atomic.AddInt64(&gold, g)
	}

	for range 8 {
		wg.Add(1)
		go worker(jobs)
	}

	curX, curY, curD := sx, sy, util.N
	for {
		var nx, ny int
		nd := curD

		alreadyVisited := g.MustGet(curX, curY) == 'X'
		next, ok := g.Get(curX+int(curD.X), curY+int(curD.Y))

		if !alreadyVisited {
			g.Set(curX, curY, 'X')
			silver++
		}

		if !ok {
			break
		}

		if next != 'X' && next != '#' {
			jobs <- rec{
				X:  curX,
				Y:  curY,
				DX: int(curD.X),
				DY: int(curD.Y),
			}
		}

		if next == '#' {
			nd = util.Point{
				X: -curD.Y,
				Y: curD.X,
			}
			nx, ny = curX, curY
		} else {
			nx, ny = curX+int(curD.X), curY+int(curD.Y)
		}

		curX, curY, curD = nx, ny, nd
	}

	close(jobs)
	wg.Wait()

	fmt.Fprintln(out, "silver:", silver)
	fmt.Fprintln(out, "gold:", gold)
}

func Problem06(in io.Reader, out io.Writer) {
	Problem06ComplexG(in, out)
}
