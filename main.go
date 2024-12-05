package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"math"
	"os"
	"slices"
	"time"
)

var inputFile = flag.String("i", "input", "Input file")
var problem = flag.Int("p", -1, "Problem number (omit to run all problems)")
var trials = flag.Int("t", 1, "Number of trials")

type ProblemFunc func(in io.Reader, out io.Writer)

type Stats struct {
	Min    float64
	Median float64
	Mean   float64
	Max    float64
	StdDev float64
}

var problems = []ProblemFunc{
	Problem01,
	Problem02,
	Problem03,
	Problem04,
	Problem05,
	Problem06,
	Problem07,
	Problem08,
	Problem09,
	Problem10,
	Problem11,
	Problem12,
	Problem13,
	Problem14,
	Problem15,
	Problem16,
	Problem17,
	Problem18,
	Problem19,
	Problem20,
	Problem21,
	Problem22,
	Problem23,
	Problem24,
	Problem25,
}

// Wrap a problem to catch panics and recover gracefully
func WrapProblem(p ProblemFunc, problemNumber int, input string) (runtime float64, err error) {
	var startTime time.Time
	defer func() {
		if r := recover(); r != nil {
			runtime = float64(time.Now().Sub(startTime).Nanoseconds()) / (1000 * 1000)
			err = fmt.Errorf("Panic: %v", r)
		}
	}()

	dest := fmt.Sprintf("input/aoc%02d/%v", problemNumber, input)
	fmt.Printf("Opening %v...\n", dest)
	f, err := os.Open(dest)
	if err != nil {
		return 0.0, err
	}
	defer f.Close()

	startTime = time.Now()
	p(f, os.Stdout)
	return float64(time.Now().Sub(startTime).Nanoseconds()) / (1000 * 1000), nil
}

func RunProblem(prob int, input string) {
	fmt.Printf("Running problem %v...\n", prob+1)
	runtime, err := WrapProblem(problems[prob], prob+1, input)
	if err != nil {
		fmt.Println(err)
		return
	}
	if *trials <= 1 {
		fmt.Printf("Runtime: %v ms\n", runtime)
	} else {
		stats, err := Profile(problems[prob], prob+1, input)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("Trials:", *trials)
		fmt.Println("Min:", stats.Min)
		fmt.Println("Median:", stats.Median)
		fmt.Println("Mean:", stats.Mean)
		fmt.Println("Max:", stats.Max)
		fmt.Println("StdDev:", stats.StdDev)
	}
}

func Profile(p ProblemFunc, problemNumber int, input string) (stats Stats, err error) {
	count := 0
	mean := 0.0
	m2 := 0.0
	median := 0.0
	ts := make([]float64, *trials)

	f, err := os.Open(fmt.Sprintf("input/aoc%02d/%v", problemNumber, input))
	if err != nil {
		return Stats{}, err
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		return Stats{}, err
	}

	for i := range *trials {
		r := bytes.NewReader(data)
		startTime := time.Now()
		p(r, io.Discard)
		samp := float64(time.Now().Sub(startTime).Nanoseconds()) / (1000 * 1000)
		ts[i] = samp

		count++
		delta := samp - mean
		mean += delta / float64(count)
		delta2 := samp - mean
		m2 = delta * delta2

		fmt.Printf("Trial %v: %v ms\n", i+1, samp)
	}

	sampleVariance := m2 / float64(count)

	slices.Sort(ts)

	if *trials%2 == 0 {
		median = (ts[*trials/2] + ts[*trials/2+1]) / 2
	} else {
		median = ts[*trials/2]
	}

	return Stats{
		Min:    ts[0],
		Median: median,
		Mean:   mean,
		Max:    ts[*trials-1],
		StdDev: math.Sqrt(sampleVariance),
	}, nil
}

func main() {
	flag.Parse()

	if *problem == -1 {
		for i := range problems {
			RunProblem(i, *inputFile)
		}
	} else if *problem >= 1 && *problem <= 25 {
		RunProblem(*problem-1, *inputFile)
	} else {
		flag.Usage()
	}
}
