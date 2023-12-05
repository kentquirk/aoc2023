package main

import (
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type FarmMapRange struct {
	destStart   int
	sourceStart int
	count       int
}

type FarmMap struct {
	from   string
	to     string
	ranges []FarmMapRange
}

// for a given map, look up the value in the ranges
// and return the new value
func (m *FarmMap) Lookup(value int) (string, int) {
	for _, r := range m.ranges {
		if r.sourceStart <= value && value < r.sourceStart+r.count {
			offset := value - r.sourceStart
			return m.to, r.destStart + offset
		}
	}
	return m.to, value
}

type Table map[string]FarmMap

func (t Table) Convert(have string, want string, value int) int {
	for have != want {
		// fmt.Printf("%s(%d) -> ", have, value)
		if m, ok := t[have]; ok {
			have, value = m.Lookup(value)
		}
	}
	// fmt.Printf("END: %s(%d)\n", have, value)
	return value
}

func parse(data string) (Table, []int) {
	table := make(Table)
	var seeds []int
	blocks := strings.Split(data, "\n\n")
	numpat := regexp.MustCompile(`\d+`)
	for _, s := range numpat.FindAllString(blocks[0], -1) {
		n, _ := strconv.Atoi(s)
		seeds = append(seeds, n)
	}
	blocks = blocks[1:]

	for _, block := range blocks {
		lines := strings.Split(block, "\n")
		namepat := regexp.MustCompile(`\w+`)
		names := namepat.FindAllString(lines[0], -1)
		from, to := names[0], names[2]
		lines = lines[1:]
		m := FarmMap{from, to, []FarmMapRange{}}
		for _, line := range lines {
			var destStart, srcStart, count int
			fmt.Sscanf(line, "%d %d %d", &destStart, &srcStart, &count)
			m.ranges = append(m.ranges, FarmMapRange{
				destStart:   destStart,
				sourceStart: srcStart,
				count:       count,
			})
		}
		table[from] = m
	}
	return table, seeds
}

func part1(t Table, seeds []int) int {
	lowest := math.MaxInt64
	for _, seed := range seeds {
		v := t.Convert("seed", "location", seed)
		if v < lowest {
			lowest = v
		}
	}
	return lowest
}

// This is super slow. It could be solved by
// dividing the ranges into linear sub-parts and finding
// the minimum of those. But I'm also kind of done with
// this problem.
func part2(t Table, seeds []int) int {
	lowest := math.MaxInt64
	for i := 0; i < len(seeds); i += 2 {
		seed := seeds[i]
		for j := 0; j < seeds[i+1]; j++ {
			v := t.Convert("seed", "location", seed+j)
			if v < lowest {
				lowest = v
			}
		}
	}
	return lowest
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
	t, seeds := parse(string(b))
	fmt.Println(part1(t, seeds))
	fmt.Println(part2(t, seeds))
}
