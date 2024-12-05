package main

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strings"
	"sync"

	"github.com/Fekinox/aoc-2024/util"
)

var mulRegex = regexp.MustCompile("mul\\((\\d+),(\\d+)\\)")

func addMuls(data []byte) int64 {
	var res int64
	for _, m := range mulRegex.FindAllSubmatch(data, -1) {
		fmt.Println(string(m[0]))
		res +=  util.MustParseInt(string(m[1]))
			util.MustParseInt(string(m[2]))
	}
	return res
}

func Problem03Fancy(in io.Reader, out io.Writer) {
	var gold int64
	var mu sync.Mutex
	var w sync.WaitGroup
	w.Add(5)
	scanner := bufio.NewScanner(in)

	// mulRegex := regexp.MustCompile("mul\\((\\d+),(\\d+)\\)")
	enabled := true
	data := make([]byte, 0, 64*1024)

	jobs := make(chan []byte, 20)
	worker := func(jobs <-chan []byte) {
		for j := range jobs {
			func(){
				mu.Lock()
				defer mu.Unlock()
				gold += addMuls(j)
			}()
		}
		func(){
			mu.Lock()
			defer mu.Unlock()
			w.Done()
		}()
	}

	for range 5 {
		go worker(jobs)
	}

	for scanner.Scan() {
		l := scanner.Bytes()
		j := 0
		for j < len(l) {
			if enabled {
				if idx := strings.Index(string(l[j:]), "don't()"); idx != -1 {
					enabled = false
					data = append(data, l[j:j+idx]...)
					jobs <- data
					data = make([]byte, 0, 64*1024)
					j += idx + len("don't()")
					continue
				}
				data = append(data, l[j:]...)
				j = len(l)
			} else {
				if idx := strings.Index(string(l[j:]), "do()"); idx != -1 {
					enabled = true
					j += idx + len("do()")
					continue
				}
				j = len(l)
			}
		}
	}

	if len(data) > 0 {
		jobs <- data
	}
	close(jobs)
	w.Wait()

	fmt.Fprintf(out, "gold: %v\n", gold)
}

func Problem03Regular(in io.Reader, out io.Writer) {
	var silver, gold int64
	scanner := bufio.NewScanner(in)

	mulRegex := regexp.MustCompile("mul\\((\\d+),(\\d+)\\)|do\\(\\)|don't\\(\\)")
	mulEnabled := true

	for scanner.Scan() {
		l := scanner.Text()
		for _, tokens := range mulRegex.FindAllStringSubmatchIndex(l, -1) {
			if strings.HasPrefix(l[tokens[0]:], "do()") {
				mulEnabled = true
			} else if strings.HasPrefix(l[tokens[0]:], "don't()") {
				mulEnabled = false
			} else {
				num := util.MustParseInt(l[tokens[2]:tokens[3]]) *
					util.MustParseInt(l[tokens[4]:tokens[5]])
				silver += num
				if mulEnabled {
					gold += num
				}
			}
		}
	}

	fmt.Fprintf(out, "silver: %v\n", silver)
	fmt.Fprintf(out, "gold: %v\n", gold)
}

func Problem03(in io.Reader, out io.Writer) {
	Problem03Fancy(in, out)
}
