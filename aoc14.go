package main

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"strings"

	"github.com/Fekinox/aoc-2024/util"
)

func Problem14(in io.Reader, out io.Writer) {
	var silver, gold int64

	sc := bufio.NewScanner(in)

	const WIDTH int64 = 606
	const HEIGHT int64 = 666

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
	minOverlap := math.MaxInt
	argminOverlap := int64(-1)
	minQuadrantMetric := math.MaxInt
	argminQuadrantMetric := int64(-1)
	for i := range WIDTH * HEIGHT {
		fmt.Println(i)
		var pairwiseDistance, overlaps, quadrantMetric int

		overlapMap := make(map[pos]int)
		fps := make([]pos, len(robots))
		quadrants := make([]int, 4)

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
			overlapMap[fps[j]]++

			if fps[j].X != int(WIDTH/2) && fps[j].Y != int(HEIGHT/2) {
				q := 0
				if fps[j].X > int(WIDTH/2) {
					q += 1
				}
				if fps[j].Y > int(HEIGHT/2) {
					q += 2
				}
				quadrants[q]++
			}
		}
		for _, o := range overlapMap {
			if o > 1 {
				overlaps += o - 1
			}
		}
		quadrantMetric += util.Abs(quadrants[0] - quadrants[1])
		quadrantMetric += util.Abs(quadrants[1] - quadrants[2])
		quadrantMetric += util.Abs(quadrants[2] - quadrants[3])
		quadrantMetric += util.Abs(quadrants[3] - quadrants[0])
		quadrantMetric += 2 * util.Abs(quadrants[0]-quadrants[2])
		quadrantMetric += 2 * util.Abs(quadrants[1]-quadrants[3])
		if pairwiseDistance < minPairwiseDistance {
			argminPairwiseDistance = i
			minPairwiseDistance = pairwiseDistance
		}
		if overlaps < minOverlap {
			argminOverlap = i
			minOverlap = overlaps
		}
		if quadrantMetric < minQuadrantMetric {
			argminQuadrantMetric = i
			minQuadrantMetric = quadrantMetric
		}

	}

	fmt.Fprintln(out, argminOverlap, argminPairwiseDistance, argminQuadrantMetric)
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

	fmt.Println(b)

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
