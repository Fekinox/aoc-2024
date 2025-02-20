package main

import (
	"fmt"
	"io"
	"slices"
	"strings"

	"github.com/Fekinox/aoc-2024/util"
)

func Problem24Washed(in io.Reader, out io.Writer) {
	var silver int64
	var gold string

	g := util.ReadNewlineSeparatedGroups(in)

	type gate struct {
		a string
		op string
		b string
	}

	var swaps []string
	gateToRegister := make(map[gate]string)
	registerToGate := make(map[string]gate)

	registerQuery := func(a, op, b string) string {
		var l, r string
		if a < b {
			l, r = a, b
		} else {
			l, r = b, a
		}
		return gateToRegister[gate{a: l, op: op, b: r}]
	}

	swap := func(a, b string) bool {
		if slices.Contains(swaps, a) || slices.Contains(swaps, b) {
			return false
		}
		swaps = append(swaps, a, b)
		var aGate, bGate gate
		for gt, out := range gateToRegister {
			if out == a {
				aGate = gt
			}
			if out == b {
				bGate = gt
			}
		}
		gateToRegister[aGate], gateToRegister[bGate] = b, a
		registerToGate[a], registerToGate[b] = bGate, aGate
		return true
	}

	registers := make(map[string]bool)
	var zRegisters int
	for _, l := range g[0] {
		tokens := strings.Split(l, ": ")
		registers[tokens[0]] = tokens[1] == "1"
	}

	for _, l := range g[1] {
		tokens := strings.Split(l, " ")
		var a, b, c, op string
		c = tokens[4]
		op = tokens[1]
		if tokens[0] < tokens[2] {
			a, b = tokens[0], tokens[2]
		} else {
			a, b = tokens[2], tokens[0]
		}

		gt := gate{
			a: a,
			op: op,
			b: b,
		}
		gateToRegister[gt] = c
		registerToGate[c] = gt

		// z registers only appear as outputs
		if c[0] == 'z' && c[1] >= '0' && c[1] <= '9' {
			zRegisters++
		}
	}

	// silver
	usedGates := make(map[gate]struct{})
	for {
		var dirty bool
		for gt, out := range gateToRegister {
			if _, ok := usedGates[gt]; ok {
				continue
			}
			l, ok := registers[gt.a]
			if !ok {
				continue
			}
			r, ok := registers[gt.b]
			if !ok {
				continue
			}
			switch gt.op {
			case "AND":
				registers[out] = l && r
			case "OR":
				registers[out] = l || r
			case "XOR":
				registers[out] = l != r
			}
			usedGates[gt] = struct{}{}
			dirty = true
		}
		if !dirty {
			break
		}
	}

	for rg, set := range registers {
		if !set || rg[0] != 'z' || rg[1] < '0' || rg[1] > '9' {
			continue
		}
		pos := util.MustParseInt(rg[1:])
		silver += 1 << pos
	}

	// gold
	restartCarry, restartIndex := "", -1
	for len(swaps) < 8 {
		var carry string
		var i int
		if restartIndex != -1 {
			i = restartIndex
			carry = restartCarry
		}
		for i < zRegisters {
			// Invariant of this loop is that the carry bit is correct.
			xi, yi, zi := fmt.Sprintf("x%02d", i), fmt.Sprintf("y%02d", i), fmt.Sprintf("z%02d", i)
			// If i = 0, carry bit is x00 & y00 and sum bit is x00 ^ y00.
			// First carry bit is correct so invariant holds.
			var adderout string
			if i == 0 {
				carry = registerQuery(xi, "AND", yi)
				adderout = registerQuery(xi, "XOR", yi)
			} else {
				// Otherwise, xi ^ yi is the half-sum, which must be XORed with the carry bit
				// (correct by invariant) to produce zi.
				halfsum := registerQuery(xi, "XOR", yi)
				adderout = registerQuery(halfsum, "XOR", carry)
				if adderout != "" {
					// The next carry bit should be the OR of xi & yi and halfsum & carry.
					// NOTE: why is this carry bit guaranteed to be correct? why can't
					// any of these have the wrong output?
					var c1, c2 string
					c1 = registerQuery(xi, "AND", yi)
					c2 = registerQuery(halfsum, "AND", carry)
					carry = registerQuery(c1, "OR", c2)
				} else {
					// If no such gate halfsum ^ carry exists, then we must replace the output
					// of xi ^ yi. We know halfsum specifically is incorrect because the carry
					// bit is correct by the invariant
					// zi must receive xi ^ yi and carry as inputs, so whichever one of those
					// inputs is not carry must be the actual output of xi ^ yi.
					gateOutput := registerToGate[zi]
					var actualAdder string
					if gateOutput.a == carry {
						actualAdder = gateOutput.b
					} else {
						actualAdder = gateOutput.a
					}
					swap(actualAdder, halfsum)
					break
				}
			}
			if adderout != zi {
				swap(adderout, zi)
				break
			}
			i++
			restartCarry, restartIndex = carry, i
		}
	}

	slices.Sort(swaps)
	gold = strings.Join(swaps, ",")
	
	fmt.Fprintln(out, "silver:", silver)
	fmt.Fprintln(out, "gold:", gold)
}

