package main

import (
	"fmt"
	"io"

	"github.com/Fekinox/aoc-2024/util"
)

func Problem10Array(in io.Reader, out io.Writer) {
	var silver, gold int64

	g := util.GridFromReader(in)

	for y := range g.Height {
	outer:
		for x := range g.Width {
			if g.MustGet(x, y) != '0' {
				continue
			}

			frontier := make(map[util.Point]int64)
			frontier[util.Point{X: int64(x), Y: int64(y)}] = 1

			for range 9 {
				if len(frontier) == 0 {
					continue outer
				}

				next := make(map[util.Point]int64)
				for k, curScore := range frontier {
					curValue := g.MustGet(int(k.X), int(k.Y))
					for _, d := range util.CardinalDirections {
						xx, yy := k.X+d.X, k.Y+d.Y
						if c, ok := g.Get(int(xx), int(yy)); ok && c-curValue == 1 {
							next[util.Point{X: xx, Y: yy}] += curScore
						}
					}
				}
				frontier = next
			}

			silver += int64(len(frontier))
			for _, v := range frontier {
				gold += v
			}
		}
	}

	fmt.Fprintln(out, "silver:", silver)
	fmt.Fprintln(out, "gold:", gold)
}

func Problem10(in io.Reader, out io.Writer) {
	Problem10Array(in, out)
}
