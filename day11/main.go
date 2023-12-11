package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

type loc struct {
	row int
	col int
}

type starmap struct {
	width      int
	height     int
	stars      []loc
	rowoffsets map[int]int
	coloffsets map[int]int
}

func load(lines []string) *starmap {
	m := &starmap{
		width:      0,
		height:     len(lines),
		stars:      make([]loc, 0),
		rowoffsets: make(map[int]int),
		coloffsets: make(map[int]int),
	}
	for i, line := range lines {
		m.width = len(line)
		for j, c := range line {
			if c == '#' {
				m.stars = append(m.stars, loc{row: i, col: j})
				m.rowoffsets[i] = 0
				m.coloffsets[j] = 0
			}
		}
	}

	offset := 0
	for r := 0; r < m.height; r++ {
		if _, ok := m.rowoffsets[r]; !ok {
			offset++
		} else {
			m.rowoffsets[r] = offset
		}
	}
	offset = 0
	for c := 0; c < m.width; c++ {
		if _, ok := m.coloffsets[c]; !ok {
			offset++
		} else {
			m.coloffsets[c] = offset
		}
	}
	return m
}

func (m *starmap) calcTotalDistance() int {
	abs := func(x int) int {
		if x < 0 {
			return -x
		}
		return x
	}

	totalDistance := 0
	for i := 0; i < len(m.stars)-1; i++ {
		for j := i + 1; j < len(m.stars); j++ {
			rowdist := abs((m.stars[i].row + m.rowoffsets[m.stars[i].row]) - (m.stars[j].row + m.rowoffsets[m.stars[j].row]))
			coldist := abs((m.stars[i].col + m.coloffsets[m.stars[i].col]) - (m.stars[j].col + m.coloffsets[m.stars[j].col]))
			dist := rowdist + coldist
			totalDistance += dist
		}
	}
	return totalDistance
}

func part1(lines []string) int {
	sm := load(lines)
	return sm.calcTotalDistance()
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
