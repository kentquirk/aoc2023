package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strings"
)

type node struct {
	name  string
	left  *node
	right *node
}

// do this in two passes so we can build the tree with pointers
// instead of relying on names for indexing
func build1(lines []string) *node {
	namepat := regexp.MustCompile(`\w+`)
	names := make(map[string]*node)
	for _, line := range lines {
		words := namepat.FindAllString(line, -1)
		name := words[0]
		node := &node{name: name}
		names[name] = node
	}
	for _, line := range lines {
		words := namepat.FindAllString(line, -1)
		name := words[0]
		left := words[1]
		right := words[2]
		n := names[name]
		n.left = names[left]
		n.right = names[right]
	}
	return names["AAA"]
}

func build2(lines []string) []*node {
	namepat := regexp.MustCompile(`\w+`)
	names := make(map[string]*node)
	var roots []*node
	for _, line := range lines {
		words := namepat.FindAllString(line, -1)
		name := words[0]
		node := &node{name: name}
		if name[2] == 'A' {
			roots = append(roots, node)
		}
		names[name] = node
	}
	for _, line := range lines {
		words := namepat.FindAllString(line, -1)
		name := words[0]
		left := words[1]
		right := words[2]
		n := names[name]
		n.left = names[left]
		n.right = names[right]
	}
	return roots
}

func part1(lines []string) int {
	sequence := lines[0]
	node := build1(lines[2:])
	steps := 0
	for ; node.name != "ZZZ"; steps++ {
		dir := sequence[steps%len(sequence)]
		switch dir {
		case 'L':
			node = node.left
		case 'R':
			node = node.right
		default:
			panic("OOPS")
		}
	}
	return steps
}

func leastCommonMultiple(a []int) int {
	lcm := a[0]
	for _, n := range a[1:] {
		lcm = lcm * n / gcd(lcm, n)
	}
	return lcm
}

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

// These sequences are periodic, so we can just find the period
// and then find the least common multiple of the periods.
// We do that by running them all long enough so that they get past the first
// period, then we find the difference between the most recent pair for
// each sequence.
func part2(lines []string) int {
	sequence := lines[0]
	roots := build2(lines[2:])
	pairs := make([]struct {
		a, b int
	}, len(roots))
	for steps := 0; steps < 100000; steps++ {
		endcount := 0
		for i, n := range roots {
			dir := sequence[steps%len(sequence)]
			switch dir {
			case 'L':
				n = n.left
			case 'R':
				n = n.right
			default:
				panic("OOPS")
			}
			if n.name[2] == 'Z' {
				endcount++
				pairs[i].a = pairs[i].b
				pairs[i].b = steps
				// fmt.Printf("root %d: %d\n", i, steps)
			}
			roots[i] = n
		}
		if endcount == len(roots) {
			return steps
		}
	}
	diffs := make([]int, len(pairs))
	for i, p := range pairs {
		diffs[i] = p.b - p.a
	}
	fmt.Println(diffs)
	return leastCommonMultiple(diffs)
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
	fmt.Println(part2(lines))
}
