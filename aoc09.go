package main

import (
	"fmt"
	"io"
	// "github.com/Fekinox/aoc-2024/util"
)

func Problem09MultiHeaps(in io.Reader, out io.Writer) {
	var silver, gold int64

	data, _ := io.ReadAll(in)
	data = data[:len(data)-1]

	var sPos, gPos int64

	done := make(chan struct{})

	moved := make([]bool, len(data)/2+1)
	heaps := make([][]int, 9)
	go func() {
		for i := 0; i < len(data); i += 2 {
			c := data[i]
			if c >= '1' && c <= '9' {
				v := int(c - '1')
				heaps[v] = append(heaps[v], i/2)
			}
		}
		done <- struct{}{}
	}()

	rightPointer := len(data) - 1
	rightCount := int(data[rightPointer] - '0')
	for i, c := range data {
		id := i / 2
		// even: data block
		// skip if already moved
		if rightPointer < i {
			break
		}
		if rightPointer == i {
			for range rightCount {
				silver += int64(id) * sPos
				sPos++
			}
			break
		}
		if i%2 == 0 {
			for range int(c - '0') {
				silver += int64(id) * sPos
				sPos++
			}
			continue
		}
		// odd: free space
	outer:
		for range int(c - '0') {
			if rightPointer < i {
				sPos++
				break outer
			}
			for rightCount == 0 {
				rightPointer -= 2
				if rightPointer < i {
					sPos++
					break outer
				}
				rightCount = int(data[rightPointer] - '0')
			}
			rightCount--
			silver += int64(rightPointer/2) * sPos
			sPos++
		}
	}

	<-done

	for i, c := range data {
		id := i / 2
		// even: data block
		// skip if already moved
		if i%2 == 0 {
			if moved[id] {
				gPos += int64(c - '0')
			} else {
				for range int(c - '0') {
					gold += int64(id) * gPos
					gPos++
				}
			}
			continue
		}

		// odd: free space
		remainingSpace := int(c - '0')
		for remainingSpace > 0 {
			bestId := -1
			bestSize := 0
			for m := range min(9, remainingSpace) {
				if len(heaps[m]) == 0 {
					continue
				}
				top := heaps[m][len(heaps[m])-1]
				if top*2 < i {
					continue
				}
				if top > bestId {
					bestId = top
					bestSize = m
				}
			}

			if bestId == -1 {
				break
			}

			for range bestSize + 1 {
				gold += int64(bestId) * gPos
				gPos++
			}

			moved[bestId] = true
			remainingSpace -= bestSize + 1
			heaps[bestSize] = heaps[bestSize][:len(heaps[bestSize])-1]
		}
		gPos += int64(remainingSpace)
	}

	fmt.Fprintln(out, "silver:", silver)
	fmt.Fprintln(out, "gold:", gold)
}

func Problem09Unwashed(in io.Reader, out io.Writer) {
	var silver, gold int64

	data, _ := io.ReadAll(in)
	data = data[:len(data)-1]

	var sPos, gPos int64

	rightPointer := len(data) - 1
	rightCount := int(data[rightPointer] - '0')
	for i, c := range data {
		id := i / 2
		// even: data block
		// skip if already moved
		if rightPointer < i {
			break
		}
		if rightPointer == i {
			for range rightCount {
				silver += int64(id) * sPos
				sPos++
			}
			break
		}
		if i%2 == 0 {
			for range int(c - '0') {
				silver += int64(id) * sPos
				sPos++
			}
			continue
		}
		// odd: free space
	outer:
		for range int(c - '0') {
			if rightPointer < i {
				sPos++
				break outer
			}
			for rightCount == 0 {
				rightPointer -= 2
				if rightPointer < i {
					sPos++
					break outer
				}
				rightCount = int(data[rightPointer] - '0')
			}
			rightCount--
			silver += int64(rightPointer/2) * sPos
			sPos++
		}
	}

	moved := make([]bool, len(data)/2+1)
	for i, c := range data {
		id := i / 2
		// even: data block
		// skip if already moved
		if i%2 == 0 {
			if moved[id] {
				gPos += int64(c - '0')
			} else {
				for range int(c - '0') {
					gold += int64(id) * gPos
					gPos++
				}
			}
			continue
		}

		// odd: free space
		remainingSpace := int(c - '0')
	goldOuter:
		for remainingSpace > 0 {
			rp := len(data) - 1
			for {
				if rp < i {
					break goldOuter
				}
				if int(data[rp]-'0') == 0 {
					rp -= 2
					continue
				}
				if !moved[rp/2] && int(data[rp]-'0') <= remainingSpace {
					break
				}
				rp -= 2
			}

			for range int(data[rp] - '0') {
				gold += int64(rp/2) * gPos
				gPos++
			}

			moved[rp/2] = true
			remainingSpace -= int(data[rp] - '0')
		}
		gPos += int64(remainingSpace)
	}

	fmt.Fprintln(out, "silver:", silver)
	fmt.Fprintln(out, "gold:", gold)
}

func Problem09(in io.Reader, out io.Writer) {
	Problem09MultiHeaps(in, out)
}
