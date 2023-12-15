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
	loads  []int
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
	d.loads = append(d.loads, d.calcLoad())
	return d
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

func (d *dish) tiltSouth() {
	for {
		nmoved := 0
		for row := d.height - 2; row >= 0; row-- {
			for col := 0; col < d.width; col++ {
				if d.data[row][col] == round && d.data[row+1][col] == empty {
					d.data[row+1][col] = round
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

func (d *dish) tiltWest() {
	for {
		nmoved := 0
		for col := 1; col < d.width; col++ {
			for row := 0; row < d.height; row++ {
				if d.data[row][col] == round && d.data[row][col-1] == empty {
					d.data[row][col-1] = round
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

func (d *dish) tiltEast() {
	for {
		nmoved := 0
		for col := d.width - 2; col >= 0; col-- {
			for row := 0; row < d.height; row++ {
				if d.data[row][col] == round && d.data[row][col+1] == empty {
					d.data[row][col+1] = round
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

func (d *dish) cycle() {
	d.tiltNorth()
	d.tiltWest()
	d.tiltSouth()
	d.tiltEast()
}

func (d *dish) getLoadForCycle(n int) int {
	if n < len(d.loads) {
		return d.loads[n]
	}
	for i := len(d.loads); i <= n; i++ {
		d.cycle()
		d.loads = append(d.loads, d.calcLoad())
	}
	return d.loads[n]
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

func part1(lines []string) int {
	d := parse(lines)
	fmt.Println(d)
	fmt.Println()
	d.tiltNorth()
	fmt.Println(d)
	return d.calcLoad()
}

// We need to find a cycle in the iterations of the dish.
// The normal cycle detection algorithms don't work here because
// the values are not unique and we need to make sure that the cycle
// is long enough.

// So we're going to grab a sequence of values and count the number of occurrences of each value within it.
// As long as all of them have a count more than 1, we probably are likely to contain a cycle within it,
// so we choose a value that has the least count, look for two occurrences of it, and then check if that's a cycle.

// If not, we grab a longer chunk and repeat.

func findCycle(d *dish) (int, int) {
	chunksize := 100
	start := 100
	minload := 0
	mincount := 0
	for {
		loads := make(map[int]int)
		for i := 0; i < chunksize; i++ {
			loads[d.getLoadForCycle(i+start)]++
		}
		for load, count := range loads {
			if mincount == 0 || count < mincount {
				minload = load
				mincount = count
			}
		}
		if mincount == 1 {
			start += chunksize
			chunksize *= 2
			mincount = 0
			// fmt.Println(loads)
			fmt.Println("increasing start to", start)
			fmt.Println("increasing chunk size to", chunksize)
			continue
		}
		fmt.Println(minload, mincount)
		break
	}
	first, second := -1, -1
	for i := start; i < start+chunksize; i++ {
		if d.getLoadForCycle(i) == minload {
			if first == -1 {
				first = i
			} else if second == -1 {
				second = i
				break
			}
		}
	}
	// now verify that this is a cycle
	for i := first; i < second; i++ {
		if d.getLoadForCycle(i) != d.getLoadForCycle(i+second-first) {
			fmt.Println("not a real cycle")
		}
	}
	return second - first, first
}

func part2(lines []string) int {
	d := parse(lines)
	cyclelen, startindex := findCycle(d)
	fmt.Println(cyclelen, startindex)
	// now calculate the load for the 1_000_000_000 cycle
	// we know that the cycle starts at startindex and has length cyclelen

	ix := startindex + (1_000_000_000-startindex)%cyclelen
	return d.getLoadForCycle(ix)
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
	// fmt.Println(part1(lines))
	fmt.Println(part2(lines))
}
