package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

type coord struct {
	r int
	c int
}

type direction struct {
	dr int
	dc int
}

type tilekind byte

type tile struct {
	kind    tilekind
	visited map[direction]struct{}
}

func newTile(kind tilekind) *tile {
	return &tile{
		kind:    kind,
		visited: make(map[direction]struct{}),
	}
}

func (t tile) String() string {
	return string(t.kind)
}

const (
	empty   tilekind = '.'
	fmirror tilekind = '/'
	bmirror tilekind = '\\'
	hsplit  tilekind = '-'
	vsplit  tilekind = '|'
)

type grid struct {
	width  int
	height int
	coords map[coord]*tile
}

func load(lines []string) grid {
	g := grid{
		width:  len(lines[0]),
		height: len(lines),
		coords: make(map[coord]*tile),
	}
	for r, line := range lines {
		for c, char := range line {
			if char == '.' {
				continue
			}
			g.coords[coord{r: r, c: c}] = newTile(tilekind(char))
		}
	}
	return g
}

func (g grid) follow(beam coord, dir direction) {
	// fmt.Println(beam, dir, g.coords[beam])
	beam.r += dir.dr
	beam.c += dir.dc
	if beam.r < 0 || beam.c < 0 {
		return
	}
	if beam.r >= g.height || beam.c >= g.width {
		return
	}
	if _, ok := g.coords[beam]; !ok {
		g.coords[beam] = newTile(empty)
	}
	if _, ok := g.coords[beam].visited[dir]; ok {
		return
	}
	g.coords[beam].visited[dir] = struct{}{}
	switch g.coords[beam].kind {
	case empty:
		g.follow(beam, dir)
	case fmirror:
		g.follow(beam, direction{dr: -dir.dc, dc: -dir.dr})
	case bmirror:
		g.follow(beam, direction{dr: dir.dc, dc: dir.dr})
	case hsplit:
		if dir.dc == 0 {
			g.follow(beam, direction{dc: 1})
			g.follow(beam, direction{dc: -1})
		} else {
			g.follow(beam, dir)
		}
	case vsplit:
		if dir.dr == 0 {
			g.follow(beam, direction{dr: 1})
			g.follow(beam, direction{dr: -1})
		} else {
			g.follow(beam, dir)
		}
	}
}

func (g grid) count() int {
	count := 0
	for _, tile := range g.coords {
		if len(tile.visited) > 0 {
			count++
		}
	}
	return count
}

func (g grid) Print(visited bool) {
	for r := 0; r < g.height; r++ {
		for c := 0; c < g.width; c++ {
			if tile, ok := g.coords[coord{r: r, c: c}]; ok {
				if visited {
					if len(tile.visited) == 0 {
						fmt.Print(".")
					} else {
						fmt.Print("#")
					}
				} else {
					fmt.Print(tile.kind)
				}
			} else {
				fmt.Print(".")
			}
		}
		fmt.Println()
	}
}

func part1(lines []string) int {
	g := load(lines)
	g.follow(coord{r: 0, c: -1}, direction{dc: 1})
	// g.Print(true)
	return g.count()
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
