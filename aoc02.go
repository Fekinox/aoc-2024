package main

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/Fekinox/aoc-2024/util"
)

func Problem2Regular(l string) (silver, gold bool) {
	tokens := strings.Fields(l)
	nums := make([]int, len(tokens))
	deltas := make([]int, len(tokens)-1)

	for i, t := range tokens {
		nums[i] = int(util.MustParseInt(t))
	}

	for i := range deltas {
		deltas[i] = nums[i+1] - nums[i]
	}

	var sign int
	if deltas[0] > 0 {
		sign = 1
	} else {
		sign = -1
	}

	silverError := false
	for _, d := range deltas {
		if sign*d < 0 || util.Abs(d) < 1 || util.Abs(d) > 3 {
			silverError = true
			break
		}
	}

	if !silverError {
		return true, true
	}

	for i := range nums {
		goldError := false
		var newSign int
		var firstDelta int

		if i == 0 {
			firstDelta = deltas[1]
		} else if i == 1 {
			firstDelta = deltas[0] + deltas[1]
		} else {
			firstDelta = deltas[0]
		}

		if firstDelta > 0 {
			newSign = 1
		} else {
			newSign = -1
		}

		for j := range len(nums) - 1 {
			// if removing last number, ignore last delta
			if i == len(nums)-1 && j == len(nums)-2 {
				continue
			}
			// if we're on a number we removed, skip it
			if j == i {
				continue
			}

			// consider the delta at the current number,
			// and if the number after the current one is the
			// removed number, then skip ahead
			del := deltas[j]
			if j+1 == i {
				del += deltas[i]
			}

			if newSign*del < 0 || util.Abs(del) < 1 || util.Abs(del) > 3 {
				goldError = true
				break
			}
		}

		if !goldError {
			return false, true
		}
	}

	return false, false
}

func Problem02(in io.Reader, out io.Writer) {
	safeSilver := int64(0)
	safeGold := int64(0)

	scanner := bufio.NewScanner(in)
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1024*1024)

outer:
	for scanner.Scan() {
		l := scanner.Text()
		tokens := strings.Fields(l)
		nums := make([]int, len(tokens))
		deltas := make([]int, len(tokens)-1)

		for i, t := range tokens {
			nums[i] = int(util.MustParseInt(t))
			if i > 0 {
				deltas[i-1] = nums[i] - nums[i-1]
			}
		}

		// max i where all deltas from [0..i-1] are safe (all same sign, abs value between 1 and 3
		longestPrefix := 0
		// min j where all deltas from [len(deltas)-1-j..len(deltas)-1] are safe
		longestSuffix := 0

		var longestPrefixDone, longestSuffixDone bool

		for i := range deltas {
			end := len(deltas) - 1 - i
			// prefix
			if !longestPrefixDone && util.Abs(deltas[i]) >= 1 && util.Abs(deltas[i]) <= 3 &&
				(i == 0 || deltas[i]*deltas[i-1] > 0) {
				longestPrefix++
			} else if !longestPrefixDone {
				longestPrefixDone = true
			}

			// suffix
			if !longestSuffixDone && util.Abs(deltas[end]) >= 1 && util.Abs(deltas[end]) <= 3 &&
				(end == len(deltas)-1 || deltas[end]*deltas[end+1] > 0) {
				longestSuffix++
			} else if !longestSuffixDone {
				longestSuffixDone = true
			}

			if longestPrefixDone && longestSuffixDone {
				break
			}
		}

		// if this is the case, then every single delta is safe
		if longestPrefix == len(deltas) {
			safeSilver++
			safeGold++
			continue outer
		}

		// if we're one short, then we can remove the number at the end to cover the entire
		// list of deltas
		if longestPrefix == len(deltas)-1 || longestSuffix == len(deltas)-1 {
			safeGold++
			continue outer
		}

		// otherwise, if we remove the ith number (1 <= i < len(deltas)-1, then
		// newDeltas[i-1] = x[i+1] - x[i-1]
		//				  = x[i+1] - x[i] + x[i] - x[i-1]
		//				  = d[i] - d[i-1]
		// we have to check the following:
		// (newDeltas[0], newDeltas[1], ... newDeltas[i-2] are all safe and have the same sign
		// newDeltas[i+1], newDeltas[i+2], ... newDeltas[len(deltas)-1] are all safe and have the same sign
		// absolute value of newDeltas[i-1] is between 1 and 3
		// everything has the same sign
		// we can rewrite the first two conditions to be in terms of i, which means we can place
		// them in the loop bounds
		loops := 0
		for i := max(1, len(nums)-2-longestSuffix); i <= min(len(nums)-2, longestPrefix+1); i++ {
			newDelta := deltas[i] + deltas[i-1]
			if util.Abs(newDelta) >= 1 && util.Abs(newDelta) <= 3 &&
				(i == 1 || newDelta*deltas[0] > 0) &&
				(i == len(nums)-2 || newDelta*deltas[len(deltas)-1] > 0) {
				safeGold++
				fmt.Fprintln(out, loops)
				fmt.Fprintln(out, min(len(nums)-2, longestPrefix+1)-max(1, len(nums)-2-longestSuffix)+1)
				continue outer
			}
			loops++
		}
	}

	fmt.Fprintln(out, "silver:", safeSilver)
	fmt.Fprintln(out, "gold:", safeGold)
}
