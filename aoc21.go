package main

import (
	"bytes"
	"fmt"
	"io"
	"math"
	"math/big"

	"github.com/Fekinox/aoc-2024/util"
)

func Problem21Bignum(in io.Reader, out io.Writer) {
	silver, gold := big.NewInt(0), big.NewInt(0)

	type pos struct {
		X int
		Y int
	}

	keypadPos := map[byte]pos{
		'0': {X: 1, Y: 3},
		'A': {X: 2, Y: 3},
		'1': {X: 0, Y: 2},
		'2': {X: 1, Y: 2},
		'3': {X: 2, Y: 2},
		'4': {X: 0, Y: 1},
		'5': {X: 1, Y: 1},
		'6': {X: 2, Y: 1},
		'7': {X: 0, Y: 0},
		'8': {X: 1, Y: 0},
		'9': {X: 2, Y: 0},
	}
	controllerPos := map[byte]pos{
		'^': {X: 1, Y: 0},
		'A': {X: 2, Y: 0},
		'<': {X: 0, Y: 1},
		'v': {X: 1, Y: 1},
		'>': {X: 2, Y: 1},
	}

	keypadRawTransitions := util.MakeGrid(11, 11, [][]byte{})
	controllerRawTransitions := util.MakeGrid(5, 5, [][]byte{})

	keypadByteToInt := func(b byte) int {
		switch b {
		case 'A':
			return 10
		default:
			return int(b - '0')
		}
	}
	controllerByteToInt := func(b byte) int {
		switch b {
		case 'A':
			return 0
		case '^':
			return 1
		case '<':
			return 2
		case 'v':
			return 3
		default:
			return 4
		}
	}
	rawTransitions := func(m map[byte]pos, src, dst byte, forbidden pos) [][]byte {
		res := [][]byte{nil}
		srcPos, dstPos := m[src], m[dst]
		var hc, vc byte
		if srcPos.X > dstPos.X {
			hc = '<'
		} else {
			hc = '>'
		}
		if srcPos.Y > dstPos.Y {
			vc = '^'
		} else {
			vc = 'v'
		}

		delX, delY := util.Abs(srcPos.X-dstPos.X), util.Abs(srcPos.Y-dstPos.Y)
		var dx, dy int
		if srcPos.X > dstPos.X {
			dx = -1
		} else {
			dx = 1
		}
		if srcPos.Y > dstPos.Y {
			dy = -1
		} else {
			dy = 1
		}

		for range delX + delY {
			var nextRes [][]byte
			for _, p := range res {
				if delX > 0 {
					nextRes = append(nextRes,
						append(bytes.Clone(p), hc),
					)
				}
				if delY > 0 {
					nextRes = append(nextRes,
						append(bytes.Clone(p), vc),
					)
				}
			}
			res = nextRes
		}

		var filteredRes [][]byte
	outer:
		for _, p := range res {
			var hct, vct int
			for _, b := range p {
				if b == hc {
					hct++
				} else {
					vct++
				}
				c := pos{X: srcPos.X + hct*dx, Y: srcPos.Y + vct*dy}
				if c == forbidden {
					continue outer
				}
			}
			if hct == delX && vct == delY {
				filteredRes = append(filteredRes, append(bytes.Clone(p), 'A'))
			}
		}

		return filteredRes
	}

	for _, sb := range []byte("0123456789A") {
		for _, db := range []byte("0123456789A") {
			ts := rawTransitions(keypadPos, sb, db, pos{X: 0, Y: 3})
			keypadRawTransitions.Set(keypadByteToInt(sb), keypadByteToInt(db), ts)
		}
	}

	for _, sb := range []byte("A^<v>") {
		for _, db := range []byte("A^<v>") {
			ts := rawTransitions(controllerPos, sb, db, pos{X: 0, Y: 0})
			controllerRawTransitions.Set(controllerByteToInt(sb), controllerByteToInt(db), ts)
		}
	}

	silverMatrix := util.MakeGridWith(11, 11, func(x, y int) *big.Int {
		return big.NewInt(0)
	})
	goldMatrix := util.MakeGridWith(11, 11, func(x, y int) *big.Int {
		return big.NewInt(0)
	})
	controllerMatrix := util.MakeGridWith(5, 5, func(x, y int) *big.Int {
		return big.NewInt(1)
	})

	for i := range 25000 {
		nextMatrix := util.MakeGrid[*big.Int](5, 5, nil)
		for _, sb := range []byte("A^<v>") {
			for _, db := range []byte("A^<v>") {
				si, di := controllerByteToInt(sb), controllerByteToInt(db)
				var bestDist *big.Int
				for _, tr := range controllerRawTransitions.MustGet(si, di) {
					dist := big.NewInt(0)
					for i, ds := range tr {
						var sr byte = 'A'
						if i > 0 {
							sr = tr[i-1]
						}
						sri, dsi := controllerByteToInt(sr), controllerByteToInt(ds)
						dist.Add(dist, controllerMatrix.MustGet(sri, dsi))
					}
					if bestDist == nil || bestDist.Cmp(dist) > 0 {
						bestDist = dist
					}
				}
				nextMatrix.Set(si, di, bestDist)
			}
		}
		controllerMatrix = nextMatrix

		if i == 1 {
			for _, sb := range []byte("0123456789A") {
				for _, db := range []byte("0123456789A") {
					si, di := keypadByteToInt(sb), keypadByteToInt(db)
					var bestDist *big.Int
					for _, tr := range keypadRawTransitions.MustGet(si, di) {
						dist := big.NewInt(0)
						for i, ds := range tr {
							var sr byte = 'A'
							if i > 0 {
								sr = tr[i-1]
							}
							sri, dsi := controllerByteToInt(sr), controllerByteToInt(ds)
							dist.Add(dist, controllerMatrix.MustGet(sri, dsi))
						}
						if bestDist == nil || bestDist.Cmp(dist) > 0 {
							bestDist = dist
						}
					}
					silverMatrix.Set(si, di, bestDist)
				}
			}
		}
	}

	for _, sb := range []byte("0123456789A") {
		for _, db := range []byte("0123456789A") {
			si, di := keypadByteToInt(sb), keypadByteToInt(db)
			var bestDist *big.Int
			for _, tr := range keypadRawTransitions.MustGet(si, di) {
				dist := big.NewInt(0)
				for i, ds := range tr {
					var sr byte = 'A'
					if i > 0 {
						sr = tr[i-1]
					}
					sri, dsi := controllerByteToInt(sr), controllerByteToInt(ds)
					dist.Add(dist, controllerMatrix.MustGet(sri, dsi))
				}
				if bestDist == nil || bestDist.Cmp(dist) > 0 {
					bestDist = dist
				}
			}
			goldMatrix.Set(si, di, bestDist)
		}
	}

	for _, l := range util.ReadLines(in) {
		silverP, goldP := big.NewInt(0), big.NewInt(0)
		for i, dst := range []byte(l) {
			var src byte = 'A'
			if i > 0 {
				src = byte(l[i-1])
			}
			si, di := keypadByteToInt(src), keypadByteToInt(dst)
			goldP.Add(goldP, goldMatrix.MustGet(si, di))
			silverP.Add(silverP, silverMatrix.MustGet(si, di))
		}
		code1 := big.NewInt(0)
		code2 := big.NewInt(0)
		code1.UnmarshalText([]byte(l[:len(l)-1]))
		code2.UnmarshalText([]byte(l[:len(l)-1]))
		silver.Add(silver, code1.Mul(code1, silverP))
		gold.Add(gold, code2.Mul(code2, goldP))
	}

	fmt.Fprintln(out, "silver:", silver)
	fmt.Fprintln(out, "gold:", gold)
}

