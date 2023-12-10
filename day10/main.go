package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

type cell struct {
	visited bool
	pipe    rune
}

type pipes [][]cell

type coord struct {
	row, col int
}

type edge struct {
	startpipe rune
	endpipe   rune
	start     int
	end       int
}

type edgelist map[int][]edge

func (e edgelist) add(r int, edge edge) {
	if edge.valid() {
		e[r] = append(e[r], edge)
	}
}

func (e edgelist) inside(r, c int) bool {
	es := e[r]
	in := false
	// don't have to worry about start because we won't call this
	// for visited cells
	for _, edge := range es {
		if c > edge.end {
			in = !in
		}
	}
	return in
}

func (e edge) valid() bool {
	switch e.startpipe {
	case '|':
		return true
	case 'F':
		return e.endpipe == 'J'
	case 'L':
		return e.endpipe == '7'
	default:
		panic("bad start pipe")
	}
}

func load(lines []string) (pipes, coord) {
	p := make(pipes, len(lines))
	var sloc coord
	for i, line := range lines {
		p[i] = make([]cell, len(line))
		for j, c := range line {
			p[i][j] = cell{pipe: c}
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
		cell := p[loc.row][loc.col]
		cell.visited = true
		p[loc.row][loc.col] = cell
		switch cell.pipe {
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
			return total + 1
		default:
			panic(fmt.Sprintf("bad pipe %v %v %v", loc, dir, p[loc.row+dir.row][loc.col+dir.col]))
		}
	}
}

func part1(lines []string) (int, int) {
	const (
		L = 1 << iota
		R
		U
		D
	)
	sm := map[int]rune{
		U + L: 'J',
		U + R: 'L',
		D + L: '7',
		D + R: 'F',
		U + D: '|',
		L + R: '-',
	}

	var schar int
	p, sloc := load(lines)
	lens := make([]int, 0)
	if sloc.row < len(p)-1 {
		switch p[sloc.row+1][sloc.col].pipe {
		case '|', 'J', 'L':
			lens = append(lens, p.follow(coord{sloc.row + 1, sloc.col}, coord{1, 0}, 0))
			schar += D
		}
	}
	if sloc.col < len(p[0])-1 {
		switch p[sloc.row][sloc.col+1].pipe {
		case '-', 'J', '7':
			lens = append(lens, p.follow(coord{sloc.row, sloc.col + 1}, coord{0, 1}, 0))
			schar += R
		}
	}
	if sloc.row > 0 {
		switch p[sloc.row-1][sloc.col].pipe {
		case '|', 'F', '7':
			lens = append(lens, p.follow(coord{sloc.row - 1, sloc.col}, coord{-1, 0}, 0))
			schar += U
		}
	}
	if sloc.col > 0 {
		switch p[sloc.row][sloc.col-1].pipe {
		case '-', 'F', 'L':
			lens = append(lens, p.follow(coord{sloc.row, sloc.col - 1}, coord{0, -1}, 0))
			schar += L
		}
	}
	// just verify that we get the same distance around in each direction
	fmt.Printf("S=%c\n", sm[schar])
	for _, l := range lens {
		fmt.Println(l)
	}

	p[sloc.row][sloc.col].pipe = sm[schar]
	contained := p.Contained(true)

	return (lens[0] + 1) / 2, contained
}

func printBox(c rune) {
	m := map[rune]rune{
		'|':  0x2502,
		'-':  0x2500,
		'L':  0x2514,
		'F':  0x250C,
		'J':  0x2518,
		'7':  0x2510,
		'S':  0x2573,
		'I':  0x2598,
		'o':  0x2591,
		'\n': '\n',
	}
	fmt.Printf("%c", m[c])
}

func (p pipes) Contained(doPrint bool) int {
	print := func(r rune) {
		if doPrint {
			printBox(r)
		}
	}
	edges := make(edgelist)
	count := 0
	// build a list of edges
	for r, row := range p {
		e := edge{}
		for c, cell := range row {
			if !cell.visited {
				continue
			}
			switch cell.pipe {
			case '|':
				edges.add(r, edge{
					startpipe: cell.pipe,
					endpipe:   cell.pipe,
					start:     c,
					end:       c,
				})
			case 'F', 'L':
				e.startpipe = cell.pipe
				e.start = c
			case 'J', '7':
				if e.startpipe != 0 {
					e.endpipe = cell.pipe
					e.end = c
					edges.add(r, e)
					e = edge{}
				}
			}
		}
	}

	// now that we have the edges we can count inside
	for r, row := range p {
		for c, cell := range row {
			if !cell.visited {
				if edges.inside(r, c) {
					count++
					print('I')
				} else {
					print('o')
				}
			} else {
				print(cell.pipe)
			}
		}
		print('\n')
	}
	return count
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
	n, cont := part1(lines)
	fmt.Println(n)
	fmt.Println(cont)
}
