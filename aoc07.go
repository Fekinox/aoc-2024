package main

import (
	"fmt"
	"io"
	"slices"
	"strings"

	"github.com/Fekinox/aoc-2024/util"
)

func concat(x, y int64) int64 {
	exp := int64(1)
	yy := y
	for yy > 0 {
		yy /= 10
		exp *= 10
	}
	return exp*x + y
}

func Problem07BFS(in io.Reader, out io.Writer) {
	var silver, gold int64
	input := util.ReadLines(in)

	for _, l := range input {
		tokens := strings.Split(l, ": ")
		testvals := util.SplitNumbers(tokens[1], " ")
		target := util.MustParseInt(tokens[0])

		i := 1
		cur := []int64{testvals[0]}
		goldCur := []int64{testvals[0]}
		for i < len(testvals) {
			if len(cur) == 0 && len(goldCur) == 0 {
				break
			}

			nextVal := testvals[i]
			next := make([]int64, 0, len(cur)*2)
			goldNext := make([]int64, 0, len(goldCur)*3)

			slices.Sort(cur)
			slices.Sort(goldCur)

			for _, c := range cur {
				if c >= target {
					break
				}
				if nextVal+c <= target {
					next = append(next, nextVal+c)
				}
				if nextVal*c <= target {
					next = append(next, nextVal*c)
				}
			}

			for _, c := range goldCur {
				if c >= target {
					break
				}
				if nextVal+c <= target {
					goldNext = append(goldNext, nextVal+c)
				}
				if nextVal*c <= target {
					goldNext = append(goldNext, nextVal*c)
				}
				cc := concat(c, nextVal)
				if cc <= target {
					goldNext = append(goldNext, cc)
				}
			}

			cur = next
			goldCur = goldNext
			i++
		}
		if i == len(testvals) {
			for _, c := range cur {
				if c == target {
					silver += target
					break
				}
			}
			for _, c := range goldCur {
				if c == target {
					gold += target
					break
				}
			}
		}
	}

	fmt.Fprintln(out, "silver:", silver)
	fmt.Fprintln(out, "gold:", gold)
}

func sevenDfs(target int64, vals []int64, current int64, usedConcat bool) (ok bool, cc bool) {
	if current == target && len(vals) == 0 {
		return true, usedConcat
	} else if current > target || len(vals) == 0 {
		return false, false
	}

	ok, cc = sevenDfs(target, vals[1:], current*vals[0], usedConcat)
	if ok {
		return
	}
	ok, cc = sevenDfs(target, vals[1:], current+vals[0], usedConcat)
	if ok {
		return
	}
	return sevenDfs(target, vals[1:], concat(current, vals[0]), true)
}

func Problem07DFS(in io.Reader, out io.Writer) {
	var silver, gold int64
	input := util.ReadLines(in)

	for _, l := range input {
		tokens := strings.Split(l, ": ")
		testvals := util.SplitNumbers(tokens[1], " ")
		target := util.MustParseInt(tokens[0])

		ok, concat := sevenDfs(target, testvals[1:], testvals[0], false)

		if ok {
			if !concat {
				silver += target
			}
			gold += target
		}
	}

	fmt.Fprintln(out, "silver:", silver)
	fmt.Fprintln(out, "gold:", gold)
}

func Problem07Stack(in io.Reader, out io.Writer) {
	var silver, gold int64
	input := util.ReadLines(in)

	type dfsIn struct {
		Current int64
		Index   int
		Concat  bool
	}

	for _, l := range input {
		tokens := strings.Split(l, ": ")
		testvals := util.SplitNumbers(tokens[1], " ")
		target := util.MustParseInt(tokens[0])

		stack := []dfsIn{{
			Current: testvals[0],
			Index:   1,
		}}

		var ok, cc bool

		for len(stack) > 0 {
			top := stack[len(stack)-1]
			stack = stack[:len(stack)-1]

			if top.Current == target && top.Index >= len(testvals) {
				ok = true
				cc = top.Concat
				break
			} else if top.Current > target || top.Index >= len(testvals) {
				continue
			}

			stack = append(stack,
				dfsIn{
					Current: concat(top.Current, testvals[top.Index]),
					Index:   top.Index + 1,
					Concat:  true,
				},
				dfsIn{
					Current: top.Current + testvals[top.Index],
					Index:   top.Index + 1,
					Concat:  top.Concat,
				},
				dfsIn{
					Current: top.Current * testvals[top.Index],
					Index:   top.Index + 1,
					Concat:  top.Concat,
				},
			)
		}

		if ok {
			if !cc {
				silver += target
			}
			gold += target
		}
	}

	fmt.Fprintln(out, "silver:", silver)
	fmt.Fprintln(out, "gold:", gold)
}

func Problem07(in io.Reader, out io.Writer) {
	Problem07Stack(in, out)
}
