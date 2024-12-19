package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"runtime"
	"runtime/pprof"
	"slices"
	"time"
)

var inputFile = flag.String("i", "input", "Input file")
var problem = flag.Int("p", -1, "Problem number (omit to run all problems)")
var trials = flag.Int("t", 1, "Number of trials")
var variant = flag.String("v", "default", "Problem variant to run")

var cpuprof = flag.String("cpuprofile", "", "write cpu profile to `file`")
var memprof = flag.String("memprofile", "", "write mem profile to `file`")

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

var variants = map[int]map[string]ProblemFunc{
	6: {
		"gosper":  Problem06ComplexG,
		"brent":   Problem06ComplexB,
		"floyd":   Problem06ComplexTH,
		"simple":  Problem06Simple,
		"complex": Problem06Complex,
	},
	7: {
		"bfs":   Problem07BFS,
		"dfs":   Problem07DFS,
		"stack": Problem07Stack,
	},
	8: {
		"hashmap": Problem08Hashmap,
		"array":   Problem08Array,
	},
	9: {
		"normal": Problem09Unwashed,
		"heaps":  Problem09MultiHeaps,
	},
	11: {
		"string":   Problem11String,
		"nostring": Problem11NoString,
	},
	13: {
		"regex":   Problem13Regex,
		"noregex": Problem13NoRegex,
	},
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
		stats, err := Profile(problems[prob], prob+1, input, false)
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

func RunProblemVariant(prob int, input string, variant string, silent bool) {
	varList, ok := variants[prob+1]
	if !ok {
		fmt.Println("Solution does not have additional variants")
		return
	}
	if variant == "all" {
		for vn := range varList {
			RunProblemVariant(prob, input, vn, true)
		}
		return
	}
	vr, ok := varList[variant]
	if !ok {
		fmt.Printf("Variant %v does not exist\n", variant)
		return
	}
	fmt.Printf("Running problem %v... (variant %v)\n", prob+1, variant)
	runtime, err := WrapProblem(vr, prob+1, input)
	if err != nil {
		fmt.Println(err)
		return
	}
	if *trials <= 1 {
		fmt.Printf("Runtime: %v ms\n", runtime)
	} else {
		stats, err := Profile(vr, prob+1, input, silent)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("Trials: %v ms\n", *trials)
		fmt.Printf("Min: %v ms\n", stats.Min)
		fmt.Printf("Median: %v ms\n", stats.Median)
		fmt.Printf("Mean: %v ms\n", stats.Mean)
		fmt.Printf("Max: %v ms\n", stats.Max)
		fmt.Printf("StdDev: %v ms\n", stats.StdDev)
	}
}

func Profile(p ProblemFunc, problemNumber int, input string, silent bool) (stats Stats, err error) {
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

		if !silent {
			fmt.Printf("Trial %v: %v ms\n", i+1, samp)
		}
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
	if *cpuprof != "" {
		f, err := os.Create(*cpuprof)
		if err != nil {
			log.Fatalln("Could not create CPU profile", err)
		}
		defer f.Close()
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatalln("could not start CPU profile", err)
		}
		defer pprof.StopCPUProfile()
	}

	if *problem == -1 {
		for i := range problems {
			RunProblem(i, *inputFile)
		}
	} else if *problem >= 1 && *problem <= 25 {
		if *variant != "default" {
			RunProblemVariant(*problem-1, *inputFile, *variant, false)
		} else {
			RunProblem(*problem-1, *inputFile)
		}
	} else {
		flag.Usage()
	}

	if *memprof != "" {
		f, err := os.Create(*memprof)
		if err != nil {
			log.Fatalln("Could not create memory profile", err)
		}
		defer f.Close()
		runtime.GC()
		if err := pprof.WriteHeapProfile(f); err != nil {
			log.Fatalln("Could not write memory profile", err)
		}
	}
}
