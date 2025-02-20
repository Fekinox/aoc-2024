package main

import (
	"fmt"
	"io"
	"slices"
	"strings"

	"github.com/Fekinox/aoc-2024/util"
)

func Problem23TriangleMaxxing(in io.Reader, out io.Writer) {
	var silver int64
	var gold string

	type edgeset struct {
		vertices []int16
		bitset *util.BoolVector
	}

	graph := make([]edgeset, 26*26)
	triangles := make([]int64, 26*26)

	for _, l := range util.ReadLines(in) {
		tokens := strings.Split(l, "-")
		v1 := int16(tokens[0][0] - 'a')*26 + int16(tokens[0][1] - 'a')
		v2 := int16(tokens[1][0] - 'a')*26 + int16(tokens[1][1] - 'a')
		graph[v1].vertices = append(graph[v1].vertices, v2)
		graph[v2].vertices = append(graph[v2].vertices, v1)
		if graph[v1].bitset == nil {
			graph[v1].bitset = util.NewBoolVector(26*26)
		}
		if graph[v2].bitset == nil {
			graph[v2].bitset = util.NewBoolVector(26*26)
		}
		graph[v1].bitset.Set(int(v2))
		graph[v2].bitset.Set(int(v1))
	}

	vToString := func(i int16) string {
		return fmt.Sprintf("%c%c", byte(i/26)+'a', byte(i%26)+'a')
	}

	// silver
	for i := range 26*26 {
		for _, j := range graph[i].vertices {
			if i > int(j) {
				continue
			}
			for _, k := range graph[int(j)].vertices {
				if j > k {
					continue
				}
				if !graph[i].bitset.Get(int(k)) {
					continue
				}
				triangles[i]++
				triangles[j]++
				triangles[k]++
				if i/26+'a' == 't' || j/26+'a' == 't' || k/26+'a' == 't' {
					silver++
				}
			}
		}
	}

	// This exploits a few properties of the input:
	// - There is one max-sized clique
	// - The clique is comprised of a bunch of triangles
	// - All nodes have the same degree
	// In this setting, a node is in the max clique if it has the maximum number of
	// triangles incident on it.
	// If it's fewer than the maximum, it cannot be in the clique.
	var best []string
	var bestTriangleCount int64
	for i, t := range triangles {
		if t > bestTriangleCount {
			bestTriangleCount = t
			best = []string{vToString(int16(i))}
		} else if t == bestTriangleCount {
			best = append(best, vToString(int16(i)))
		}
	}

	slices.Sort(best)
	gold = strings.Join(best, ",")

	fmt.Fprintln(out, "silver:", silver)
	fmt.Fprintln(out, "gold:", gold)
}

func Problem23BronKerbosch(in io.Reader, out io.Writer) {
	var silver int64
	var gold string

	type edgeset struct {
		vertices []int16
		bitset *util.BoolVector
	}

	graph := make([]edgeset, 26*26)

	for _, l := range util.ReadLines(in) {
		tokens := strings.Split(l, "-")
		v1 := int16(tokens[0][0] - 'a')*26 + int16(tokens[0][1] - 'a')
		v2 := int16(tokens[1][0] - 'a')*26 + int16(tokens[1][1] - 'a')
		graph[v1].vertices = append(graph[v1].vertices, v2)
		graph[v2].vertices = append(graph[v2].vertices, v1)
		if graph[v1].bitset == nil {
			graph[v1].bitset = util.NewBoolVector(26*26)
		}
		if graph[v2].bitset == nil {
			graph[v2].bitset = util.NewBoolVector(26*26)
		}
		graph[v1].bitset.Set(int(v2))
		graph[v2].bitset.Set(int(v1))
	}

	vToString := func(i int16) string {
		return fmt.Sprintf("%c%c", byte(i/26)+'a', byte(i%26)+'a')
	}

	bitsetToStrings := func(bs *util.BoolVector) string {
		var strs []string
		for i := range 26*26 {
			if len(graph[i].vertices) == 0 {
				continue
			}
			if bs != nil && bs.Get(i) {
				strs = append(strs, vToString(int16(i)))
			}
		}
		return strings.Join(strs, ",")
	}

	// silver
	for i := range 26*26 {
		for _, j := range graph[i].vertices {
			if i > int(j) {
				continue
			}
			for _, k := range graph[int(j)].vertices {
				if j > k {
					continue
				}
				if !graph[i].bitset.Get(int(k)) {
					continue
				}
				if i/26+'a' == 't' || j/26+'a' == 't' || k/26+'a' == 't' {
					silver++
				}
			}
		}
	}

	type state struct {
		r *util.BoolVector
		p *util.BoolVector
		x *util.BoolVector
	}

	stack := []state{{
		r: util.NewBoolVector(26*26),
		p: util.NewBoolVector(26*26),
		x: util.NewBoolVector(26*26),
	}}
	for i := range 26*26 {
		if len(graph[i].vertices) != 0 {
			stack[0].p.Set(i)
		}
	}

	var bestClique *util.BoolVector
	var bestSize int

	for len(stack) > 0 {
		cur := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		rr, pp, xx := cur.r, cur.p, cur.x

		var count int
		u := -1
		for i := range 26*26 {
			if u == -1 && pp.Get(i) || xx.Get(i) {
				u = i
				break
			}
			if rr.Get(i) {
				count++
			}
		}
		if u == -1 {
			if count > bestSize {
				bestClique = rr.Clone()
				bestSize = count
			}
			continue
		}
		ubs := graph[u].bitset
		for i := range 26*26 {
			if !(pp.Get(i) && !ubs.Get(i)) {
				continue
			}
			rrr, ppp, xxx := rr.Clone(), pp.Clone(), xx.Clone()
			rrr.Set(i)
			ppp.Intersection(graph[i].bitset)
			xxx.Intersection(graph[i].bitset)
			stack = append(stack, state{
				r: rrr,
				p: ppp,
				x: xxx,
			})
			pp.Unset(i)
			xx.Set(i)
		}
	}

	st := bitsetToStrings(bestClique)
	if len(st) > len(gold) {
		gold = st
	}

	fmt.Fprintln(out, "silver:", silver)
	fmt.Fprintln(out, "gold:", gold)
}

