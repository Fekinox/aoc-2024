# Advent of Code 2024

All solutions to Advent of Code 2024, written in Go.

## Usage

Visit the [Advent of Code website](https://adventofcode.com/2024), download the
input, and save it to `input/aocXX/input`, where XX is the problem number.

Build the program with `go build .` and run it with `aoc-2024`. By default, it
will run all 25 solvers on the default input (`input/aocXX/input`) and return
the output accordingly. Provide input flags to change the default behavior:

* `-p n`: Only run solvers for a specific problem number.
* `-i inputFile`: Solvers should look for the input file located at
  `input/aocXX/inputFile` instead.
* `-t trials`: Repeatedly run the solution for the given number of trials and
  return the min, median, mean, max, and standard deviation of the run times.
