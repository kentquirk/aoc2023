package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

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

type block struct {
	width  int
	height int
	data   []int
}

func loadBlock(s []string) block {
	var b block
	b.height = len(s)
	for _, line := range s {
		b.width = len(line)
		n := 0
		for i := 0; i < b.width; i++ {
			c := line[i]
			if c == '#' {
				n |= (1 << uint(b.width-i-1))
			}
		}
		b.data = append(b.data, n)
	}
	return b
}

func (b block) String() string {
	s := ""
	for _, row := range b.data {
		s += bstr(row, b.width) + "\n"
	}
	return s
}

func (b block) column(c int) int {
	n := 0
	bit := 1 << (b.width - c - 1)
	for row := 0; row < b.height; row++ {
		if b.data[row]&bit != 0 {
			n |= (1 << uint(b.height-row-1))
		}
	}
	fmt.Println("Column", c, "is", bstr(n, b.height))
	return n
}

func (b block) checkEqualRows(r1, r2 int) bool {
	if r1 < 0 || r2 >= b.height {
		return true
	}
	if b.data[r1] == b.data[r2] {
		return b.checkEqualRows(r1-1, r2+1)
	}
	return false
}

func (b block) checkEqualColumns(c1, c2 int) bool {
	if c1 < 0 || c2 >= b.width {
		return true
	}
	if b.column(c1) == b.column(c2) {
		return b.checkEqualColumns(c1-1, c2+1)
	}
	return false
}

func part1(blocks []string) int {
	total := 0
outer:
	for i, b := range blocks {
		block := loadBlock(strings.Split(b, "\n"))
		fmt.Printf("Block %d:\n%s\n", i, block)
		for r := 0; r < block.height-1; r++ {
			if block.checkEqualRows(r, r+1) {
				fmt.Println("Reflection found in block ", i, "at row", r+1)
				total += 100 * (r + 1)
				continue outer
			}
		}
		for c := 0; c < block.width-1; c++ {
			if block.checkEqualColumns(c, c+1) {
				fmt.Println("Reflection found in block ", i, "at col", c+1)
				total += c + 1
				continue outer
			}
		}
		fmt.Println("No reflection found in block ", i)
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
	blocks := strings.Split(string(b), "\n\n")
	fmt.Printf("Part 1: %d\n", part1(blocks))
}
