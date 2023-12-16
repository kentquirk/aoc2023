package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

// We're going to model this with binary.
// Each arrangement has a set of masks.

// make a mask with 0s where the .s are and 1s where the # and ? are
// generate all the possible arrangements within the given length by expanding 0s
// iterate through all possibilities and AND with the mask -- if it's unchanged, it's valid

type row struct {
	length       int
	source       string
	mask         int
	groups       []int
	possibleGaps []int
}

func NewRow(line string) *row {
	splits := strings.Split(line, " ")
	r := &row{}
	r.source = splits[0]
	r.length = len(splits[0])
	r.mask = 0
	r.groups = []int{}
	bit := 1 << uint(r.length-1)
	for _, c := range line {
		switch c {
		case '#', '?':
			r.mask |= bit
		}
		bit >>= 1
	}
	for _, num := range strings.Split(splits[1], ",") {
		n, _ := strconv.Atoi(num)
		r.groups = append(r.groups, n)
	}
	r.possibleGaps = []int{}
	if r.source[0] == '.' {
		r.possibleGaps = append(r.possibleGaps, 0)
	}
	for i := 0; i < len(r.groups); i++ {
		r.possibleGaps = append(r.possibleGaps, 1)
	}
	if r.source[len(r.source)-1] == '.' {
		r.possibleGaps = append(r.possibleGaps, 0)
	}
	return r
}

func (r *row) String() string {
	return fmt.Sprintf("%2d %s %v %v %s", r.length, bstr(r.mask, r.length), r.groups, r.possibleGaps, r.source)
}

func bstr(i, l int) string {
	s := ""
	b := 1 << uint(l-1)
	for b != 0 {
		if (b & i) == 0 {
			s += "."
		} else {
			s += "#"
		}
		b >>= 1
	}
	return s
}

type gapmap map[string][]int

func (g gapmap) add(gap []int) {
	g[fmt.Sprintf("%v", gap)] = gap
}

func (g gapmap) all() [][]int {
	all := [][]int{}
	for _, v := range g {
		all = append(all, v)
	}
	return all
}

func generateAllGaps(gaps []int, nextra int) [][]int {
	if nextra == 0 {
		return [][]int{gaps}
	}
	gm := gapmap{}
	for i := 0; i < len(gaps); i++ {
		g := make([]int, len(gaps))
		copy(g, gaps)
		g[i]++
		all := generateAllGaps(g, nextra-1)
		for _, gap := range all {
			gm.add(gap)
		}
	}
	return gm.all()
}

func (r *row) makeArrangement(gaps []int) int {
	// add an extra 0-length gap at the end so the counts match
	gaps = append(gaps, 0)
	arr := 0
	bit := 1 << uint(r.length-1)
	for i, g := range r.groups {
		for j := 0; j < g; j++ {
			arr |= bit
			bit >>= 1
		}
		for j := 0; j < gaps[i]; j++ {
			bit >>= 1
		}
	}
	return arr
}

func (r *row) arrangements() []int {
	// generate all the possible arrangements within the given length by expanding 0s
	// arrs := []int{}
	// the minimum len is the sum of the groups plus the number of gaps
	minLen := 0
	for _, g := range r.groups {
		minLen += g
	}
	minLen += len(r.possibleGaps)
	nextra := r.length - minLen
	allgaps := generateAllGaps(r.possibleGaps, nextra)
	arrs := []int{}
	for _, gap := range allgaps {
		arr := r.makeArrangement(gap)
		fmt.Println(bstr(arr, r.length))
		fmt.Println(bstr(r.mask, r.length))
		fmt.Println()
		if arr&r.mask == arr {
			arrs = append(arrs, arr)
		}
	}

	return arrs
}

func part1(lines []string) int {
	total := 0
	for _, line := range lines {
		r := NewRow(line)
		a := r.arrangements()
		fmt.Println(r)
		for _, arr := range a {
			fmt.Println("  ", bstr(arr, r.length))
		}
		total += len(a)
	}
	return total
}

func part2(lines []string) int {
	return 0
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
	fmt.Println(part1(lines))
}
