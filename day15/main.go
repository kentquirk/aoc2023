package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strings"
)

type slot struct {
	label string
	lens  int
}

type box struct {
	slots []slot
}

func (b *box) add(label string, lens int) {
	for i := range b.slots {
		if b.slots[i].label == label {
			b.slots[i].lens = lens
			return
		}
	}
	b.slots = append(b.slots, slot{label, lens})
}

func (b *box) remove(label string) {
	for i := range b.slots {
		if b.slots[i].label == label {
			b.slots = append(b.slots[:i], b.slots[i+1:]...)
			return
		}
	}
}

func (b *box) totalFocusingPower(boxID byte) int {
	total := 0
	for i, s := range b.slots {
		total += (int(boxID) + 1) * (i + 1) * s.lens
	}
	return total
}

func HASH(s string) byte {
	var h byte = 0
	for _, c := range s {
		h += byte(c)
		h *= 17
	}
	return h
}

func part1(steps []string) int {
	total := 0
	for _, s := range steps {
		total += int(HASH(s))
	}
	return total
}

func part2(steps []string) int {
	boxes := make(map[byte]*box)
	pat := regexp.MustCompile(`^(\w+)([=-])(\d)?$`)
	for _, s := range steps {
		parts := pat.FindStringSubmatch(s)
		label := parts[1]
		action := parts[2]
		bi := HASH(label)
		b, ok := boxes[bi]
		if !ok {
			b = &box{}
			boxes[bi] = b
		}
		switch action {
		case "=":
			lens := int(parts[3][0] - '0')
			b.add(label, lens)
		case "-":
			b.remove(label)
		default:
			log.Fatalf("unknown action %s", action)
		}
	}

	total := 0
	for bi, b := range boxes {
		// fmt.Printf("box %d: %v\n", bi, b)
		total += b.totalFocusingPower(bi)
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
	pat := regexp.MustCompile(`[,\n]+`)
	steps := pat.Split(strings.TrimSpace(string(b)), -1)
	fmt.Println(part1(steps))
	fmt.Println(part2(steps))
}
