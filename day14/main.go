package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

type cell byte

const (
	empty cell = '.'
	cube  cell = '#'
	round cell = 'O'
)

type dish struct {
	width  int
	height int
	data   [][]cell
}

func parse(lines []string) *dish {
	d := &dish{
		width:  len(lines[0]),
		height: len(lines),
		data:   make([][]cell, len(lines)),
	}
	for i, line := range lines {
		d.data[i] = make([]cell, len(line))
		for j, c := range line {
			d.data[i][j] = cell(c)
		}
	}
	return d
}

func (d *dish) tiltNorth() {
	for {
		nmoved := 0
		for row := 1; row < d.height; row++ {
			for col := 0; col < d.width; col++ {
				if d.data[row][col] == round && d.data[row-1][col] == empty {
					d.data[row-1][col] = round
					d.data[row][col] = empty
					nmoved++
				}
			}
		}
		if nmoved == 0 {
			break
		}
	}
}

func (d *dish) calcLoad() int {
	load := 0
	for row := 0; row < d.height; row++ {
		for col := 0; col < d.width; col++ {
			if d.data[row][col] == round {
				load += d.height - row
			}
		}
	}
	return load
}

func (d *dish) String() string {
	var sb strings.Builder
	for _, line := range d.data {
		for _, c := range line {
			sb.WriteByte(byte(c))
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func part1(lines []string) int {
	d := parse(lines)
	fmt.Println(d)
	fmt.Println()
	d.tiltNorth()
	fmt.Println(d)
	return d.calcLoad()
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
