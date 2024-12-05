package main

import (
	"fmt"
	"io"
	"slices"
	"strings"
	"sync"

	"github.com/Fekinox/aoc-2024/util"
)

func Problem01(in io.Reader, out io.Writer) {
	input := util.ReadLines(in)

	listL := make([]int64, len(input))
	listR := make([]int64, len(input))
	occurrencesInRight := make(map[int64]int64)

	for i, l := range input {
		tokens := strings.Fields(l)
		listL[i] = util.MustParseInt(tokens[0])
		listR[i] = util.MustParseInt(tokens[1])
		if x, ok := occurrencesInRight[listR[i]]; ok {
			occurrencesInRight[listR[i]] = x + 1
		} else {
			occurrencesInRight[listR[i]] = 1
		}
	}

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		slices.Sort(listL)
		wg.Done()
	}()
	go func() {
		slices.Sort(listR)
		wg.Done()
	}()
	wg.Wait()

	var totalDist int64
	var similarity int64
	for i := range input {
		dist := util.Abs(listL[i] - listR[i])
		totalDist += dist

		if occs, ok := occurrencesInRight[listL[i]]; ok {
			similarity += listL[i] * occs
		}
	}

	fmt.Fprintln(out, "silver:", totalDist)
	fmt.Fprintln(out, "gold:", similarity)
}