func Problem23Numeric(in io.Reader, out io.Writer) {
	var silver int64
	var gold string

	graph := make([][]int16, 26*26)

	for _, l := range util.ReadLines(in) {
		tokens := strings.Split(l, "-")
		v1 := int16(tokens[0][0] - 'a')*26 + int16(tokens[0][1] - 'a')
		v2 := int16(tokens[1][0] - 'a')*26 + int16(tokens[1][1] - 'a')
		graph[v1] = append(graph[v1], v2)
		graph[v2] = append(graph[v2], v1)
	}

	var maxDegree int
	for _, es := range graph {
		maxDegree = max(maxDegree, len(es))
	}

	vToString := func(i int16) string {
		return fmt.Sprintf("%c%c", byte(i/26)+'a', byte(i%26)+'a')
	}

	cliques := [][]int16{nil}
	var i int
	for {
		var nextCliques [][]int16
		for _, cl := range cliques {
			var verts []int16
			if len(cl) == 0 {
				for v, es := range graph {
					if len(es) == maxDegree {
						verts = append(verts, int16(v))
					}
				}
			} else {
				verts = graph[cl[0]]
			}
			outer:
			for _, v := range verts {
				if len(cl) != 0 && v <= cl[len(cl)-1] {
					continue
				}
				for _, v2 := range cl {
					if !slices.Contains(graph[v2], v) {
						continue outer
					}
				}
				newClique := make([]int16, len(cl))
				copy(newClique, cl)
				newClique = append(newClique, v)
				nextCliques = append(nextCliques, newClique)
			}
		}

		if len(nextCliques) == 0 {
			break
		}
		cliques = nextCliques
		i++

		if i == 3 {
			for _, cl := range cliques {
				if cl[0]/26 + 'a' == 't' || cl[1]/26 + 'a' == 't' || cl[2]/26 + 'a' == 't' {
					silver++
				}
			}
		}
	}

	slices.Sort(cliques[0])
	var sb strings.Builder
	for i, c := range cliques[0] {
		if i == len(cliques[0])-1 {
			fmt.Fprintf(&sb, "%s", vToString(c))
		} else {
			fmt.Fprintf(&sb, "%s,", vToString(c))
		}
	}

	gold = sb.String()

	fmt.Fprintln(out, "silver:", silver)
	fmt.Fprintln(out, "gold:", gold)
}


func Problem23BFS(in io.Reader, out io.Writer) {
	var silver int64
	var gold string

	graph := make(map[string][]string)

	for _, l := range util.ReadLines(in) {
		tokens := strings.Split(l, "-")
		graph[tokens[0]] = append(graph[tokens[0]], tokens[1])
		graph[tokens[1]] = append(graph[tokens[1]], tokens[0])
	}

	var maxDegree int
	for _, es := range graph {
		maxDegree = max(maxDegree, len(es))
	}

	cliques := [][]string{nil}
	var i int
	for {
		var nextCliques [][]string
		for _, cl := range cliques {
			var verts []string
			if len(cl) == 0 {
				for v, es := range graph {
					if len(es) == maxDegree {
						verts = append(verts, v)
					}
				}
			} else {
				verts = graph[cl[0]]
			}
			outer:
			for _, v := range verts {
				if len(cl) != 0 && v <= cl[len(cl)-1] {
					continue
				}
				for _, v2 := range cl {
					if !slices.Contains(graph[v2], v) {
						continue outer
					}
				}
				newClique := make([]string, len(cl))
				copy(newClique, cl)
				newClique = append(newClique, v)
				nextCliques = append(nextCliques, newClique)
			}
		}

		if len(nextCliques) == 0 {
			break
		}
		cliques = nextCliques
		i++

		if i == 3 {
			for _, cl := range cliques {
				if cl[0][0] == 't' || cl[1][0] == 't' || cl[2][0] == 't' {
					silver++
				}
			}
		}
	}

	slices.Sort(cliques[0])

	gold = strings.Join(cliques[0], ",")

	fmt.Fprintln(out, "silver:", silver)
	fmt.Fprintln(out, "gold:", gold)
}

func Problem23(in io.Reader, out io.Writer) {
	Problem23TriangleMaxxing(in, out)
}
