package main

import (
	"container/heap"
	"fmt"
	"io"
	"math"
	"sync"

	"github.com/Fekinox/aoc-2024/util"
)

func Problem16Grid(in io.Reader, out io.Writer) {
	var silver, gold int64
	g := util.GridFromReader(in)

	type pos struct {
		X int
		Y int
	}

	dirs := []pos{
		{X: 1, Y: 0},
		{X: 0, Y: -1},
		{X: -1, Y: 0},
		{X: 0, Y: 1},
	}

	type state struct {
		X int
		Y int
		D int
	}

	var count int64
	var sx, sy, ex, ey int
	for y := range g.Height {
		for x := range g.Width {
			if g.MustGet(x, y) == 'S' {
				sx, sy = x, y
			}
			if g.MustGet(x, y) == 'E' {
				ex, ey = x, y
			}
			if g.MustGet(x, y) != '#' {
				count++
			}
		}
	}

	startState := state{
		X: sx,
		Y: sy,
		D: 0,
	}

	var endStates []state
	for d := range dirs {
		endStates = append(endStates, state{
			X: ex,
			Y: ey,
			D: d,
		})
	}

	sSPD := make([]int64, g.Width*g.Height*4)
	eSPD := make([]int64, g.Width*g.Height*4)

	for i := range sSPD {
		sSPD[i] = -1
		eSPD[i] = -1
	}

	var wg sync.WaitGroup
	wg.Add(2)

	dijk := func(distmap []int64, initStates ...state) {
		defer wg.Done()
		var frontier util.PriorityQueue[int64, state]

		for _, st := range initStates {
			heap.Push(&frontier, &util.PriorityQueueItem[int64, state]{
				Value:    st,
				Priority: 0,
			})
		}

		// var vis int
		for frontier.Len() > 0 {
			// fmt.Fprintf(out, "%d/%d\n", vis, count*4)
			it := heap.Pop(&frontier).(*util.PriorityQueueItem[int64, state])
			idx := it.Value.D + 4*(it.Value.X+g.Width*it.Value.Y)
			if v := distmap[idx]; v != -1 {
				continue
			}
			// vis++
			distmap[idx] = it.Priority
			dr := dirs[it.Value.D]
			nx, ny := it.Value.X+dr.X, it.Value.Y+dr.Y
			if v, ok := g.Get(nx, ny); ok && v != '#' {
				heap.Push(&frontier, &util.PriorityQueueItem[int64, state]{
					Value: state{
						X: nx,
						Y: ny,
						D: it.Value.D,
					},
					Priority: it.Priority + 1,
				})
			}
			heap.Push(&frontier, &util.PriorityQueueItem[int64, state]{
				Value: state{
					X: it.Value.X,
					Y: it.Value.Y,
					D: (it.Value.D + 1) % 4,
				},
				Priority: it.Priority + 1000,
			})
			heap.Push(&frontier, &util.PriorityQueueItem[int64, state]{
				Value: state{
					X: it.Value.X,
					Y: it.Value.Y,
					D: (it.Value.D + 3) % 4,
				},
				Priority: it.Priority + 1000,
			})
		}
	}

	go dijk(sSPD, startState)
	go dijk(eSPD, endStates...)
	wg.Wait()

	silver = int64(math.MaxInt64)
	for d := range util.CardinalDirections {
		idx := d + 4*(ex+g.Width*ey)
		if v := sSPD[idx]; v != -1 {
			silver = min(silver, v)
		}
	}

	for y := range g.Height {
		for x := range g.Width {
			if v, ok := g.Get(x, y); !ok || v == '#' {
				continue
			}
			best := int64(math.MaxInt64)
			for d := range util.CardinalDirections {
				sidx := d + 4*(x+g.Width*y)
				eidx := ((d + 2) % 4) + 4*(x+g.Width*y)

				sDist := sSPD[sidx]
				eDist := eSPD[eidx]
				if sDist != -1 && eDist != -1 {
					best = min(best, sDist+eDist)
				}
			}
			if best < silver {
				panic("impossible")
			}
			if best == silver {
				gold++
			}
		}
	}

	fmt.Fprintln(out, "silver:", silver)
	fmt.Fprintln(out, "gold:", gold)
}

