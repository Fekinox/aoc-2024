package main

import (
	"fmt"
	"io"

	"github.com/Fekinox/aoc-2024/util"
)

func Problem20Washed(in io.Reader, out io.Writer) {
	var silver, gold int64

	type state struct {
		X int
		Y int
	}

	g := util.GridFromReader(in)

	var sx, sy, ex, ey int
	for y := range g.Height {
		for x := range g.Width {
			if g.MustGet(x, y) == 'S' {
				sx, sy = x, y
			}
			if g.MustGet(x, y) == 'E' {
				ex, ey = x, y
			}
		}
	}

	var path []state
	cur := state{
		X: sx,
		Y: sy,
	}
	for {
		path = append(path, cur)

		if cur.X == ex && cur.Y == ey {
			break
		}

		for _, d := range util.CardinalDirections {
			xx, yy := cur.X+int(d.X), cur.Y+int(d.Y)
			if v, ok := g.Get(xx, yy); !ok || v == '#' {
				continue
			}
			if len(path) > 1 && path[len(path)-2].X == xx && path[len(path)-2].Y == yy {
				continue
			}
			cur = state{
				X: xx,
				Y: yy,
			}
			break
		}
	}

	totalDist := len(path) - 1

	for d1 := 0; d1 < len(path); d1++ {
		for d2 := d1 + 100; d2 < len(path); d2++ {
			p1, p2 := path[d1], path[d2]
			d := util.Abs(p1.X-p2.X) + util.Abs(p1.Y-p2.Y)
			cheatDist := d1 + d + (totalDist - d2)
			if (totalDist - cheatDist) >= 100 {
				if d <= 2 {
					silver++
				}
				if d <= 20 {
					gold++
				}
			}
		}
	}

	fmt.Fprintln(out, "silver:", silver)
	fmt.Fprintln(out, "gold:", gold)
}

func Problem20Unwashed(in io.Reader, out io.Writer) {
	var silver, gold int64

	type pos struct {
		X int
		Y int
	}

	g := util.GridFromReader(in)

	var sx, sy, ex, ey int

	for y := range g.Height {
		for x := range g.Width {
			if g.MustGet(x, y) == 'S' {
				sx, sy = x, y
			}
			if g.MustGet(x, y) == 'E' {
				ex, ey = x, y
			}
		}
	}

	distances := util.MakeGrid(g.Width, g.Height, -1)
	var path []pos
	cx, cy := sx, sy
	dist := 0
	for cx != ex || cy != ey {
		fmt.Println(cx, cy)
		distances.Set(cx, cy, dist)
		path = append(path, pos{
			X: cx,
			Y: cy,
		})
	dir:
		for _, d := range util.CardinalDirections {
			xx, yy := cx+int(d.X), cy+int(d.Y)
			if v, ok := g.Get(xx, yy); ok && v != '#' && distances.MustGet(xx, yy) == -1 {
				cx, cy = xx, yy
				dist++
				break dir
			}
		}
	}

	path = append(path, pos{
		X: ex,
		Y: ey,
	})
	distances.Set(ex, ey, dist)

	savings := make(map[int]int)

	for _, ps := range path {
		regDist := distances.MustGet(ps.X, ps.Y)
		for _, d := range util.CardinalDirections {
			xx, yy := ps.X+2*int(d.X), ps.Y+2*int(d.Y)
			dd, ok := distances.Get(xx, yy)
			if !ok || dd == -1 {
				continue
			}
			cheatDist := 2 + (dist - dd)
			if dist-(regDist+cheatDist) >= 100 {
				silver++
			}

			fmt.Println(dist, regDist+cheatDist)
		}
	}

	for _, p1 := range path {
		for _, p2 := range path {
			d := util.Abs(p1.X-p2.X) + util.Abs(p1.Y-p2.Y)
			if d > 20 {
				continue
			}
			d1, d2 := distances.MustGet(p1.X, p1.Y), distances.MustGet(p2.X, p2.Y)
			cheatDist := d1 + d + (dist - d2)
			if dist-cheatDist >= 100 {
				gold++
			}

			if dist-cheatDist >= 0 {
				savings[dist-cheatDist]++
			}
		}
	}

	fmt.Println(savings)

	fmt.Fprintln(out, "silver:", silver)
	fmt.Fprintln(out, "gold:", gold)
}

func Problem20(in io.Reader, out io.Writer) {
	Problem20Washed(in, out)
}
