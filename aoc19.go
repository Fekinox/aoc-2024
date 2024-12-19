package main

import (
	"fmt"
	"io"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/Fekinox/aoc-2024/util"
)

func Problem19Bytes(in io.Reader, out io.Writer) {
	var silver, gold int64

	var wg sync.WaitGroup
	const GROUP_SIZE int = 16

	groups := util.ReadNewlineSeparatedGroups(in)

	towels := make([][][]byte, 5)
	for _, t := range strings.Split(groups[0][0], ", ") {
		var index int
		switch t[0] {
		case 'u':
			index = 1
		case 'b':
			index = 2
		case 'r':
			index = 3
		case 'g':
			index = 4
		}
		towels[index] = append(towels[index], []byte(t))
	}

	var designs [][]byte
	for _, g := range groups[1] {
		designs = append(designs, []byte(g))
	}

	worker := func(data [][]byte) {
		defer wg.Done()
		var s, g int64

		for _, d := range data {
			possible := make([]int64, len(d)+1)
			possible[len(d)] = 1
			for j := range len(d) {
				jj := len(d) - j - 1
				var index int
				switch d[jj] {
				case 'u':
					index = 1
				case 'b':
					index = 2
				case 'r':
					index = 3
				case 'g':
					index = 4
				}
			outer:
				for _, t := range towels[index] {
					if jj+len(t) > len(d) {
						continue
					}
					for i, k := range t {
						if k != d[jj+i] {
							continue outer
						}
					}
					possible[jj] += possible[jj+len(t)]
				}
			}

			if possible[0] > 0 {
				s++
			}
			g += possible[0]
		}
		atomic.AddInt64(&silver, s)
		atomic.AddInt64(&gold, g)
	}

	for i := 0; i <= len(groups[1]); i += GROUP_SIZE {
		wg.Add(1)
		go worker(designs[i:min(len(groups[1]), i+GROUP_SIZE)])
	}

	wg.Wait()

	fmt.Fprintln(out, "silver:", silver)
	fmt.Fprintln(out, "gold:", gold)
}

func Problem19Strings(in io.Reader, out io.Writer) {
	var silver, gold int64

	var wg sync.WaitGroup
	const GROUP_SIZE int = 16

	groups := util.ReadNewlineSeparatedGroups(in)

	towels := make([][]string, 5)
	for _, t := range strings.Split(groups[0][0], ", ") {
		var index int
		switch t[0] {
		case 'u':
			index = 1
		case 'b':
			index = 2
		case 'r':
			index = 3
		case 'g':
			index = 4
		}
		towels[index] = append(towels[index], t)
	}

	worker := func(data []string) {
		defer wg.Done()
		var s, g int64

		for _, d := range data {
			possible := make([]int64, len(d)+1)
			possible[len(d)] = 1
			for j := range len(d) {
				jj := len(d) - j - 1
				var index int
				switch d[jj] {
				case 'u':
					index = 1
				case 'b':
					index = 2
				case 'r':
					index = 3
				case 'g':
					index = 4
				}
				for _, t := range towels[index] {
					if jj+len(t) <= len(d) && strings.HasPrefix(d[jj:], t) {
						possible[jj] = possible[jj] + possible[jj+len(t)]
					}
				}
			}

			if possible[0] > 0 {
				s++
			}
			g += possible[0]
		}
		atomic.AddInt64(&silver, s)
		atomic.AddInt64(&gold, g)
	}

	for i := 0; i <= len(groups[1]); i += GROUP_SIZE {
		wg.Add(1)
		go worker(groups[1][i:min(len(groups[1]), i+GROUP_SIZE)])
	}

	wg.Wait()

	fmt.Fprintln(out, "silver:", silver)
	fmt.Fprintln(out, "gold:", gold)
}

func Problem19(in io.Reader, out io.Writer) {
	Problem19Bytes(in, out)
}