func Problem16Hashmap(in io.Reader, out io.Writer) {
	var silver, gold int64
	g := util.GridFromReader(in)

	type pos struct {
		X int
		Y int
	}

	dirs := []pos{
		{X: 1, Y: 0},
		{X: 0, Y: -1},
		{X: -1, Y: 0},
		{X: 0, Y: 1},
	}

	type state struct {
		X int
		Y int
		D int
	}

	var count int64
	var sx, sy, ex, ey int
	for y := range g.Height {
		for x := range g.Width {
			if g.MustGet(x, y) == 'S' {
				sx, sy = x, y
			}
			if g.MustGet(x, y) == 'E' {
				ex, ey = x, y
			}
			if g.MustGet(x, y) != '#' {
				count++
			}
		}
	}

	fmt.Fprintln(out, count)

	startState := state{
		X: sx,
		Y: sy,
		D: 0,
	}

	var endStates []state
	for d := range dirs {
		endStates = append(endStates, state{
			X: ex,
			Y: ey,
			D: d,
		})
	}

	//sspd[state]: shortest path from (S, east) to (state.pos, state.dir)
	sShortestPathDistance := make(map[state]int64)
	//espd[state]: shortest path from (state.pos, state.dir) to (E, any pos)
	eShortestPathDistance := make(map[state]int64)

	var wg sync.WaitGroup
	wg.Add(2)

	dijk := func(distmap map[state]int64, initStates ...state) {
		defer wg.Done()
		var frontier util.PriorityQueue[int64, state]

		for _, st := range initStates {
			heap.Push(&frontier, &util.PriorityQueueItem[int64, state]{
				Value:    st,
				Priority: 0,
			})
		}

		for frontier.Len() > 0 {
			it := heap.Pop(&frontier).(*util.PriorityQueueItem[int64, state])
			if _, ok := distmap[it.Value]; ok {
				continue
			}
			distmap[it.Value] = it.Priority
			dr := dirs[it.Value.D]
			nx, ny := it.Value.X+dr.X, it.Value.Y+dr.Y
			if v, ok := g.Get(nx, ny); ok && v != '#' {
				heap.Push(&frontier, &util.PriorityQueueItem[int64, state]{
					Value: state{
						X: nx,
						Y: ny,
						D: it.Value.D,
					},
					Priority: it.Priority + 1,
				})
			}
			heap.Push(&frontier, &util.PriorityQueueItem[int64, state]{
				Value: state{
					X: it.Value.X,
					Y: it.Value.Y,
					D: (it.Value.D + 1) % 4,
				},
				Priority: it.Priority + 1000,
			})
			heap.Push(&frontier, &util.PriorityQueueItem[int64, state]{
				Value: state{
					X: it.Value.X,
					Y: it.Value.Y,
					D: (it.Value.D + 3) % 4,
				},
				Priority: it.Priority + 1000,
			})
		}
	}

	go dijk(sShortestPathDistance, startState)
	go dijk(eShortestPathDistance, endStates...)
	wg.Wait()

	silver = int64(math.MaxInt64)
	for d := range util.CardinalDirections {
		if v, ok := sShortestPathDistance[state{
			X: ex,
			Y: ey,
			D: d,
		}]; ok {
			silver = min(silver, v)
		}
	}

	for y := range g.Height {
		for x := range g.Width {
			if v, ok := g.Get(x, y); !ok || v == '#' {
				continue
			}
			best := int64(math.MaxInt64)
			for d := range util.CardinalDirections {
				st := state{
					X: x,
					Y: y,
					D: d,
				}
				rst := state{
					X: x,
					Y: y,
					D: (d + 2) % 4,
				}

				sDist := sShortestPathDistance[st]
				eDist := eShortestPathDistance[rst]
				best = min(best, sDist+eDist)
			}
			if best < silver {
				panic("impossible")
			}
			if best == silver {
				gold++
			}
		}
	}

	fmt.Fprintln(out, "silver:", silver)
	fmt.Fprintln(out, "gold:", gold)
}

func Problem16(in io.Reader, out io.Writer) {
	Problem16Grid(in, out)
}
