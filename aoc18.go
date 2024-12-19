package main

import (
	"fmt"
	"io"

	"github.com/Fekinox/aoc-2024/util"
)

func Problem18Washed(in io.Reader, out io.Writer) {
	var silver int64
	var gold string

	g := util.MakeGrid(71, 71, -1)
	lines := util.ReadLines(in)
	for i, l := range lines {
		tokens := util.SplitNumbers(l, ",")
		g.Set(int(tokens[0]), int(tokens[1]), i)
	}

	type pos struct {
		X int
		Y int
	}

	search := func(lim int) int64 {
		vis := util.MakeGrid(71, 71, false)
		frontier := []pos{{
			X: 0,
			Y: 0,
		}}

		var d int64
		for len(frontier) > 0 {
			var nextFrontier []pos
			for _, f := range frontier {
				if v, ok := vis.Get(f.X, f.Y); !ok || v {
					continue
				}
				if wall, ok := g.Get(f.X, f.Y); !ok || (wall != -1 && wall <= lim) {
					continue
				}
				vis.Set(f.X, f.Y, true)
				if f.X == 70 && f.Y == 70 {
					return d
				}
				for _, dir := range util.CardinalDirections {
					nextFrontier = append(nextFrontier, pos{
						X: f.X + int(dir.X),
						Y: f.Y + int(dir.Y),
					})
				}
			}
			frontier = nextFrontier
			d++
		}

		return -1
	}

	lo, hi := 0, len(lines)
	for lo < hi {
		m := lo + (hi-lo)/2
		if search(m) != -1 {
			lo = m + 1
		} else {
			hi = m
		}
	}
	gold = lines[lo]

	silver = search(1023)

	// render(1024)
	// render(lo)

	fmt.Fprintln(out, "silver:", silver)
	fmt.Fprintln(out, "gold:", gold)
}

func Problem18Unwashed(in io.Reader, out io.Writer) {
	var silver int64
	var gold string

	g := util.MakeGrid(71, 71, false)
	lines := util.ReadLines(in)
	for _, l := range lines[:1024] {
		tokens := util.SplitNumbers(l, ",")
		g.Set(int(tokens[0]), int(tokens[1]), true)
	}

	vis := util.MakeGrid(71, 71, false)

	type pos struct {
		X int
		Y int
	}

	frontier := []pos{{
		X: 0,
		Y: 0,
	}}

	var d int64
outer:
	for len(frontier) > 0 {
		var nextFrontier []pos
		for _, f := range frontier {
			if v, ok := vis.Get(f.X, f.Y); !ok || v {
				continue
			}
			if wall, ok := g.Get(f.X, f.Y); !ok || wall {
				continue
			}
			vis.Set(f.X, f.Y, true)
			if f.X == 70 && f.Y == 70 {
				break outer
			}
			for _, dir := range util.CardinalDirections {
				nextFrontier = append(nextFrontier, pos{
					X: f.X + int(dir.X),
					Y: f.Y + int(dir.Y),
				})
			}
		}
		frontier = nextFrontier
		d++
	}

	gg := util.MakeGrid(71, 71, false)
gold:
	for _, ln := range lines {
		tokens := util.SplitNumbers(ln, ",")
		gg.Set(int(tokens[0]), int(tokens[1]), true)

		vis := util.MakeGrid(71, 71, false)

		frontier := []pos{{
			X: 0,
			Y: 0,
		}}

		for len(frontier) > 0 {
			var nextFrontier []pos
			for _, f := range frontier {
				if v, ok := vis.Get(f.X, f.Y); !ok || v {
					continue
				}
				if wall, ok := gg.Get(f.X, f.Y); !ok || wall {
					continue
				}
				vis.Set(f.X, f.Y, true)
				if f.X == 70 && f.Y == 70 {
					continue gold
				}
				for _, dir := range util.CardinalDirections {
					nextFrontier = append(nextFrontier, pos{
						X: f.X + int(dir.X),
						Y: f.Y + int(dir.Y),
					})
				}
			}
			frontier = nextFrontier
		}

		gold = ln
		break
	}

	silver = d

	fmt.Fprintln(out, "silver:", silver)
	fmt.Fprintln(out, "gold:", gold)
}

func Problem18(in io.Reader, out io.Writer) {
	Problem18Washed(in, out)
}
