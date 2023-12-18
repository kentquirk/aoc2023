package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/dgryski/go-wyhash"
)

type grouplist []int

func (g grouplist) String() string {
	s := ""
	for _, n := range g {
		s += fmt.Sprintf("%d ", n)
	}
	return s
}

func (g grouplist) MinLen() int {
	if len(g) == 0 {
		return 0
	}
	total := 0
	for _, n := range g {
		total += n
	}
	return total + len(g) - 1
}

type memofunc func(grouplist, int) int

func hash(groups grouplist, start int) uint64 {
	var seed uint64 = 253295235
	b := make([]byte, len(groups)+1)
	for i, g := range groups {
		b[i] = byte(g)
	}
	b[len(b)-1] = byte(start)
	return wyhash.Hash(b, seed)
}

func memoize(cache map[uint64]int, f memofunc) memofunc {
	return func(groups grouplist, start int) int {
		h := hash(groups, start)
		v, ok := cache[h]
		if ok {
			return v
		}
		v = f(groups, start)
		cache[h] = v
		return v
	}
}

type row struct {
	source       []byte
	possibleLocs map[int]struct{}
	groups       grouplist
	caf          memofunc
	caa          memofunc
	cafCache     map[uint64]int
	caaCache     map[uint64]int
}

func groups(b []byte) []int {
	g := []int{}
	lastb := byte(0)
	for i := 0; i < len(b); i++ {
		if b[i] == '#' {
			if lastb != '#' {
				g = append(g, 0)
			}
			g[len(g)-1]++
		}
		lastb = b[i]
	}
	return g
}

func (r *row) isValid(arrangement []byte) bool {
	// check that the arrangement matches the groups
	groups := groups(arrangement)
	if len(groups) != len(r.groups) {
		return false
	}
	for i, g := range groups {
		if g != r.groups[i] {
			return false
		}
	}
	return true
}

func NewRow(line string, count int) *row {
	splits := strings.Split(line, " ")
	src := splits[0]
	grps := splits[1]
	for i := 1; i < count; i++ {
		src += "?" + splits[0]
		grps += "," + splits[1]
	}
	r := &row{
		source:       []byte(src),
		groups:       make(grouplist, 0),
		possibleLocs: make(map[int]struct{}),
		cafCache:     make(map[uint64]int),
		caaCache:     make(map[uint64]int),
	}
	r.caf = memoize(r.cafCache, r.countArrangementsFrom)
	r.caa = memoize(r.caaCache, r.countAllArrangements)

	for i, b := range r.source {
		if b != '.' {
			r.possibleLocs[i] = struct{}{}
		}
	}

	for _, num := range strings.Split(grps, ",") {
		n, _ := strconv.Atoi(num)
		r.groups = append(r.groups, n)
	}
	return r
}

func (r *row) String() string {
	return fmt.Sprintf("%s %v (p%v)", string(r.source), r.groups, r.possibleLocs)
}

// return the number of valid arrangements for the groups, anchored at start
func (r *row) countArrangementsFrom(groups grouplist, start int) int {
	// fmt.Printf("countArrangements(%v, %d)\n", groups, start)
	if len(groups) == 0 {
		// check that the rest of the row is not required
		for i := start; i < len(r.source); i++ {
			if r.source[i] == '#' {
				return 0
			}
		}
		return 1
	}
	_, ok := r.possibleLocs[start]
	if !ok {
		return 0
	}
	if start+groups.MinLen() > len(r.source) {
		return 0
	}
	// have a potential start location
	for i := 0; i < groups[0]; i++ {
		// check that the next i locations are valid
		// fmt.Printf("checking %d: %s\n", start, string(r.source[start:start+i]))
		if r.source[start+i] == '.' {
			return 0
		}
	}
	// check that the next location is not required (if it exists)
	if start+groups[0] < len(r.source) && r.source[start+groups[0]] == '#' {
		return 0
	}
	// it's valid, so recurse with the next group
	// fmt.Println(string(r.source), groups, start)
	return r.caa(groups[1:], start+groups[0]+1)
}

func (r *row) countAllArrangements(groups grouplist, start int) int {
	if len(groups) == 0 {
		return r.caf(groups, start)
	}

	total := 0
	last := bytes.IndexByte(r.source[start:], '#')
	if last == -1 {
		last = len(r.source) - groups.MinLen()
	} else {
		last += start
	}
	for i := start; i <= last; i++ {
		total += r.caf(groups, i)
	}
	return total
}

func solve(lines []string, count int) int {
	total := 0
	for _, line := range lines {
		r := NewRow(line, count)
		// fmt.Println(r)
		arr := r.caa(r.groups, 0)
		// fmt.Printf("%s: %d\n", line, arr)
		total += arr
	}
	return total
}

func main() {
	args := os.Args[1:]
	name := "sample"
	if len(args) > 0 {
		name = args[0]
	}
	f, err := os.Open(fmt.Sprintf("./data/%s.txt", name))
	if err != nil {
		log.Fatal(err)
	}
	b, err := io.ReadAll(f)
	if err != nil {
		log.Fatal(err)
	}
	lines := strings.Split(string(b), "\n")
	fmt.Println(solve(lines, 1))
	fmt.Println(solve(lines, 5))
}
