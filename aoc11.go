package main

import (
	"fmt"
	"io"
	"strings"

	"github.com/Fekinox/aoc-2024/util"
)

func Problem11NoString(in io.Reader, out io.Writer) {
	var silver, gold int64

	ln := util.ReadString(in)

	mapping := make(map[int64]int64)

	for _, s := range util.SplitNumbers(strings.TrimSpace(ln), " ") {
		mapping[s]++
	}

	for i := range 75 {
		if i == 25 {
			for _, c := range mapping {
				silver += c
			}
		}
		next := make(map[int64]int64)
	outer:
		for s, count := range mapping {
			if count == 0 {
				continue
			}
			if s == 0 {
				next[1] += count
			} else {
				ss := s
				pow := int64(1)
				for ss > 0 {
					if ss < 10 {
						next[s*2024] += count
						continue outer
					}
					ss /= 100
					pow *= 10
				}
				next[s/pow] += count
				next[s%pow] += count
			}
		}
		mapping = next
	}

	for _, c := range mapping {
		gold += c
	}

	fmt.Fprintln(out, "silver:", silver)
	fmt.Fprintln(out, "gold:", gold)
}

func Problem11String(in io.Reader, out io.Writer) {
	var silver, gold int64

	ln := util.ReadString(in)

	mapping := make(map[int64]int64)

	for _, s := range util.SplitNumbers(strings.TrimSpace(ln), " ") {
		mapping[s]++
	}

	for i := range 75 {
		if i == 25 {
			for _, c := range mapping {
				silver += c
			}
		}
		next := make(map[int64]int64)
		for s, count := range mapping {
			if count == 0 {
				continue
			}
			if s == 0 {
				next[1] += count
			} else if conv := fmt.Sprintf("%d", s); len(conv)%2 == 0 {
				next[util.MustParseInt(conv[len(conv)/2:])] += count
				next[util.MustParseInt(conv[:len(conv)/2])] += count
			} else {
				next[s*2024] += count
			}
		}
		mapping = next
	}

	for _, c := range mapping {
		gold += c
	}

	fmt.Fprintln(out, "silver:", silver)
	fmt.Fprintln(out, "gold:", gold)
}

func Problem11(in io.Reader, out io.Writer) {
	Problem11NoString(in, out)
}
