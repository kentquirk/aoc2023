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

type td struct {
	raceTime   int
	recordDist int
}

func (t td) DistanceForPress(b int) int {
	return (t.raceTime - b) * b
}

func (t td) D4PDerivative(b int) int {
	return -2*b + t.raceTime
}

func (t td) DistanceForPressAboveRecord(b int) int {
	return ((t.raceTime - b) * b) - t.recordDist
}

func (t td) BeatsRecord(b int) bool {
	return t.DistanceForPress(b) > t.recordDist
}

func parse1(lines []string) []td {
	numpat := regexp.MustCompile(`\d+`)
	ts := numpat.FindAllString(lines[0], -1)
	ds := numpat.FindAllString(lines[1], -1)
	times := make([]td, len(ts))
	for i := range ts {
		times[i].raceTime, _ = strconv.Atoi(ts[i])
		times[i].recordDist, _ = strconv.Atoi(ds[i])
	}
	return times
}

func parse2(lines []string) td {
	numpat := regexp.MustCompile(`\d+`)
	ts := numpat.FindAllString(lines[0], -1)
	ds := numpat.FindAllString(lines[1], -1)
	times, _ := strconv.Atoi(strings.Join(ts, ""))
	dists, _ := strconv.Atoi(strings.Join(ds, ""))
	return td{times, dists}
}

func part1(races []td) int {
	product := 1
	for _, race := range races {
		count := 0
		for t := 1; t < race.raceTime; t++ {
			if race.BeatsRecord(t) {
				count++
			}
		}
		product *= count
	}
	return product
}

func newtonsMethod(race td, startingGuess int) int {
	// use newton's method to find the zero of the DistanceForPress function
	// f(x) = (t-x)*x - d
	// f'(x) = -2x + t
	// x_n+1 = x_n - f(x_n)/f'(x_n)

	attempts := 0
	lastguess := -1
	guess := startingGuess
	for {
		attempts++
		if attempts > 1000 {
			fmt.Println("too many attempts")
			break
		}
		if guess == lastguess {
			if !race.BeatsRecord(guess) {
				if race.BeatsRecord(guess + 1) {
					return guess + 1
				}
				if race.BeatsRecord(guess - 1) {
					return guess - 1
				}
				fmt.Println("failure!")
			}
			return guess
		}
		fg := race.DistanceForPressAboveRecord(guess)
		if fg == 0 {
			break
		}
		dfg := race.D4PDerivative(guess)
		if dfg == 0 {
			guess++
			continue
		}
		lastguess = guess
		guess = guess - int(float64(fg)/float64(dfg))
	}
	return guess
}

func part2(race td) int {
	fmt.Println(race)
	min := newtonsMethod(race, 10)
	max := newtonsMethod(race, race.raceTime)
	fmt.Println(min, max)
	return max - min + 1
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
	fmt.Println(part1(parse1(lines)))
	fmt.Println(part2(parse2(lines)))
}
