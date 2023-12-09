package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

type sequence []int
type triangle []sequence

func (s sequence) Deltas() (sequence, bool) {
	var deltas sequence
	allzeros := true
	for i := 0; i < len(s)-1; i++ {
		diff := s[i+1] - s[i]
		deltas = append(deltas, diff)
		if diff != 0 {
			allzeros = false
		}
	}
	return deltas, allzeros
}

func BuildSequence(line string) sequence {
	var s sequence
	for _, v := range strings.Split(line, " ") {
		var i int
		fmt.Sscanf(v, "%d", &i)
		s = append(s, i)
	}
	return s
}

func BuildTriangle(s sequence) triangle {
	var t triangle
	t = append(t, s)
	for next, allzeros := s.Deltas(); !allzeros; next, allzeros = next.Deltas() {
		t = append(t, next)
	}
	return t
}

func (t triangle) NextValue() int {
	row := len(t) - 1
	next := t[row][len(t[row])-1]
	for row--; row >= 0; row-- {
		next += t[row][len(t[row])-1]
	}
	return next
}

func (t triangle) PrevValue() int {
	row := len(t) - 1
	prev := t[row][0]
	for row--; row >= 0; row-- {
		prev = t[row][0] - prev
	}
	return prev
}

func (t triangle) String() string {
	var s strings.Builder
	for _, row := range t {
		for _, v := range row {
			fmt.Fprintf(&s, "%d ", v)
		}
		fmt.Fprintln(&s)
	}
	return s.String()
}

func part1(lines []string) int {
	total := 0
	for _, line := range lines {
		s := BuildSequence(line)
		t := BuildTriangle(s)
		n := t.NextValue()
		total += n
		// fmt.Println(t)
		// fmt.Println(n, total)
	}
	return total
}

func part2(lines []string) int {
	total := 0
	for _, line := range lines {
		s := BuildSequence(line)
		t := BuildTriangle(s)
		n := t.PrevValue()
		total += n
		// fmt.Println(t)
		// fmt.Println(n, total)
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
	fmt.Println(part1(lines))
	fmt.Println(part2(lines))
}
