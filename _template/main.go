package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func part1(lines []string) int {
	return 0
}

func part2(lines []string) int {
	return 0
}

func main() {
	args := os.Args
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
