package main

import (
	"fmt"
	"io"
	"sync"
	"sync/atomic"

	"github.com/Fekinox/aoc-2024/util"
)

func Problem22(in io.Reader, out io.Writer) {
	var silver, gold int64

	type state struct {
		count  int64
		buyers *util.BoolVector
	}

	inp := util.ReadLines(in)

	prices := make([][]int, len(inp))

	// sales := make(map[int]*state)
	sales := make([]*state, 130321)
	for i := range sales {
		sales[i] = &state{
			buyers: util.NewBoolVector(len(inp)),
		}
	}

	const CELL_SIZE = 64
	var wg sync.WaitGroup

	worker := func(l, r int) {
		defer wg.Done()
		for i := l; i < r; i++ {
			secret := int(util.MustParseInt(inp[i]))
			prices[i] = append(prices[i], secret%10)

			for j := range 2000 {
				newSecret := secret
				newSecret = (newSecret ^ (newSecret * 64)) % 16777216
				newSecret = (newSecret ^ (newSecret / 32)) % 16777216
				newSecret = (newSecret ^ (newSecret * 2048)) % 16777216

				if newSecret < 0 {
					panic(newSecret)
				}

				prices[i] = append(prices[i], newSecret%10)
				secret = newSecret

				if j >= 3 {
					ch1, ch2, ch3, ch4 :=
						prices[i][j-2]-prices[i][j-3]+9,
						prices[i][j-1]-prices[i][j-2]+9,
						prices[i][j]-prices[i][j-1]+9,
						prices[i][j+1]-prices[i][j]+9
					ss := ch1 + 19*(ch2+19*(ch3+19*ch4))
					if !sales[ss].buyers.Get(i) {
						atomic.AddInt64(&sales[ss].count, int64(prices[i][j+1]))
						sales[ss].buyers.Set(i)
					}
				}
			}

			atomic.AddInt64(&silver, int64(secret))
		}
	}

	for i := 0; i < len(inp); i += CELL_SIZE {
		j := min(len(inp), i+CELL_SIZE)
		wg.Add(1)
		go worker(i, j)
	}

	wg.Wait()

	for _, st := range sales {
		gold = max(gold, st.count)
	}

	fmt.Fprintln(out, "silver:", silver)
	fmt.Fprintln(out, "gold:", gold)
}
