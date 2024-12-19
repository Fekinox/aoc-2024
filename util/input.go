package util

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

func ReadLines(r io.Reader) []string {
	res := make([]string, 0)
	scanner := bufio.NewScanner(r)
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1024*1024)
	for scanner.Scan() {
		res = append(res, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}

	return res
}

type GroupReader struct {
	scanner *bufio.Scanner
	currentGroup []string
	err error
}

func NewGroupReader(r io.Reader) *GroupReader {
	scanner := bufio.NewScanner(r)
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1024*1024)
	return &GroupReader{
		scanner: scanner,
		currentGroup: []string{},
	}
}

func (g *GroupReader) Scan() bool {
	g.currentGroup = make([]string, 0)
	for g.scanner.Scan() {
		if err := g.scanner.Err(); err != nil {
			g.err = err
			return false
		}
		t := g.scanner.Text()
		if t == "" {
			return true
		}
		g.currentGroup = append(g.currentGroup, t)
	}
	return false
}

func (g *GroupReader) Group() []string {
	return g.currentGroup
}

func (g *GroupReader) Err() error {
	return g.err
}

func ReadNewlineSeparatedGroups(r io.Reader) [][]string {
	res := make([][]string, 0)
	curGroup := make([]string, 0)
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		t := scanner.Text()
		if t == "" {
			res = append(res, curGroup)
			curGroup = make([]string, 0)
			continue
		}
		curGroup = append(curGroup, t)
	}

	res = append(res, curGroup)

	return res
}

func ReadString(r io.Reader) string {
	var buf strings.Builder
	_, err := io.Copy(&buf, r)
	if err != nil {
		panic(err)
	}
	return buf.String()
}

func ReadGrid(r io.Reader) Grid[rune] {
	return GridFromStrings(ReadLines(r)...)
}

func SplitNumbers(s string, sep string) []int64 {
	splits := strings.Split(s, sep)
	nums := make([]int64, len(splits))
	for i, sp := range splits {
		nums[i] = MustParseInt(sp)
	}

	return nums
}
