package util

import (
	"bufio"
	"io"
	"os"
	"strings"
)

func ReadLines(path string) []string {
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	res := make([]string, 0)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		res = append(res, scanner.Text())
	}

	return res
}

func ReadNewlineSeparatedGroups(path string) [][]string {
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	res := make([][]string, 0)
	curGroup := make([]string, 0)
	scanner := bufio.NewScanner(f)
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

func ReadString(path string) string {
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	var buf strings.Builder
	_, err = io.Copy(&buf, f)
	if err != nil {
		panic(err)
	}
	return buf.String()
}

func ReadGrid(path string) Grid[rune] {
	return GridFromStrings(ReadLines(path)...)
}
