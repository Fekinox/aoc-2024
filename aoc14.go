package main

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"strings"

	"github.com/Fekinox/aoc-2024/util"
)

func Problem14CRT(in io.Reader, out io.Writer) {
	var silver, gold int64

	sc := bufio.NewScanner(in)

	const WIDTH int64 = 101
	const HEIGHT int64 = 103

	type pos struct {
		X int
		Y int
	}

	type robot struct {
		Pos pos
		Vel pos
	}

	quadrants := []int64{0, 0, 0, 0}
	var robots []robot

	for sc.Scan() {
		vectors := strings.Split(sc.Text(), " ")
		p, v := util.SplitNumbers(vectors[0][2:], ","), util.SplitNumbers(vectors[1][2:], ",")
		ps := pos{
			X: int(p[0]),
			Y: int(p[1]),
		}
		vl := pos{
			X: int(v[0]),
			Y: int(v[1]),
		}

		if vl.X < 0 {
			vl.X += int(WIDTH)
		}
		if vl.Y < 0 {
			vl.Y += int(HEIGHT)
		}

		robots = append(robots, robot{Pos: ps, Vel: vl})

		fp := pos{
			X: (ps.X + vl.X*100) % int(WIDTH),
			Y: (ps.Y + vl.Y*100) % int(HEIGHT),
		}

		q := 0
		if fp.X != int(WIDTH)/2 && fp.Y != int(HEIGHT)/2 {
			if fp.X > int(WIDTH)/2 {
				q += 2
			}
			if fp.Y > int(HEIGHT)/2 {
				q += 1
			}
			quadrants[q]++
		}

	}

	minXVariance, minYVariance := math.MaxInt, math.MaxInt
	argminXVariance, argminYVariance := int64(-1), int64(-1)
	for i := range max(WIDTH, HEIGHT) {
		var xvar, yvar int

		fps := make([]pos, len(robots))

		for j, r := range robots {
			fps[j] = pos{
				X: (r.Pos.X + r.Vel.X*int(i)) % int(WIDTH),
				Y: (r.Pos.Y + r.Vel.Y*int(i)) % int(HEIGHT),
			}
			if fps[j].X < 0 || fps[j].Y < 0 {
				panic("what")
			}
			for k := range j {
				xvar += util.Abs(fps[k].X-fps[j].X)
				yvar += util.Abs(fps[k].Y-fps[j].Y)
			}

		}
		if xvar < minXVariance {
			argminXVariance = i
			minXVariance = xvar 
		}
		if yvar < minYVariance {
			argminYVariance = i
			minYVariance = yvar 
		}
	}

	// We know that t ~ t1 mod 101 and t ~ t2 mod 103
	// By first congruence, t = t1 + 101k
	// so t1 + 101k ~ t2 mod 103
	// 101k ~ t2 - t1 mod 103
	// k ~ (101)^-1 * (t2 - t1) mod 103
	// note that (101)^-1 ~ 51 mod 103, as
	// 101 ~ -2 mod 103
	// 51 * 101 ~ -102 ~ 1 mod 103
	// so k = 51 * (xvar + yvar) + 103k
	k := (51 * (argminYVariance - argminXVariance)) % 103
	for k < 0 {
		k += 103
	}
	gold = (argminXVariance + 101*k)

	g := util.MakeGrid(int(WIDTH), int(HEIGHT), false)
	b := 0
	for _, r := range robots {
		fp := pos{
			X: int((int64(r.Pos.X) + int64(r.Vel.X)*(gold+WIDTH*HEIGHT)) % WIDTH),
			Y: int((int64(r.Pos.Y) + int64(r.Vel.Y)*(gold+WIDTH*HEIGHT)) % HEIGHT),
		}
		g.Set(fp.X, fp.Y, true)
		b++
	}

	for y := range g.Height {
		for x := range g.Width {
			ch := '.'
			if g.MustGet(x, y) {
				ch = '#'
			}
			fmt.Fprintf(out, "%c", ch)
		}
		fmt.Fprintf(out, "\n")
	}

	silver = 1
	for _, v := range quadrants {
		silver *= v
	}

	fmt.Fprintln(out, "silver:", silver)
	fmt.Fprintln(out, "gold:", gold)
}

func Problem14Naive(in io.Reader, out io.Writer) {
	var silver, gold int64

	sc := bufio.NewScanner(in)

	const WIDTH int64 = 101
	const HEIGHT int64 = 103

	type pos struct {
		X int
		Y int
	}

	type robot struct {
		Pos pos
		Vel pos
	}

	quadrants := []int64{0, 0, 0, 0}
	var robots []robot

	for sc.Scan() {
		vectors := strings.Split(sc.Text(), " ")
		p, v := util.SplitNumbers(vectors[0][2:], ","), util.SplitNumbers(vectors[1][2:], ",")
		ps := pos{
			X: int(p[0]),
			Y: int(p[1]),
		}
		vl := pos{
			X: int(v[0]),
			Y: int(v[1]),
		}

		if vl.X < 0 {
			vl.X += int(WIDTH)
		}
		if vl.Y < 0 {
			vl.Y += int(HEIGHT)
		}

		robots = append(robots, robot{Pos: ps, Vel: vl})

		fp := pos{
			X: (ps.X + vl.X*100) % int(WIDTH),
			Y: (ps.Y + vl.Y*100) % int(HEIGHT),
		}

		q := 0
		if fp.X != int(WIDTH)/2 && fp.Y != int(HEIGHT)/2 {
			if fp.X > int(WIDTH)/2 {
				q += 2
			}
			if fp.Y > int(HEIGHT)/2 {
				q += 1
			}
			quadrants[q]++
		}

	}

	minPairwiseDistance := math.MaxInt
	argminPairwiseDistance := int64(-1)
	for i := range WIDTH * HEIGHT {
		var pairwiseDistance int

		fps := make([]pos, len(robots))

		for j, r := range robots {
			fps[j] = pos{
				X: (r.Pos.X + r.Vel.X*int(i)) % int(WIDTH),
				Y: (r.Pos.Y + r.Vel.Y*int(i)) % int(HEIGHT),
			}
			if fps[j].X < 0 || fps[j].Y < 0 {
				panic("what")
			}
			for k := range j {
				pairwiseDistance += util.Abs(fps[k].X-fps[j].X) + util.Abs(fps[k].Y-fps[j].Y)
			}

		}
		if pairwiseDistance < minPairwiseDistance {
			argminPairwiseDistance = i
			minPairwiseDistance = pairwiseDistance
		}

	}

	gold = int64(argminPairwiseDistance)

	g := util.MakeGrid(int(WIDTH), int(HEIGHT), false)
	b := 0
	for _, r := range robots {
		fp := pos{
			X: int((int64(r.Pos.X) + int64(r.Vel.X)*(argminPairwiseDistance+WIDTH*HEIGHT)) % WIDTH),
			Y: int((int64(r.Pos.Y) + int64(r.Vel.Y)*(argminPairwiseDistance+WIDTH*HEIGHT)) % HEIGHT),
		}
		g.Set(fp.X, fp.Y, true)
		b++
	}

	for y := range g.Height {
		for x := range g.Width {
			ch := '.'
			if g.MustGet(x, y) {
				ch = '#'
			}
			fmt.Fprintf(out, "%c", ch)
		}
		fmt.Fprintf(out, "\n")
	}

	silver = 1
	for _, v := range quadrants {
		silver *= v
	}

	fmt.Fprintln(out, "silver:", silver)
	fmt.Fprintln(out, "gold:", gold)
}

func Problem14(in io.Reader, out io.Writer) {
	Problem14CRT(in, out)
}
