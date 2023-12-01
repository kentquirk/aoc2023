package main

import (
	"flag"
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
	args := flag.Args()
	name := "sample"
	if len(args) > 1 {
		name = args[1]
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