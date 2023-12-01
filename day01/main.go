package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strings"
)

func part1(lines []string) int {
	pat := regexp.MustCompile(`[\d]`)

	total := 0
	for _, line := range lines {
		m := pat.FindAllString(line, -1)
		if m == nil {
			fmt.Println("no match - ", line)
			continue
		}
		n1 := m[0][0] - '0'
		n2 := m[len(m)-1][0] - '0'
		fmt.Println(n1, n2)
		total += int(n1*10 + n2)
	}

	return total
}

var digits []string = []string{"zero", "one", "two", "three", "four", "five", "six", "seven", "eight", "nine"}

func part2(lines []string) int {
	total := 0
	for _, line := range lines {
		firstix := 999999
		lastix := 0
		firstval := 0
		lastval := 0
		// fmt.Println("---", line)
		pat := regexp.MustCompile(`[\d]`)
		m := pat.FindAllStringIndex(line, -1)
		if m != nil {
			firstix = m[0][0]
			firstval = int(line[m[0][0]] - '0')
			lastix = m[len(m)-1][1]
			lastval = int(line[m[len(m)-1][0]] - '0')
		}
		for d, w := range digits {
			pat = regexp.MustCompile(w)
			m = pat.FindAllStringIndex(line, -1)
			// fmt.Println(d, w, m)
			if m != nil {
				ix1 := m[0][0]
				if ix1 < firstix {
					firstix = ix1
					firstval = d
				}
				ix2 := m[len(m)-1][1]
				if ix2 > lastix {
					lastix = ix2
					lastval = d
				}
			}
		}
		// fmt.Println(firstval, lastval, line)
		total += int(firstval*10 + lastval)
	}

	return total
}

func main() {
	args := os.Args[1:]
	name := "sample2"
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
	fmt.Println(part2(lines))
}
