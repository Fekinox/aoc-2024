package main

import (
	"fmt"
	"io"
	"math"
	"strings"

	"github.com/Fekinox/aoc-2024/util"
)

func Problem17Washed(in io.Reader, out io.Writer) {
	var silver string
	var gold int64

	lines := util.ReadLines(in)

	var program []int64
	a := util.MustParseInt(strings.Split(lines[0], " ")[2])
	program = util.SplitNumbers(strings.Split(lines[4], " ")[1], ",")
	var output []int64
	k1, k2 := program[3], program[11]

	for a > 0 {
		b := (a % 8) ^ k1
		c := (a >> b)
		a = a / 8
		output = append(output, (b^c^k2)%8)
	}

	type item struct {
		currentA int64
		index    int
	}

	stack := []item{{
		currentA: 0,
		index:    len(program) - 1,
	}}

	targetA := int64(math.MaxInt64)
outer:
	for len(stack) > 0 {
		cur := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		// for range len(program) -  1 - cur.index {
		// 	fmt.Fprintf(out, "  ")
		// }
		// cc := cur.currentA
		// values := make([]byte, 16)
		// for i := range 16 {
		// 	values[len(values)-1-i] = byte(cc % 8)
		// 	cc /= 8
		// }
		// fmt.Fprintln(out, values)

		for t := range 8 {
			tt := int64(7 - t)
			testVal := cur.currentA*8 + tt
			pv := program[cur.index]

			b2 := pv ^ k2
			b1 := tt ^ k1
			if b2%8 == (b1^(testVal>>b1))%8 {
				if cur.index == 0 {
					targetA = min(targetA, testVal)
					break outer
				} else {
					stack = append(stack, item{
						currentA: testVal,
						index:    cur.index - 1,
					})
				}
			}
		}
	}

	var sb strings.Builder
	for i, v := range output {
		if i == len(output)-1 {
			fmt.Fprintf(&sb, "%v", v)
		} else {
			fmt.Fprintf(&sb, "%v,", v)
		}
	}

	silver = sb.String()
	gold = targetA

	fmt.Fprintln(out, "silver:", silver)
	fmt.Fprintln(out, "gold:", gold)
}

func Problem17Unwashed(in io.Reader, out io.Writer) {
	var silver string
	var gold int64

	lines := util.ReadLines(in)

	var a, b, c int64
	var program []int64
	a = util.MustParseInt(strings.Split(lines[0], " ")[2])
	b = util.MustParseInt(strings.Split(lines[1], " ")[2])
	c = util.MustParseInt(strings.Split(lines[2], " ")[2])
	program = util.SplitNumbers(strings.Split(lines[4], " ")[1], ",")
	var ip int

	var output []int64

	for {
		fmt.Println(a, b, c)
		if ip >= len(program) {
			break
		}
		op := program[ip]
		var literalValue int64 = program[ip+1]
		var comboValue int64
		if literalValue <= 3 {
			comboValue = literalValue
		} else if literalValue == 4 {
			comboValue = a
		} else if literalValue == 5 {
			comboValue = b
		} else if literalValue == 6 {
			comboValue = c
		}

		switch op {
		case 0:
			a = a >> comboValue
			ip += 2
		case 1:
			b = b ^ literalValue
			ip += 2
		case 2:
			b = comboValue % 8
			ip += 2
		case 3:
			if a == 0 {
				ip += 2
			} else {
				ip = int(literalValue)
			}
		case 4:
			b = b ^ c
			ip += 2
		case 5:
			output = append(output, comboValue%8)
			ip += 2
		case 6:
			b = a >> comboValue
			ip += 2
		case 7:
			c = a >> comboValue
			ip += 2
		}
	}

	type item struct {
		currentA int64
		index    int
	}

	stack := []item{{
		currentA: 0,
		index:    len(program) - 1,
	}}

	k1, k2 := program[3], program[11]

	targetA := int64(math.MaxInt64)
	var count int
	for len(stack) > 0 {
		cur := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		count++

		for t := range 8 {
			tt := int64(t)
			testVal := cur.currentA*8 + tt
			pv := program[cur.index]

			b2 := pv ^ k2
			b1 := tt ^ k1
			if b2%8 == (b1^(testVal>>b1))%8 {
				if cur.index == 0 {
					targetA = min(targetA, testVal)
				} else {
					stack = append(stack, item{
						currentA: testVal,
						index:    cur.index - 1,
					})
				}
			}
		}
	}

	fmt.Fprintln(out, count)

	var sb strings.Builder
	for i, v := range output {
		if i == len(output)-1 {
			fmt.Fprintf(&sb, "%v", v)
		} else {
			fmt.Fprintf(&sb, "%v,", v)
		}
	}

	silver = sb.String()
	gold = targetA

	fmt.Fprintln(out, "silver:", silver)
	fmt.Fprintln(out, "gold:", gold)
}

func Problem17(in io.Reader, out io.Writer) {
	Problem17Washed(in, out)
}
