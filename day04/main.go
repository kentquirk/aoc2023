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

type Set[T comparable] map[T]struct{}

func (s Set[T]) Add(i T) {
	s[i] = struct{}{}
}

func (s Set[T]) AddAll(i ...T) {
	for _, v := range i {
		s.Add(v)
	}
}

func (s Set[T]) Contains(i T) bool {
	_, ok := s[i]
	return ok
}

func (s Set[T]) Len() int {
	return len(s)
}

func part1(lines []string) int {
	splitpat := regexp.MustCompile(`:|\|`)
	numpat := regexp.MustCompile(`\d+`)
	totalPoints := 0
	for _, line := range lines {
		winners := make(Set[string])
		numWinners := 0
		parts := splitpat.Split(line, -1)
		if len(parts) != 3 {
			continue
		}
		wins := numpat.FindAllString(parts[1], -1)
		winners.AddAll(wins...)
		for _, card := range numpat.FindAllString(parts[2], -1) {
			if winners.Contains(card) {
				numWinners++
			}
		}
		if numWinners > 0 {
			totalPoints += 1 << (numWinners - 1)
		}
	}
	return totalPoints
}

type card struct {
	id           int
	numWinners   int
	numInstances int
}

func part2(lines []string) int {
	splitpat := regexp.MustCompile(`:|\|`)
	numpat := regexp.MustCompile(`\d+`)
	cards := make(map[int]*card)
	lastid := 0
	for _, line := range lines {
		winners := make(Set[string])
		numWinners := 0
		parts := splitpat.Split(line, -1)
		if len(parts) != 3 {
			continue
		}
		id, _ := strconv.Atoi(numpat.FindString(parts[0]))
		wins := numpat.FindAllString(parts[1], -1)
		winners.AddAll(wins...)
		for _, card := range numpat.FindAllString(parts[2], -1) {
			if winners.Contains(card) {
				numWinners++
			}
		}
		cards[id] = &card{
			id:           id,
			numWinners:   numWinners,
			numInstances: 1,
		}
		lastid = id
	}
	totalPoints := 0
	for id := 1; id <= lastid; id++ {
		c := cards[id]
		if c.numWinners > 0 {
			for i := 1; i <= c.numWinners; i++ {
				cards[id+i].numInstances += c.numInstances
			}
		}
		totalPoints += c.numInstances
	}

	return totalPoints
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
	fmt.Println(part2(lines))
}
