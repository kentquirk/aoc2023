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

type colorset struct {
	red   int
	green int
	blue  int
}

func (c colorset) String() string {
	return fmt.Sprintf("red: %d, green: %d, blue: %d", c.red, c.green, c.blue)
}

func (c colorset) isPossibleWith(bag colorset) bool {
	return c.red <= bag.red && c.green <= bag.green && c.blue <= bag.blue
}

func (c colorset) power() int {
	return c.red * c.green * c.blue
}

func (c *colorset) accumulate(draw colorset) {
	if draw.red > c.red {
		c.red = draw.red
	}
	if draw.green > c.green {
		c.green = draw.green
	}
	if draw.blue > c.blue {
		c.blue = draw.blue
	}
}

type game []colorset

func parse(lines []string) map[int]game {
	pat := regexp.MustCompile(`(\d+) (red|green|blue)`)
	games := make(map[int]game)
	for _, line := range lines {
		if line == "" {
			continue
		}
		sp := strings.Split(line, ":")
		gamenum, _ := strconv.Atoi(sp[0][5:])
		var g game
		for _, draws := range strings.Split(sp[1], ";") {
			ma := pat.FindAllStringSubmatch(draws, -1)
			var dr colorset
			for _, m := range ma {
				switch m[2] {
				case "red":
					dr.red, _ = strconv.Atoi(m[1])
				case "green":
					dr.green, _ = strconv.Atoi(m[1])
				case "blue":
					dr.blue, _ = strconv.Atoi(m[1])
				}
			}
			g = append(g, dr)
		}
		games[gamenum] = g
	}
	return games
}

func part1(games map[int]game) int {
	bag := colorset{red: 12, green: 13, blue: 14}
	total := 0
outer:
	for id, g := range games {
		for _, d := range g {
			if !d.isPossibleWith(bag) {
				continue outer
			}
		}
		total += id
	}
	return total
}

func part2(games map[int]game) int {
	total := 0

	for _, g := range games {
		bag := colorset{}
		for _, d := range g {
			bag.accumulate(d)
		}
		total += bag.power()
	}
	return total
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
	games := parse(lines)
	// fmt.Println(games)
	fmt.Println(part1(games))
	fmt.Println(part2(games))
}
