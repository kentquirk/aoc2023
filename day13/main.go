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
	score  int
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
	s := fmt.Sprintf("%d x %d (%d)\n", b.width, b.height, b.score)
	for _, row := range b.data {
		s += bstr(row, b.width) + "\n"
	}
	return s
}

func (b block) clone() block {
	var nb block
	nb.width = b.width
	nb.height = b.height
	nb.data = append([]int(nil), b.data...)
	return nb
}

func (b block) column(c int) int {
	n := 0
	bit := 1 << (b.width - c - 1)
	for row := 0; row < b.height; row++ {
		if b.data[row]&bit != 0 {
			n |= (1 << uint(b.height-row-1))
		}
	}
	return n
}

func (b block) rotate() block {
	var nb block
	nb.width = b.height
	nb.height = b.width
	for c := 0; c < b.width; c++ {
		nb.data = append(nb.data, b.column(c))
	}
	return nb
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

func (b block) findHorizontalReflection(old int) (int, bool) {
	for r := 0; r < b.height-1; r++ {
		if b.checkEqualRows(r, r+1) {
			if r+1 == old {
				// it's the same one, try again
				continue
			}
			return r + 1, true
		}
	}
	return 0, false
}

func (b block) findVerticalReflection() (int, bool) {
	for c := 0; c < b.width-1; c++ {
		if b.checkEqualColumns(c, c+1) {
			return c + 1, true
		}
	}
	return 0, false
}

func part1(blocks []block) int {
	total := 0
	for i, block := range blocks {
		// fmt.Printf("Block %d:\n%s\n", i, block)
		if r, ok := block.findHorizontalReflection(-1); ok {
			// fmt.Println("Reflection found in block ", i, "at row", r)
			blocks[i].score = 100 * r
			total += blocks[i].score
			continue
		}
		if c, ok := block.findVerticalReflection(); ok {
			// fmt.Println("Reflection found in block ", i, "at col", c)
			blocks[i].score = c
			total += blocks[i].score
			continue
		}
		fmt.Println("No reflection found in block ", i)
	}

	return total
}

func (b block) findNewHorizontalReflection(old int) (int, bool) {
	// XOR all pairs of rows together and count bits; if there's 1 bit then we can flip it and
	// try to see if we have a valid reflection
	for r1 := 0; r1 < b.height-1; r1++ {
		for r2 := r1 + 1; r2 < b.height; r2++ {
			x := b.data[r1] ^ b.data[r2]
			// fast way to tell if a number is a power of 2 (has only 1 bit set)
			if (x != 0) && (x&(x-1) == 0) {
				// We have a pair of rows that differ by 1 bit; try flipping one of them
				// and see if we have a valid reflection
				nb := b.clone()
				nb.data[r1] ^= x
				if r, ok := nb.findHorizontalReflection(old); ok {
					// fmt.Printf("FHR\n%s\n", nb)
					return r, true
				}
			}
		}
	}
	return 0, false
}

func part2(blocks []block) int {
	total := 0
	for i, block := range blocks {
		fmt.Printf("Block %d:\n%s\n", i, block)

		if r, ok := block.findNewHorizontalReflection(block.score / 100); ok {
			fmt.Println("New reflection found in block ", i, "for score ", r*100)
			total += 100 * r
			continue
		}

		nb := block.rotate()
		if r, ok := nb.findNewHorizontalReflection(block.score); ok {
			fmt.Println("New reflection found in block ", i, "for score ", r)
			total += r
			continue
		}

		fmt.Println("No new reflection found in block ", i)
	}

	return total
}

func main() {
	args := os.Args[1:]
	name := "sample7"
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
	sections := strings.Split(string(b), "\n\n")
	blocks := make([]block, 0)
	for _, s := range sections {
		blocks = append(blocks, loadBlock(strings.Split(s, "\n")))
	}
	fmt.Printf("Part 1: %d\n", part1(blocks))
	fmt.Printf("Part 2: %d\n", part2(blocks))
}
