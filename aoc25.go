package main

import (
	"fmt"
	"io"

	"github.com/Fekinox/aoc-2024/util"
)

func Problem25(in io.Reader, out io.Writer) {
	var silver int64

	groups := util.ReadNewlineSeparatedGroups(in)
	var locks, keys [][]int
	for _, g := range groups {
		gr := util.GridFromStrings(g...)
		var isLock bool
		for i := range gr.Width {
			top, bot := gr.MustGet(i, 0), gr.MustGet(i, gr.Height-1)
			if top == '.' {
				isLock = false
				break
			}
			if bot == '.' {
				isLock = true
				break
			}
		}

		if isLock {
			var heights []int
			for i := range gr.Width {
				var h int
				for j := range gr.Height {
					if v, ok := gr.Get(i, j+1); !ok || v == '.' {
						break
					}
					h++
				}
				heights = append(heights, h)
			}
			locks = append(locks, heights)
		} else {
			var heights []int
			for i := range gr.Width {
				var h int
				for j := range gr.Height {
					if v, ok := gr.Get(i, gr.Height-2-j); !ok || v == '.' {
						break
					}
					h++
				}
				heights = append(heights, h)
			}
			keys = append(keys, heights)
		}
	}

	for _, l := range locks {
		outer:
		for _, k := range keys {
			for i := range len(l) {
				if (5-l[i]) < k[i] {
					continue outer		
				}
			}
			silver++
		}
	}

	var gold string

	gold = "Merry Christmas!"
	fmt.Fprintln(out, "silver:", silver)
	fmt.Fprintln(out, "gold:", gold)
}