func Problem21Main(in io.Reader, out io.Writer) {
	var silver, gold int64

	type pos struct {
		X int
		Y int
	}

	keypadPos := map[byte]pos{
		'0': {X: 1, Y: 3},
		'A': {X: 2, Y: 3},
		'1': {X: 0, Y: 2},
		'2': {X: 1, Y: 2},
		'3': {X: 2, Y: 2},
		'4': {X: 0, Y: 1},
		'5': {X: 1, Y: 1},
		'6': {X: 2, Y: 1},
		'7': {X: 0, Y: 0},
		'8': {X: 1, Y: 0},
		'9': {X: 2, Y: 0},
	}
	controllerPos := map[byte]pos{
		'^': {X: 1, Y: 0},
		'A': {X: 2, Y: 0},
		'<': {X: 0, Y: 1},
		'v': {X: 1, Y: 1},
		'>': {X: 2, Y: 1},
	}

	keypadRawTransitions := util.MakeGrid(11, 11, [][]byte{})
	controllerRawTransitions := util.MakeGrid(5, 5, [][]byte{})

	keypadByteToInt := func(b byte) int {
		switch b {
		case 'A':
			return 10
		default:
			return int(b - '0')
		}
	}
	controllerByteToInt := func(b byte) int {
		switch b {
		case 'A':
			return 0
		case '^':
			return 1
		case '<':
			return 2
		case 'v':
			return 3
		default:
			return 4
		}
	}
	rawTransitions := func(m map[byte]pos, src, dst byte, forbidden pos) [][]byte {
		var res [][]byte
		srcPos, dstPos := m[src], m[dst]
		var hc, vc byte
		if srcPos.X > dstPos.X {
			hc = '<'
		} else {
			hc = '>'
		}
		if srcPos.Y > dstPos.Y {
			vc = '^'
		} else {
			vc = 'v'
		}

		delX, delY := util.Abs(srcPos.X-dstPos.X), util.Abs(srcPos.Y-dstPos.Y)
		var dx, dy int
		if srcPos.X > dstPos.X {
			dx = -1
		} else {
			dx = 1
		}
		if srcPos.Y > dstPos.Y {
			dy = -1
		} else {
			dy = 1
		}

		hFirst, vFirst := make([]byte, delX+delY), make([]byte, delX+delY)
		for i := range delX + delY {
			if i < delX {
				hFirst[i] = hc
			} else {
				hFirst[i] = vc
			}
			if i < delY {
				vFirst[i] = vc
			} else {
				vFirst[i] = hc
			}
		}
		if dx == 0 || dy == 0 {
			return [][]byte{hFirst}
		}
		hCorner, vCorner :=
			pos{
				X: srcPos.X + dx*delX,
				Y: srcPos.Y,
			},
			pos{
				X: srcPos.X,
				Y: srcPos.Y + dy*delY,
			}

		if hCorner != forbidden {
			res = append(res, append(hFirst, 'A'))
		}
		if vCorner != forbidden {
			res = append(res, append(vFirst, 'A'))
		}

		return res
	}

	for _, sb := range []byte("0123456789A") {
		for _, db := range []byte("0123456789A") {
			ts := rawTransitions(keypadPos, sb, db, pos{X: 0, Y: 3})
			keypadRawTransitions.Set(keypadByteToInt(sb), keypadByteToInt(db), ts)
		}
	}

	for _, sb := range []byte("A^<v>") {
		for _, db := range []byte("A^<v>") {
			ts := rawTransitions(controllerPos, sb, db, pos{X: 0, Y: 0})
			controllerRawTransitions.Set(controllerByteToInt(sb), controllerByteToInt(db), ts)
		}
	}

	silverMatrix := util.MakeGrid(11, 11, int64(0))
	goldMatrix := util.MakeGrid(11, 11, int64(0))
	controllerMatrix := util.MakeGrid(5, 5, int64(1))

	for i := range 25 {
		nextMatrix := util.MakeGrid(5, 5, int64(0))
		for _, sb := range []byte("A^<v>") {
			for _, db := range []byte("A^<v>") {
				si, di := controllerByteToInt(sb), controllerByteToInt(db)
				var bestDist int64 = math.MaxInt64
				for _, tr := range controllerRawTransitions.MustGet(si, di) {
					var dist int64
					for i, ds := range tr {
						var sr byte = 'A'
						if i > 0 {
							sr = tr[i-1]
						}
						sri, dsi := controllerByteToInt(sr), controllerByteToInt(ds)
						dist += controllerMatrix.MustGet(sri, dsi)
					}
					bestDist = min(bestDist, dist)
				}
				nextMatrix.Set(si, di, bestDist)
			}
		}
		controllerMatrix = nextMatrix

		if i == 1 {
			for _, sb := range []byte("0123456789A") {
				for _, db := range []byte("0123456789A") {
					si, di := keypadByteToInt(sb), keypadByteToInt(db)
					var bestDist int64 = math.MaxInt64
					for _, tr := range keypadRawTransitions.MustGet(si, di) {
						var dist int64
						for i, ds := range tr {
							var sr byte = 'A'
							if i > 0 {
								sr = tr[i-1]
							}
							sri, dsi := controllerByteToInt(sr), controllerByteToInt(ds)
							dist += controllerMatrix.MustGet(sri, dsi)
						}
						bestDist = min(bestDist, dist)
					}
					silverMatrix.Set(si, di, bestDist)
				}
			}
		}
	}

	for _, sb := range []byte("0123456789A") {
		for _, db := range []byte("0123456789A") {
			si, di := keypadByteToInt(sb), keypadByteToInt(db)
			var bestDist int64 = math.MaxInt64
			for _, tr := range keypadRawTransitions.MustGet(si, di) {
				var dist int64
				for i, ds := range tr {
					var sr byte = 'A'
					if i > 0 {
						sr = tr[i-1]
					}
					sri, dsi := controllerByteToInt(sr), controllerByteToInt(ds)
					dist += controllerMatrix.MustGet(sri, dsi)
				}
				bestDist = min(bestDist, dist)
			}
			goldMatrix.Set(si, di, bestDist)
		}
	}

	for _, l := range util.ReadLines(in) {
		var silverP, goldP int64
		for i, dst := range []byte(l) {
			var src byte = 'A'
			if i > 0 {
				src = byte(l[i-1])
			}
			goldP += goldMatrix.MustGet(keypadByteToInt(src), keypadByteToInt(dst))
			silverP += silverMatrix.MustGet(keypadByteToInt(src), keypadByteToInt(dst))
		}
		code := util.MustParseInt(l[:len(l)-1])
		silver += code * silverP
		gold += code * goldP
	}

	fmt.Fprintln(out, "silver:", silver)
	fmt.Fprintln(out, "gold:", gold)
}

