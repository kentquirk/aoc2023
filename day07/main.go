package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"slices"
	"strconv"
	"strings"
)

type handType int

const (
	none handType = iota
	highCard
	onePair
	twoPair
	threeOfAKind
	fullHouse
	fourOfAKind
	fiveOfAKind
)

func (h handType) String() string {
	switch h {
	case none:
		return "none"
	case highCard:
		return "high card"
	case onePair:
		return "one pair"
	case twoPair:
		return "two pair"
	case threeOfAKind:
		return "three of a kind"
	case fullHouse:
		return "full house"
	case fourOfAKind:
		return "four of a kind"
	case fiveOfAKind:
		return "five of a kind"
	}
	return "unknown"
}

type hand struct {
	cards           string
	bid             int
	typ             handType
	comparableCards string
}

const ordering = "23456789TJQKA"
const jokerOrdering = "J23456789TQKA"

func NewHand(cards string, bid string, withJokers bool) *hand {
	b, _ := strconv.Atoi(bid)
	ordering := ordering
	if withJokers {
		ordering = jokerOrdering
	}
	h := &hand{
		cards: cards,
		bid:   b,
		comparableCards: strings.Map(func(r rune) rune {
			return rune(strings.IndexRune(ordering, r)) + 'a'
		}, cards),
	}
	h.SetType(withJokers)
	return h
}

func (h *hand) String() string {
	return fmt.Sprintf("[%s %s %4d %15s]", h.cards, h.comparableCards, h.bid, h.typ)
}

func (h *hand) SetType(withJokers bool) {
	m := make(map[rune]int)
	for _, c := range h.cards {
		m[c]++
	}
	paircount := 0
	jokercount := 0
	if withJokers {
		jokercount = m['J']
	}
	for _, v := range m {
		switch v {
		case 5:
			h.typ = fiveOfAKind
			return
		case 4:
			switch jokercount {
			case 0:
				h.typ = fourOfAKind
			case 1, 4:
				h.typ = fiveOfAKind
			}
			return
		case 3:
			switch jokercount {
			case 0:
				h.typ = threeOfAKind
			case 1:
				h.typ = fourOfAKind
				return
			case 2:
				h.typ = fiveOfAKind
				return
			case 3:
				if len(m) == 3 {
					h.typ = fourOfAKind
				} else {
					h.typ = fiveOfAKind
				}
				return
			}
		case 2:
			paircount++
		}
	}
	// we only get here if we still haven't determined the hand type
	switch paircount {
	case 2:
		switch jokercount {
		case 2:
			h.typ = fourOfAKind
		case 1:
			h.typ = fullHouse
		case 0:
			h.typ = twoPair
		}
		return
	case 1:
		switch jokercount {
		case 0:
			if h.typ == threeOfAKind {
				h.typ = fullHouse
				return
			}
		case 1, 2:
			// if there's one pair and one joker, it's 3 of a kind
			// if there's one pair and two jokers, the pair is the jokers so it's also 3 of a kind
			h.typ = threeOfAKind
			return
		case 3:
			h.typ = fiveOfAKind
			return
		}
		h.typ = onePair
		return
	case 0:
		if h.typ == none {
			if jokercount == 1 {
				h.typ = onePair
			} else {
				h.typ = highCard
			}
		}
		return
	}
}

func Compare(lhs, rhs *hand) int {
	if lhs.typ != rhs.typ {
		return int(rhs.typ) - int(lhs.typ)
	}
	return strings.Compare(rhs.comparableCards, lhs.comparableCards)
}

func eval(lines []string, withJokers bool) int {
	hands := make([]*hand, 0)
	for _, line := range lines {
		parts := strings.Split(line, " ")
		h := NewHand(parts[0], parts[1], withJokers)
		hands = append(hands, h)
	}
	slices.SortFunc(hands, Compare)
	for _, h := range hands {
		fmt.Println(h)
	}
	winnings := 0
	for i, h := range hands {
		winnings += h.bid * (len(hands) - i)
	}
	return winnings
}

func part1(lines []string) int {
	return eval(lines, false)
}

func part2(lines []string) int {
	return eval(lines, true)
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
