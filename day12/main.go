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

type row struct {
	source []byte
	// length      int
	damagedLocs []int
	unknownLocs []int
	groups      []int

	// mask         int
	// possibleGaps []int
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

func NewRow(line string) *row {
	splits := strings.Split(line, " ")
	r := &row{}
	r.source = []byte(splits[0])
	for i, b := range r.source {
		if b == '?' {
			r.unknownLocs = append(r.unknownLocs, i)
		} else if b == '#' {
			r.damagedLocs = append(r.damagedLocs, i)
		}
	}

	r.groups = []int{}
	for _, num := range strings.Split(splits[1], ",") {
		n, _ := strconv.Atoi(num)
		r.groups = append(r.groups, n)
	}
	return r
}

func (r *row) String() string {
	return fmt.Sprintf("%s %v (u%v d%v)", string(r.source), r.groups, r.unknownLocs, r.damagedLocs)
}

func (r *row) countValidArrangements() int {
	// count the number of valid arrangements
	// start with the number of unknowns
	nvalid := 0
	unk := len(r.unknownLocs)
	for i := 0; i < (1 << unk); i++ {
		// now start with a clone of the source and fill in the unknowns based on i
		arr := make([]byte, len(r.source))
		copy(arr, r.source)

		bit := 1 << uint(unk-1)
		for ix := 0; ix < unk; ix++ {
			if i&bit != 0 {
				arr[r.unknownLocs[ix]] = '#'
			} else {
				arr[r.unknownLocs[ix]] = '.'
			}
			bit >>= 1
		}
		// fmt.Printf("arr %d: %s %t\n", i, string(arr), r.isValid(arr))
		if r.isValid(arr) {
			nvalid++
		}
	}
	return nvalid
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

func part1(lines []string) int {
	total := 0
	for _, line := range lines {
		r := NewRow(line)
		// fmt.Println(r)
		total += r.countValidArrangements()
	}
	return total
}

// func part2(lines []string) int {
// 	return 0
// }

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
