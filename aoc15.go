package main

import (
	"fmt"
	"io"
	"slices"
	"strings"

	"github.com/Fekinox/aoc-2024/util"
)

func Problem15Washed(in io.Reader, out io.Writer) {
}

func Problem15Unwashed(in io.Reader, out io.Writer) {
	var silver, gold int64

	groups := util.ReadNewlineSeparatedGroups(in)
	grid := util.GridFromStrings(groups[0]...)
	wideGrid := util.MakeGrid(grid.Width*2, grid.Height, '.')

	moves := strings.Join(groups[1], "")

	var robX, robY int
	var wideRobX, wideRobY int

	for y := range grid.Height {
		for x := range grid.Width {
			switch grid.MustGet(x, y) {
			case '@':
				robX, robY = x, y
				wideRobX, wideRobY = x*2, y
				wideGrid.Set(x*2, y, '@')
				break
			case '#':
				wideGrid.Set(x*2, y, '#')
				wideGrid.Set(x*2+1, y, '#')
				break
			case 'O':
				wideGrid.Set(x*2, y, '[')
				wideGrid.Set(x*2+1, y, ']')
				break
			}
		}
	}

	type pos struct {
		X int
		Y int
	}

normgrid:
	for _, m := range moves {
		// for y := range grid.Height {
		// 	for x := range grid.Width {
		// 		fmt.Fprintf(out, "%c", grid.MustGet(x, y))
		// 	}
		// 	fmt.Fprintln(out, "")
		// }
		var dx, dy int
		switch m {
		case '<':
			dx = -1
		case '>':
			dx = +1
		case '^':
			dy = -1
		case 'v':
			dy = 1
		}

		var nextMoves, stack []pos

		stack = append(stack, pos{
			X: robX,
			Y: robY,
		})

		for len(stack) > 0 {
			cur := stack[len(stack)-1]
			stack = stack[:len(stack)-1]

			nextCur := pos{
				X: cur.X + dx,
				Y: cur.Y + dy,
			}
			switch grid.MustGet(nextCur.X, nextCur.Y) {
			case '#':
				continue normgrid
			case 'O':
				stack = append(stack, nextCur)
			}
			nextMoves = append(nextMoves, cur)
		}

		for mv := range nextMoves {
			thePos := nextMoves[len(nextMoves)-1-mv]
			grid.Set(thePos.X+dx, thePos.Y+dy, grid.MustGet(thePos.X, thePos.Y))
			grid.Set(thePos.X, thePos.Y, '.')
		}

		robX += dx
		robY += dy
	}

widegrid:
	for _, m := range moves {
		// for y := range wideGrid.Height {
		// 	for x := range wideGrid.Width {
		// 		fmt.Fprintf(out, "%c", wideGrid.MustGet(x, y))
		// 	}
		// 	fmt.Fprintln(out, "")
		// }
		var dx, dy int
		switch m {
		case '<':
			dx = -1
		case '>':
			dx = +1
		case '^':
			dy = -1
		case 'v':
			dy = 1
		}

		var nextMoves, queue []pos

		queue = append(queue, pos{
			X: wideRobX,
			Y: wideRobY,
		})

		for len(queue) > 0 {
			cur := queue[0]
			queue = queue[1:]

			if !slices.Contains(nextMoves, cur) {
				nextMoves = append(nextMoves, cur)
			}
			nextCur := pos{
				X: cur.X + dx,
				Y: cur.Y + dy,
			}
			switch wideGrid.MustGet(nextCur.X, nextCur.Y) {
			case '#':
				continue widegrid
			case '[':
				if dy == 0 {
					queue = append(queue, nextCur)
				} else {
					queue = append(queue, nextCur, pos{
						X: nextCur.X + 1,
						Y: nextCur.Y,
					})
				}
				break
			case ']':
				if dy == 0 {
					queue = append(queue, nextCur)
				} else {
					queue = append(queue, nextCur, pos{
						X: nextCur.X - 1,
						Y: nextCur.Y,
					})
				}
				break
			}
		}

		fmt.Fprintln(out, nextMoves)

		for mv := range nextMoves {
			thePos := nextMoves[len(nextMoves)-1-mv]
			wideGrid.Set(thePos.X+dx, thePos.Y+dy, wideGrid.MustGet(thePos.X, thePos.Y))
			wideGrid.Set(thePos.X, thePos.Y, '.')
		}

		wideRobX += dx
		wideRobY += dy
	}

	for y := range grid.Height {
		for x := range grid.Width {
			if grid.MustGet(x, y) == 'O' {
				silver += int64(x) + 100*int64(y)
			}
		}
	}

	for y := range wideGrid.Height {
		for x := range wideGrid.Width {
			if wideGrid.MustGet(x, y) == '[' {
				gold += int64(x) + 100*int64(y)
			}
		}
	}

	fmt.Fprintln(out, "silver:", silver)
	fmt.Fprintln(out, "gold:", gold)
}

func Problem15(in io.Reader, out io.Writer) {
	Problem15Washed(in, out)
}
