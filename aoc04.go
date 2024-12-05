package main

import (
	"fmt"
	"io"
	"sync"
	"sync/atomic"

	"github.com/Fekinox/aoc-2024/util"
)

var xmas = []byte("XMAS")
const CELL_WIDTH = 64

func Problem04(in io.Reader, out io.Writer) {
	var silver, gold int64
	var wg sync.WaitGroup

	input := util.GridFromReader(in)

	search := func(x, y int, dx, dy int) bool {
		xx := x
		yy := y
		dist := 0

		for range 4 {
			if c, ok := input.Get(xx, yy); !ok || c != xmas[dist] {
				return false
			}
			xx += dx
			yy += dy
			dist++
		}

		return true
	}

	mas := func(x, y int, dx, dy int) bool {
		pre, preOk := input.Get(x-dx, y-dy)
		post, postOk := input.Get(x+dx, y+dy)
		return preOk && postOk && ((pre == 'M' && post == 'S') || (pre == 'S' && post == 'M'))
	}

	sub := func(ox, oy int, ow, oh int) (silver, gold int64) {
		for y := oy; y < min(input.Height, oy + oh); y++ {
			for x := ox; x < min(input.Width, ox + ow); x++ {
				var c = input.MustGet(x, y)
				if c == 'X' {
					for _, d := range util.AllDirections {
						if search(x, y, int(d.X), int(d.Y)) {
							silver++
						}
					}
				}
				if c == 'A' && mas(x, y, 1, 1) && mas(x, y, 1, -1) {
					gold++
				}
			}
		}

		return silver, gold
	}

	for sy := 0; sy <= input.Height; sy += CELL_WIDTH {
		for sx := 0; sx <= input.Width; sx += CELL_WIDTH {
			wg.Add(1)
			go func() {
				defer wg.Done()
				s, g := sub(sx, sy, CELL_WIDTH, CELL_WIDTH)

				atomic.AddInt64(&silver, s)
				atomic.AddInt64(&gold, g)
			}()
		}
	}
	wg.Wait()

	fmt.Fprintln(out, "silver:", silver)
	fmt.Fprintln(out, "gold:", gold)
}