func Problem24Unwashed(in io.Reader, out io.Writer) {
	var silver int64
	var gold string

	g := util.ReadNewlineSeparatedGroups(in)
	registers := make(map[string]bool)

	type gate struct {
		a string
		b string
		c string
		op string
	}

	for _, l := range g[0] {
		tokens := strings.Split(l, ": ")
		registers[tokens[0]] = tokens[1] == "1"
	}

	var gates, remainingGates []gate
	var halfsumCarryGates, halfcarryOverGates []int

	var potentialCarries []string

	for i, l := range g[1] {
		tokens := strings.Split(l, " ")
		gt := gate{
			a: tokens[0],
			b: tokens[2],
			c: tokens[4],
			op: tokens[1],
		}
		remainingGates = append(remainingGates, gt)
		gates = append(gates, gt)

		if gt.a[1] < '0' || gt.a[1] > '9' {
			if gt.op == "OR" {
				halfcarryOverGates = append(halfcarryOverGates, i)
			} else {
				halfsumCarryGates = append(halfsumCarryGates, i)
				if gt.op == "XOR" {
					potentialCarries = append(potentialCarries, gt.a)
					potentialCarries = append(potentialCarries, gt.b)
				}

			}
		}
	}

	slices.Sort(potentialCarries)

	for len(remainingGates) > 0 {
		var newGates []gate
		for _, g := range remainingGates {
			a, ok := registers[g.a]
			if !ok {
				newGates = append(newGates, g)
				continue
			}
			b, ok := registers[g.b]
			if !ok {
				newGates = append(newGates, g)
				continue
			}
			switch g.op {
			case "AND":
				registers[g.c] = a && b
			case "OR":
				registers[g.c] = a || b
			case "XOR":
				registers[g.c] = a != b
			}
		}

		remainingGates = newGates
		fmt.Println(len(remainingGates))
	}

	var zTokens []string
	for z := range registers {
		if z[0] == 'z' {
			zTokens = append(zTokens, z)
		}
	}
	slices.Sort(zTokens)
	for i, z := range zTokens {
		fmt.Println(z, registers[z])
		if registers[z] {
			silver += 1<<i
		}
	}

	// a00 = x00 ^ y00
	// b00 = x00 & y00
	//
	// a01 = x01 ^ y01
	// b01 = x01 & y01
	//
	// c01 = a01 ^ b00
	//
	// d01 = a01 & b00
	// d01 = (x01 ^ y01) & b00
	// e01 = b01 | d01
	// e01 = (x01 & y01) | (d01)
	//
	// skt = x00 & y00 (b00)
	// kjn = x01 ^ y01 (a01)
	// hth = x01 & y01 (b01)
	// z01 = kjn ^ skt (c01)
	// hjc = skt & kjn (d01)
	// pgc = hth | hjc (e01)
	//
	// given
	// CIN = prev step
	// SS1 = xkk ^ ykk
	// CC1 = xkk & ykk
	// then
	// ZKK = SS1 ^ CIN
	// DDD = SS1 & CIN
	// EEE = DDD | CC1
	//
	//
	// (not last or first bit)
	// ... = CARRYIN
	// 1HALFSUM ^ 0CARRYIN = SUM
	// HALFCARRY | OVER = CARRYOUT
	// 1HALFSUM & 0CARRYIN = OVER
	// Xi & Yi = HALFCARRY
	// Xi ^ Yi = HALFSUM
	//
	// (first bit)
	// Xi ^ Yi = SUM
	// Xi & Yi = CARRYOUT
	//
	// (last bit)
	//
	// if starting at 0, then ZKK = 
	// only the outputs are potentially wrong for these 5
	//
	// you know what the first carry bit is
	//
	// the strings HALFSUM and CARRYIN are both non-terminal (don't contain numbers) and exist
	// in XOR-AND pairs
	// each CARRYIN must pair with a corresponding CARRYOUT.
	// for z00, SUM = HALFSUM. for z44, CARRYOUT should just be z45
	// HALFSUM ^ CARRYIN must output to its corresponding terminal- if not, it must be swapped
	// HALFSUM & CARRYIN must output to one of (HALFCARRY, OVER)- if not, it must be swapped
	// HALFCARRY and OVER only appear once, HALFSUM, CARRYIN, XI, and YI all appear twice
	//
	// with the carry bit, you can deduce the halfsum's register.
	// - if there is a match but it doesn't output to the terminal, the terminal has to be swapped
	// - if there are no matches, the carry bit has to be swapped
	// Xi & Yi and HALFSUM & CARRYIN should output into HALFCARRY and OVER. if not, you have
	// to swap them
	// - if Xi & Yi and HALFSUM & CARRYIN both output in there, we're good
	// - if Xi & Yi doesn't output in there but HALFSUM & CARRYIN does, then HALFSUM & CARRYIN
	// correctly maps to OVER and Xi & Yi's output must be swapped
	// - ditto
	// - if Xi & Yi and HALFSUM & CARRYIN are both wrong, they both need to be swapped (doesn't matter
	// the order)
	
	// var carry string
	// carry = rightCarries[0]
	
	type state struct {
		carry string
		terminal int
		swaps []string
	}

	fmt.Println(potentialCarries)

	var initCarry string
	for _, g := range gates {
		if g.op == "AND" && (g.a == "x00" || g.a == "y00") {
			initCarry = g.c
		}
	}

	stack := []state{{
		carry: initCarry,
		terminal: 1,
		swaps: []string{},
	}}

	for len(stack) > 0 {
		cur := stack[0]
		stack = stack[1:]
		fmt.Println(cur)

		withSwaps := func(s string) string {
			for i := 0; i < len(cur.swaps); i+=2 {
				if s == cur.swaps[i] {
					return cur.swaps[i+1]
				} else if s == cur.swaps[i+1] {
					return cur.swaps[i]
				}
			}
			return s
		}
		addSwaps := func(s state, a, b string) (state, bool) {
			if slices.Contains(cur.swaps, a) || slices.Contains(cur.swaps, b) {
				return state{}, false
			}
			var sw []string
			sw = append(sw, s.swaps...)
			sw = append(sw, a, b)
			return state{
				carry: s.carry,
				swaps: sw,
				terminal: s.terminal,
			}, true
		}

		// Having more than 4 swaps is impossible
		if len(cur.swaps) > 8 {
			continue
		}

		// If carry is a terminal that isn't z45, impossible
		if cur.carry[1] >= '0' && cur.carry[1] <= '9' {
			if cur.carry == "z45" {
				slices.Sort(cur.swaps)
				gold = strings.Join(cur.swaps, ",")
				break
			}
		}

		// Attempt to find HALFSUM by finding a XOR gate with carry as one of the inputs
		var halfsumXor, halfsumAnd *gate
		var halfsum string
		for _, i := range halfsumCarryGates {
			g := &gates[i]
			if g.a == cur.carry || g.b == cur.carry {
				if g.op == "XOR" {
					halfsumXor = g
				} else {
					halfsumAnd = g
				}
				if g.a == cur.carry {
					halfsum = g.b
				} else {
					halfsum = g.a
				}
			}
		}

		// If none are found, swap carry with one of the available carries
		if halfsumXor == nil || halfsumAnd == nil {
			fmt.Println("could not find carry- retrying with all swaps")
			for _, c := range potentialCarries {
				if c == cur.carry {
					continue
				}
				sw, ok := addSwaps(cur, c, cur.carry)
				if !ok {
					continue
				}
				stack = append(stack, state{
					carry: c,
					swaps: sw.swaps,
					terminal: sw.terminal,
				})
			}
			continue
		}

		// If halfsumXor does not output to the expected terminal, swap accordingly
		expectedTerminal := fmt.Sprintf("z%02d", cur.terminal)
		expectedInput := fmt.Sprintf("x%02d", cur.terminal)
		if withSwaps(halfsumXor.c) != expectedTerminal {
			fmt.Println("unexpected terminal", withSwaps(halfsumXor.c), expectedTerminal)
			if sw, ok := addSwaps(cur, withSwaps(halfsumXor.c), expectedTerminal); ok {
				stack = append(stack, sw)
			}
			continue
		}

		// If xii ^ yii does not output to halfsum, swap accordingly
		var rightHalfsum, rightHalfcarry *gate
		for i := range gates {
			g := &gates[i]
			if (g.a == expectedInput || g.b == expectedInput) {
				if g.op == "AND" {
					rightHalfcarry = g
				} else {
					rightHalfsum = g
				}
			}
		}
		if withSwaps(rightHalfsum.c) != halfsum {
			fmt.Println("halfsum incorrect")
			if sw, ok := addSwaps(cur, halfsum, withSwaps(rightHalfsum.c)); ok {
				stack = append(stack, sw)
			}
			continue
		}

		// The outputs of halfsumAnd and rightHalfCarry should both go to an OR gate.
		var carryOut *gate
		over, halfcarry := 
			withSwaps(halfsumAnd.c), withSwaps(rightHalfcarry.c)
		var trueOver, trueHalfCarry string
		for _, i := range halfcarryOverGates {
			g := &gates[i]
			if g.a == over || g.b == halfcarry {
				carryOut = g
				trueOver = g.a
				trueHalfCarry = g.b
			}
			if g.b == over || g.a == halfcarry {
				carryOut = g
				trueOver = g.b
				trueHalfCarry = g.a
			}
		}

		if over != trueOver && halfcarry != trueHalfCarry {
			fmt.Println("uuuuhhhhh")
			continue
		}

		if halfcarry != trueHalfCarry {
			fmt.Println("half carry incorrect output")
			if sw, ok := addSwaps(cur, halfcarry, trueHalfCarry); ok {
				stack = append(stack, sw)
			}
			continue
		}

		if over != trueOver {
			fmt.Println("over incorrect")
			if sw, ok := addSwaps(cur, over, trueOver); ok {
				stack = append(stack, sw)
			}
			continue
		}

		fmt.Println(halfsumAnd, halfsumXor, rightHalfsum, rightHalfcarry, carryOut)

		stack = append(stack, state{
			carry: withSwaps(carryOut.c),
			terminal: cur.terminal+1,
			swaps: cur.swaps,
		})
	}
	
	fmt.Fprintln(out, "silver:", silver)
	fmt.Fprintln(out, "gold:", gold)
}

func Problem24(in io.Reader, out io.Writer) {
	Problem24Washed(in, out)
}
