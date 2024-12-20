package main

import (
	"fmt"
	"io"

	"github.com/Fekinox/aoc-2024/util"
)

func ProblemXX(in io.Reader, out io.Writer) {
	input := util.ReadLines(in)

	for _, l := range input {
		fmt.Fprintln(out, l)
	}
}
