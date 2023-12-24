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

type part struct {
	ratings map[string]int
}

func newPart(s string) *part {
	pat := regexp.MustCompile(`(\w+)=([0-9]+)`)
	p := &part{ratings: make(map[string]int)}
	for _, m := range pat.FindAllStringSubmatch(s, -1) {
		p.ratings[m[1]], _ = strconv.Atoi(m[2])
	}
	return p
}

func (p *part) rating() int {
	r := 0
	for _, v := range p.ratings {
		r += v
	}
	return r
}

type parts []*part

type rule struct {
	condition func(p *part) bool
	action    func(p *part) (string, bool)
}

func newRule(s string, accept, reject func()) rule {
	// "a<3:b" or "b"
	// group 1: key
	// group 2: comparator
	// group 3: value
	// group 4: destination
	pat := regexp.MustCompile(`(?:(\w+)([<>]?)([0-9]+):)?(\w+)`)
	m := pat.FindStringSubmatch(s)
	var cond func(p *part) bool
	var action func(p *part) (string, bool)

	if m[1] == "" {
		cond = func(p *part) bool {
			return true
		}
	} else {
		key := m[1]
		cv, _ := strconv.Atoi(m[3])
		switch m[2] {
		case "<":
			cond = func(p *part) bool {
				return p.ratings[key] < cv
			}
		case ">":
			cond = func(p *part) bool {
				return p.ratings[key] > cv
			}
		}
	}

	switch m[4] {
	case "A":
		action = func(p *part) (string, bool) {
			accept()
			return "A", true
		}
	case "R":
		action = func(p *part) (string, bool) {
			reject()
			return "R", true
		}
	default:
		action = func(p *part) (string, bool) {
			return m[4], false
		}
	}
	return rule{condition: cond, action: action}
}

type workflow struct {
	name  string
	rules []rule
}

func newWorkflow(s string, accept, reject func()) *workflow {
	pat := regexp.MustCompile(`(\w+){(.*)}`)
	m := pat.FindStringSubmatch(s)
	w := &workflow{name: m[1]}
	for _, l := range strings.Split(m[2], ",") {
		w.rules = append(w.rules, newRule(l, accept, reject))
	}
	return w
}

func (w *workflow) run(p *part) (string, bool) {
	for _, r := range w.rules {
		if r.condition(p) {
			return r.action(p)
		}
	}
	// no rule matched - can't happen
	panic("no rule matched")
}

type workshop map[string]*workflow

func (w workshop) run(p *part) bool {
	currentWorkflow := w["in"]
	for {
		next, done := currentWorkflow.run(p)
		if done {
			return next == "A"
		}
		currentWorkflow = w[next]
	}
}

func parse(lines []string, afunc, rfunc func()) (workshop, parts) {
	workflows := make(workshop)
	parts := make([]*part, 0)
	buildParts := false
	for _, l := range lines {
		if l == "" {
			buildParts = true
			continue
		}
		if buildParts {
			p := newPart(l)
			parts = append(parts, p)
		} else {
			w := newWorkflow(l, afunc, rfunc)
			workflows[w.name] = w
		}
	}
	return workflows, parts
}

func part1(lines []string) int {
	accepted := 0
	rejected := 0
	afunc := func() {
		accepted++
	}
	rfunc := func() {
		rejected++
	}
	workflows, parts := parse(lines, afunc, rfunc)
	totalRating := 0
	for _, p := range parts {
		if workflows.run(p) {
			totalRating += p.rating()
		}
	}

	return totalRating
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
	fmt.Println(part1(lines))
}
