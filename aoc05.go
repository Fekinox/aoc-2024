package main

import (
	"fmt"
	"io"
	"slices"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/Fekinox/aoc-2024/util"
)

const JOB_SIZE = 9999999

type Pair struct {
	First int64
	Second int64
}

func Problem05Par(in io.Reader, out io.Writer) {
	var silver, gold int64
	var wg sync.WaitGroup

	groups := util.ReadNewlineSeparatedGroups(in)
	comesBefore := make(map[Pair]struct{})

	for _, l := range groups[0] {
		p := strings.Split(l, "|")
		left := util.MustParseInt(p[0])
		right := util.MustParseInt(p[1])

		comesBefore[Pair{First: left, Second: right}] = struct{}{}
	}

	sub := func(i, j int) (silver, gold int64) {
		next:
		for _, l := range groups[1][i:min(j, len(groups[1]))] {
			nums := strings.Split(l, ",")
			convNums := make([]int64, len(nums))
			for j, n := range nums {
				convNums[j] = util.MustParseInt(n)
			}

			for j := range len(convNums)-1 {
				if _, ok := comesBefore[Pair{First: convNums[j], Second: convNums[j+1]}]; !ok {
					slices.SortFunc(convNums, func(x int64, y int64) int {
						if x == y {
							return 0
						}
						if _, ok := comesBefore[Pair{First: x, Second: y}]; ok {
							return -1
						} else {
							return 1
						}
					})

					gold += convNums[len(convNums)/2]
					continue next
				}
			}

			silver += convNums[len(convNums)/2]
		}

		return silver, gold
	}

	for i := 0; i < len(groups[1]); i += JOB_SIZE {
		wg.Add(1)
		go func() {
			defer wg.Done()
			s, g := sub(i, i+JOB_SIZE)

			atomic.AddInt64(&silver, s)
			atomic.AddInt64(&gold, g)
		}()
	}

	wg.Wait()


	fmt.Fprintln(out, "silver:", silver)
	fmt.Fprintln(out, "gold:", gold)
}

func Problem05Seq(in io.Reader, out io.Writer) {
	var silver, gold int64

	groups := util.ReadNewlineSeparatedGroups(in)
	// comesBefore := make(map[Pair]struct{})
	comesBefore := util.MakeGrid(100, 100, false)

	for _, l := range groups[0] {
		p := strings.Split(l, "|")
		comesBefore.Set(int(util.MustParseInt(p[0])), int(util.MustParseInt(p[1])), true)
	}

	next:
	for _, l := range groups[1] {
		nums := strings.Split(l, ",")
		convNums := make([]int, len(nums))
		for j, n := range nums {
			convNums[j] = int(util.MustParseInt(n))
		}

		for j := range len(convNums)-1 {
			if ok := comesBefore.MustGet(int(convNums[j]), int(convNums[j+1])); !ok {
				slices.SortFunc(convNums, func(x int, y int) int {
					if x == y {
						return 0
					}
					if ok := comesBefore.MustGet(int(x), int(y)); ok {
						return -1
					} else {
						return 1
					}
				})

				gold += int64(convNums[len(convNums)/2])
				continue next
			}
		}

		silver += int64(convNums[len(convNums)/2])
	}

	fmt.Fprintln(out, "silver:", silver)
	fmt.Fprintln(out, "gold:", gold)
}

func Problem05(in io.Reader, out io.Writer) {
	Problem05Seq(in, out)
}
