package main

import (
	"fmt"
	"io"

	"github.com/Fekinox/aoc-2024/util"
)

func Problem08Hashmap(in io.Reader, out io.Writer) {
	var silver, gold int64
	g := util.GridFromReader(in)

	freqs := make(map[byte][]util.Point)

	anti := util.MakeGrid(g.Width, g.Height, false)
	goldAnti := util.MakeGrid(g.Width, g.Height, false)

	for y := range g.Height {
		for x := range g.Width {
			c := g.MustGet(x, y)
			if ('0' <= c && c <= '9') || ('a' <= c && c <= 'z') || ('A' <= c && c <= 'Z') {
				freqs[c] = append(freqs[c], util.Point{X: int64(x), Y: int64(y)})
			}
		}
	}

	for _, f := range freqs {
		for i := range len(f) {
			for j := i + 1; j < len(f); j++ {
				deltaX, deltaY := int(f[i].X-f[j].X), int(f[i].Y-f[j].Y)
				ax1, ay1 := int(f[i].X)+deltaX, int(f[i].Y)+deltaY
				ax2, ay2 := int(f[j].X)-deltaX, int(f[j].Y)-deltaY
				if anti.InBounds(ax1, ay1) {
					anti.Set(ax1, ay1, true)
					goldAnti.Set(ax1, ay1, true)
				}
				if anti.InBounds(ax2, ay2) {
					anti.Set(ax2, ay2, true)
					goldAnti.Set(ax2, ay2, true)
				}

				// gold
				var xx, yy int
				xx, yy = int(f[i].X), int(f[i].Y)
				for goldAnti.InBounds(xx, yy) {
					goldAnti.Set(xx, yy, true)
					xx, yy = xx+deltaX, yy+deltaY
				}

				xx, yy = int(f[i].X), int(f[i].Y)
				for goldAnti.InBounds(xx, yy) {
					goldAnti.Set(xx, yy, true)
					xx, yy = xx-deltaX, yy-deltaY
				}
			}
		}
	}

	for y := range g.Height {
		for x := range g.Width {
			if anti.MustGet(x, y) {
				silver++
			}
			if goldAnti.MustGet(x, y) {
				gold++
			}
		}
	}

	fmt.Fprintln(out, "silver:", silver)
	fmt.Fprintln(out, "gold:", gold)
}

func Problem08Array(in io.Reader, out io.Writer) {
	var silver, gold int64
	g := util.GridFromReader(in)

	freqs := make([][]util.Point, 10+26+26)

	anti := util.MakeGrid(g.Width, g.Height, false)
	goldAnti := util.MakeGrid(g.Width, g.Height, false)

	for y := range g.Height {
		for x := range g.Width {
			c := g.MustGet(x, y)
			var idx int
			if '0' <= c && c <= '9' {
				idx = int(c - '0')
			} else if 'a' <= c && c <= 'z' {
				idx = 10 + int(c-'a')
			} else if 'A' <= c && c <= 'Z' {
				idx = 10 + 26 + int(c-'A')
			} else {
				continue
			}
			freqs[idx] = append(freqs[idx], util.Point{X: int64(x), Y: int64(y)})
		}
	}

	for _, f := range freqs {
		for i := range len(f) {
			for j := i + 1; j < len(f); j++ {
				deltaX, deltaY := int(f[i].X-f[j].X), int(f[i].Y-f[j].Y)

				// gold
				var m int
				xx, yy := int(f[i].X), int(f[i].Y)
				for goldAnti.InBounds(xx+m*deltaX, yy+m*deltaY) {
					goldAnti.Set(xx, yy, true)
					if m == 1 {
						anti.Set(xx, yy, true)
					}
					m++
				}

				m = -1
				for goldAnti.InBounds(xx+m*deltaX, yy+m*deltaY) {
					goldAnti.Set(xx, yy, true)
					if m == -2 {
						anti.Set(xx, yy, true)
					}
					m--
				}
			}
		}
	}

	for y := range g.Height {
		for x := range g.Width {
			if anti.MustGet(x, y) {
				silver++
			}
			if goldAnti.MustGet(x, y) {
				gold++
			}
		}
	}

	fmt.Fprintln(out, "silver:", silver)
	fmt.Fprintln(out, "gold:", gold)
}

func Problem08(in io.Reader, out io.Writer) {
	Problem08Hashmap(in, out)
}
