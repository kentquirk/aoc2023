package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

type pipes [][]rune

type coord struct {
	row, col int
}

func load(lines []string) (pipes, coord) {
	p := make(pipes, len(lines))
	var sloc coord
	for i, line := range lines {
		p[i] = make([]rune, len(line))
		for j, c := range line {
			p[i][j] = c
			if c == 'S' {
				sloc = coord{i, j}
			}
		}
	}
	return p, sloc
}

// follow returns the number of steps to follow the path from start to finish
// start is the starting point
// loc is the current location being evaluated
// dir is the direction we are moving
func (p pipes) follow(loc, dir coord, total int) int {
	for {
		switch p[loc.row][loc.col] {
		case '|':
			if dir.row == -1 {
				return p.follow(coord{loc.row - 1, loc.col}, coord{-1, 0}, total+1)
			}
			return p.follow(coord{loc.row + 1, loc.col}, coord{1, 0}, total+1)
		case '-':
			if dir.col == -1 {
				return p.follow(coord{loc.row, loc.col - 1}, coord{0, -1}, total+1)
			}
			return p.follow(coord{loc.row, loc.col + 1}, coord{0, 1}, total+1)
		case 'F':
			if dir.col == -1 {
				return p.follow(coord{loc.row + 1, loc.col}, coord{1, 0}, total+1)
			}
			return p.follow(coord{loc.row, loc.col + 1}, coord{0, 1}, total+1)
		case 'L':
			if dir.col == -1 {
				return p.follow(coord{loc.row - 1, loc.col}, coord{-1, 0}, total+1)
			}
			return p.follow(coord{loc.row, loc.col + 1}, coord{0, 1}, total+1)
		case 'J':
			if dir.row == 1 {
				return p.follow(coord{loc.row, loc.col - 1}, coord{0, -1}, total+1)
			}
			return p.follow(coord{loc.row - 1, loc.col}, coord{-1, 0}, total+1)
		case '7':
			if dir.row == -1 {
				return p.follow(coord{loc.row, loc.col - 1}, coord{0, -1}, total+1)
			}
			return p.follow(coord{loc.row + 1, loc.col}, coord{1, 0}, total+1)
		case 'S':
			return total
		default:
			panic(fmt.Sprintf("bad pipe %v %v %v", loc, dir, p[loc.row+dir.row][loc.col+dir.col]))
		}
	}
}

func part1(lines []string) int {
	p, sloc := load(lines)
	lens := make([]int, 0)
	if sloc.row < len(p)-1 {
		switch p[sloc.row+1][sloc.col] {
		case '|', 'J', 'L':
			lens = append(lens, p.follow(coord{sloc.row + 1, sloc.col}, coord{1, 0}, 0))
		}
	}
	if sloc.col < len(p[0])-1 {
		switch p[sloc.row][sloc.col+1] {
		case '-', 'J', '7':
			lens = append(lens, p.follow(coord{sloc.row, sloc.col + 1}, coord{0, 1}, 0))
		}
	}
	if sloc.row > 0 {
		switch p[sloc.row-1][sloc.col] {
		case '|', 'F', '7':
			lens = append(lens, p.follow(coord{sloc.row - 1, sloc.col}, coord{-1, 0}, 0))
		}
	}
	if sloc.col > 0 {
		switch p[sloc.row][sloc.col-1] {
		case '-', 'F', 'L':
			lens = append(lens, p.follow(coord{sloc.row, sloc.col - 1}, coord{0, -1}, 0))
		}
	}
	// just verify that we get the same distance around in each direction
	for _, l := range lens {
		fmt.Println(l)
	}
	return (lens[0] + 1) / 2
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