func Problem21Naive(in io.Reader, out io.Writer) {
	var silver, gold int64

	type pos struct {
		X int
		Y int
	}

	type state struct {
		rob1  pos
		rob2  pos
		rob3  pos
		index int
	}

	keypad := map[pos]byte{
		{X: 0, Y: 0}: '7',
		{X: 1, Y: 0}: '8',
		{X: 2, Y: 0}: '9',
		{X: 0, Y: 1}: '4',
		{X: 1, Y: 1}: '5',
		{X: 2, Y: 1}: '6',
		{X: 0, Y: 2}: '1',
		{X: 1, Y: 2}: '2',
		{X: 2, Y: 2}: '3',
		{X: 1, Y: 3}: '0',
		{X: 2, Y: 3}: 'A',
	}
	keypadA := pos{X: 2, Y: 3}

	controller := map[pos]byte{
		{X: 1, Y: 0}: '^',
		{X: 2, Y: 0}: 'A',
		{X: 0, Y: 1}: '<',
		{X: 1, Y: 1}: 'v',
		{X: 2, Y: 1}: '>',
	}
	controllerA := pos{X: 2, Y: 0}

	directions := map[byte]pos{
		'^': {X: 0, Y: -1},
		'<': {X: -1, Y: 0},
		'>': {X: 1, Y: 0},
		'v': {X: 0, Y: 1},
	}

	for _, l := range util.ReadLines(in) {
		visited := make(map[state][]byte, 275*4)
		frontier := map[state][]byte{{
			rob1:  pos{X: 2, Y: 3},
			rob2:  pos{X: 2, Y: 0},
			rob3:  pos{X: 2, Y: 0},
			index: 0,
		}: nil}

		var dist int
		var reached bool
	outer:
		for len(frontier) > 0 {
			nextFrontier := make(map[state][]byte)
			for cur, path := range frontier {
				if _, ok := visited[cur]; ok {
					continue
				}
				visited[cur] = path
				// test all legal moves
				// directional for robot 3
				for dc, d := range directions {
					pp := pos{
						X: cur.rob3.X + int(d.X),
						Y: cur.rob3.Y + int(d.Y),
					}
					p := append(bytes.Clone(path), dc)
					if _, ok := controller[pp]; ok {
						nextFrontier[state{
							rob1:  cur.rob1,
							rob2:  cur.rob2,
							rob3:  pp,
							index: cur.index,
						}] = p
					}
				}
				// activate
				p := append(bytes.Clone(path), 'A')
				// case 1: rob3 is on a direction
				if cur.rob3 != controllerA {
					pp := pos{
						X: cur.rob2.X + directions[controller[cur.rob3]].X,
						Y: cur.rob2.Y + directions[controller[cur.rob3]].Y,
					}
					if _, ok := controller[pp]; ok {
						nextFrontier[state{
							rob1:  cur.rob1,
							rob2:  pp,
							rob3:  cur.rob3,
							index: cur.index,
						}] = p
					}
					continue
				}
				// case 2: rob3 is on activate, rob 2 on direction
				if cur.rob2 != controllerA {
					pp := pos{
						X: cur.rob1.X + directions[controller[cur.rob2]].X,
						Y: cur.rob1.Y + directions[controller[cur.rob2]].Y,
					}
					if _, ok := keypad[pp]; ok {
						nextFrontier[state{
							rob1:  pp,
							rob2:  cur.rob2,
							rob3:  cur.rob3,
							index: cur.index,
						}] = p
					}
					continue
				}
				// case 3: rob2 and 3 on activate, rob 1 on number we want
				if cur.rob1 != keypadA && keypad[cur.rob1] == byte(l[cur.index]) {
					nextFrontier[state{
						rob1:  cur.rob1,
						rob2:  cur.rob2,
						rob3:  cur.rob3,
						index: cur.index + 1,
					}] = p
					continue
				}
				// case 4: all robs on activate and index is at the end
				if cur.rob1 == keypadA && l[cur.index] == 'A' {
					dist++
					reached = true
					break outer
				}
			}
			frontier = nextFrontier
			dist++
		}
		fmt.Println(dist, reached, util.MustParseInt(l[:len(l)-1]))

		code := util.MustParseInt(l[:len(l)-1])
		silver += code * int64(dist)
	}

	fmt.Fprintln(out, "silver:", silver)
	fmt.Fprintln(out, "gold:", gold)
}

func Problem21(in io.Reader, out io.Writer) {
	Problem21Main(in, out)
}
