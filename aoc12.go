package main

import (
	"fmt"
	"io"

	"github.com/Fekinox/aoc-2024/util"
)

func Problem12SpanFloodFill(in io.Reader, out io.Writer) {
}

func Problem12StackFloodFill(in io.Reader, out io.Writer) {
	var silver, gold int64

	g := util.GridFromReader(in)
	visited := util.NewBoolVector(g.Width * g.Height)

	for y := range g.Height {
		for x := range g.Width {
			if visited.Get(y*g.Width + x) {
				continue
			}

			frontier := []util.Point{{
				X: int64(x),
				Y: int64(y),
			}}

			curVal := g.MustGet(x, y)

			var area, perimeter, numSides int64

			for len(frontier) > 0 {
				cur := frontier[0]
				frontier = frontier[1:]

				if visited.Get(int(cur.Y)*g.Width + int(cur.X)) {
					continue
				}

				visited.Set(int(cur.Y)*g.Width + int(cur.X))
				area++

				for _, d := range util.CardinalDirections {
					xx, yy := cur.X+d.X, cur.Y+d.Y
					if val, ok := g.Get(int(xx), int(yy)); !ok || val != curVal {
						perimeter++

						newD := util.Point{
							X: d.Y,
							Y: -d.X,
						}
						sideVal, sideOk := g.Get(int(cur.X+newD.X), int(cur.Y+newD.Y))
						if !sideOk || sideVal != curVal {
							numSides++
							continue
						}
						stepVal, stepOk := g.Get(int(cur.X+newD.X+d.X), int(cur.Y+newD.Y+d.Y))
						if stepOk && stepVal == curVal {
							numSides++
							continue
						}
					} else {
						frontier = append(frontier, util.Point{X: xx, Y: yy})
					}
				}
			}
			silver += area * perimeter
			gold += area * numSides
		}
	}

	fmt.Fprintln(out, "silver:", silver)
	fmt.Fprintln(out, "gold:", gold)
}

func Problem12(in io.Reader, out io.Writer) {
	Problem12SpanFloodFill(in, out)
}
