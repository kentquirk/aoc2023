package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type number struct {
	value        int
	row          int
	firstcol     int
	lastcol      int
	isPartNumber bool
}

func (n *number) String() string {
	return fmt.Sprintf("{%d (%d, %d-%d) %t}", n.value, n.row, n.firstcol, n.lastcol, n.isPartNumber)
}

type special struct {
	row int
	col int
}

func (n *number) checkAdjacent(special *special) bool {
	if special.row >= n.row-1 && special.row <= n.row+1 {
		if special.col >= n.firstcol-1 && special.col <= n.lastcol {
			return true
		}
	}
	return false
}

func checkAdjacent(numbers []*number, specials []*special) {
	for _, s := range specials {
		for _, n := range numbers {
			if n.checkAdjacent(s) {
				n.isPartNumber = true
				continue
			}
		}
	}
}

func parse(lines []string) ([]*number, []*special) {
	var numbers []*number
	var specials []*special

	numpat := regexp.MustCompile(`\d+`)
	specialpat := regexp.MustCompile(`[^0-9.]`)
	for i, line := range lines {
		nums := numpat.FindAllStringIndex(line, -1)
		for _, num := range nums {
			v, _ := strconv.Atoi(line[num[0]:num[1]])
			numbers = append(numbers, &number{
				value:    v,
				row:      i,
				firstcol: num[0],
				lastcol:  num[1],
			})
		}
		sps := specialpat.FindAllStringIndex(line, -1)
		for _, sp := range sps {
			specials = append(specials, &special{
				row: i,
				col: sp[0],
			})
		}
	}
	return numbers, specials
}

func total(numbers []*number) int {
	total := 0
	for _, n := range numbers {
		if n.isPartNumber {
			total += n.value
		}
	}
	return total
}

func part1(lines []string) int {
	return 0
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
	numbers, specials := parse(lines)
	checkAdjacent(numbers, specials)
	fmt.Println(numbers)
	fmt.Println(total(numbers))
}
