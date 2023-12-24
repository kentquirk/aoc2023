package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

type color string

type run struct {
	first int
	last  int
	color color
	fill  bool
}

type contigRun struct {
	first int
	last  int
	cross bool
}

type row struct {
	runs []*run
}

func newRow() *row {
	return &row{
		runs: []*run{},
	}
}

func (r *row) add(ru *run) {
	r.runs = append(r.runs, ru)
	sort.Slice(r.runs, func(i, j int) bool {
		return r.runs[i].first < r.runs[j].first
	})
}

func (r *row) colorAt(c int) color {
	for _, ru := range r.runs {
		if ru.first <= c && c <= ru.last {
			return ru.color
		}
		if ru.first > c {
			return ""
		}
	}
	return ""
}

func (r *row) isTrench(c int) bool {
	for _, ru := range r.runs {
		if ru.first <= c && c <= ru.last {
			return !ru.fill
		}
		if ru.first > c {
			return false
		}
	}
	return false
}

func (r *row) isFilled(c int) bool {
	for _, ru := range r.runs {
		if ru.first <= c && c <= ru.last {
			return ru.color != ""
		}
		if ru.first > c {
			return false
		}
	}
	return false
}

// returns a list of contiguous runs, ignoring color.
// so "#..#####..#" returns (0,0), (3, 8), (11,11)
// each run is also tagged with whether it's a cross or not.
func (l *lagoon) contiguous(rix int) []contigRun {
	contigs := []contigRun{}
	cur := contigRun{first: l.rows[rix].runs[0].first}
	for i := 1; i < len(l.rows[rix].runs); {
		if l.rows[rix].runs[i-1].last+1 == l.rows[rix].runs[i].first {
			i++
		} else {
			cur.last = l.rows[rix].runs[i-1].last
			// before we add it to contigs, we need to check the previous row
			// at the edges of the run to see if if this run is a cross or a cap.
			// it's a cross if it's only 1 wide
			if cur.first == cur.last {
				cur.cross = true
			}
			// it's a cross if the previous row has one edge up and one edge down
			if rix > 0 && (l.rows[rix-1].isTrench(cur.first) != l.rows[rix-1].isTrench(cur.last)) {
				cur.cross = true
			}
			contigs = append(contigs, cur)
			cur = contigRun{first: l.rows[rix].runs[i].first}
			i++
		}
	}
	cur.last = l.rows[rix].runs[len(l.rows[rix].runs)-1].last
	contigs = append(contigs, cur)
	if rix > 18 && rix < 30 {
		fmt.Println(rix, contigs)
	}
	return contigs
}

type lagoon struct {
	height int
	rows   []*row
}

func newLagoon() *lagoon {
	rows := []*row{{runs: []*run{}}}

	return &lagoon{
		height: 1,
		rows:   rows,
	}
}

func (l *lagoon) left() int {
	left := 0
	for _, row := range l.rows {
		if len(row.runs) == 0 {
			continue
		}

		if left > row.runs[0].first {
			left = row.runs[0].first
		}
	}
	return left
}

func (l *lagoon) right() int {
	right := 0
	for _, row := range l.rows {
		if len(row.runs) == 0 {
			continue
		}

		if right < row.runs[len(row.runs)-1].last {
			right = row.runs[len(row.runs)-1].last
		}
	}
	return right
}

func (l *lagoon) count() int {
	c := 0
	for _, row := range l.rows {
		for _, ru := range row.runs {
			c += ru.last - ru.first + 1
		}
		// fmt.Println(c)
	}
	return c
}

func (l *lagoon) fill() {
	for ix, row := range l.rows {
		cont := l.contiguous(ix)
		inside := cont[0].cross
		for i := 1; i < len(cont); i++ {
			if inside {
				row.add(&run{
					first: cont[i-1].last + 1,
					last:  cont[i].first - 1,
					color: "#ffffff",
					fill:  true,
				})
			}
			if cont[i].cross {
				inside = !inside
			}
		}
	}
}

func (l *lagoon) String() string {
	var sb strings.Builder
	lt := l.left()
	rt := l.right()
	for _, row := range l.rows {
		for i := lt; i <= rt; i++ {
			if !row.isFilled(i) {
				sb.WriteString(".")
			} else {
				sb.WriteString("#")
			}
		}
		sb.WriteString("\n")
	}
	return sb.String()
}

type Instruction struct {
	operation string
	argument  int
	color     color
}

func parseInstruction(line string) Instruction {
	parts := strings.Split(line, " ")
	i, _ := strconv.Atoi(parts[1])
	return Instruction{
		operation: parts[0],
		argument:  i,
		color:     color(parts[2][1 : len(parts[2])-1]),
	}
}

func parseInstructions(lines []string) []Instruction {
	var instructions []Instruction
	for _, line := range lines {
		instructions = append(instructions, parseInstruction(line))
	}
	return instructions
}

// expands rows and returns the new row index
func (l *lagoon) maybeExpand(r, dr int) int {
	if r+dr < 0 {
		for r+dr < 0 {
			l.rows = append([]*row{newRow()}, l.rows...)
			l.height++
			r++
		}
	} else {
		for l.height < r+dr+1 {
			l.rows = append(l.rows, newRow())
			l.height++
		}
	}
	return r
}

func (l *lagoon) dig(r, c int, instruction Instruction) (int, int) {
	switch instruction.operation {
	case "R":
		l.rows[r].add(&run{
			first: c + 1,
			last:  c + instruction.argument,
			color: instruction.color,
		})
		return r, c + instruction.argument
	case "L":
		l.rows[r].add(&run{
			first: c - instruction.argument,
			last:  c - 1,
			color: instruction.color,
		})
		return r, c - instruction.argument
	case "U":
		r = l.maybeExpand(r, -instruction.argument)
		for i := 1; i <= instruction.argument; i++ {
			l.rows[r-i].add(&run{
				first: c,
				last:  c,
				color: instruction.color,
			})
		}
		return r - instruction.argument, c
	case "D":
		r = l.maybeExpand(r, instruction.argument)
		for i := 1; i <= instruction.argument; i++ {
			l.rows[r+i].add(&run{
				first: c,
				last:  c,
				color: instruction.color,
			})
		}
		return r + instruction.argument, c
	}
	panic("unknown operation")
}

func part1(lines []string) int {
	lagoon := newLagoon()
	instructions := parseInstructions(lines)
	r, c := 0, 0
	for _, instruction := range instructions {
		r, c = lagoon.dig(r, c, instruction)
	}
	fmt.Println(lagoon)
	lagoon.fill()
	fmt.Println("-------")
	fmt.Println(lagoon)
	return lagoon.count()
}

func part2(lines []string) int {
	return 0
}

func main() {
	args := os.Args[1:]
	name := "input"
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
