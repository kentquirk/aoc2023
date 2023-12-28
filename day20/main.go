package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/dgryski/go-wyhash"
)

type pulse bool

const (
	high pulse = true
	low  pulse = false
	on   bool  = true
	off  bool  = false
)

func (p pulse) String() string {
	if p == high {
		return "high"
	}
	return "low"
}

type event struct {
	src  string
	pul  pulse
	dest string
}

func (e *event) String() string {
	return fmt.Sprintf("%s -%v> %s", e.src, e.pul, e.dest)
}

func (e *event) setDest(d string) *event {
	return &event{src: e.src, pul: e.pul, dest: d}
}

type module interface {
	// Receive a pulse and may return a pulse to send
	// or nil if none should be sent.
	Receive(e *event) *event
	Hash(seed uint64) uint64
}

type flipflop struct {
	name  string
	state bool
}

func newFlipFlop(name string) *flipflop {
	return &flipflop{name: name, state: off}
}

func (f *flipflop) String() string {
	return fmt.Sprintf("%s(flipflop)", f.name)
}

func (f *flipflop) Receive(e *event) *event {
	if e.pul == high {
		return nil
	}
	f.state = !f.state
	return &event{src: f.name, pul: pulse(f.state)}
}

func (f *flipflop) Hash(seed uint64) uint64 {
	b := []byte(f.name)
	if f.state {
		b = append(b, '1')
	} else {
		b = append(b, '0')
	}
	return wyhash.Hash(b, seed)
}

type conjunction struct {
	name   string
	inputs []string
	memory map[string]pulse
}

func newConjunction(name string) *conjunction {
	c := &conjunction{name: name, memory: make(map[string]pulse)}
	return c
}

func (c *conjunction) addInput(name string) {
	c.inputs = append(c.inputs, name)
	sort.Strings(c.inputs)
	c.memory[name] = low
}

func (c *conjunction) String() string {
	return fmt.Sprintf("%s(conjunction)", c.name)
}

func (c *conjunction) Receive(e *event) *event {
	c.memory[e.src] = e.pul
	for _, v := range c.memory {
		if v == low {
			return &event{src: c.name, pul: high}
		}
	}
	return &event{src: c.name, pul: low}
}

func (c *conjunction) Hash(seed uint64) uint64 {
	b := []byte(c.name)
	for _, inp := range c.inputs {
		b = append(b, inp...)
		if c.memory[inp] == low {
			b = append(b, '0')
		} else {
			b = append(b, '1')
		}
	}
	return wyhash.Hash(b, seed)
}

type broadcast struct {
	name string
}

func newBroadcast(name string) *broadcast {
	return &broadcast{name: name}
}

func (b *broadcast) String() string {
	return fmt.Sprintf("%s(broadcast)", b.name)
}

func (b *broadcast) Receive(e *event) *event {
	return &event{src: b.name, pul: e.pul}
}

func (b *broadcast) Hash(seed uint64) uint64 {
	return wyhash.Hash([]byte(b.name), seed)
}

type network struct {
	names   []string
	modules map[string]module
	dests   map[string][]string
	queue   chan *event
}

func newNetwork(lines []string) *network {
	conjunctions := make(map[string]*conjunction)
	n := &network{
		names:   make([]string, 0),
		modules: make(map[string]module),
		dests:   make(map[string][]string),
		queue:   make(chan *event, 100),
	}
	for _, l := range lines {
		parts := strings.Split(l, " -> ")
		src := parts[0]
		dests := strings.Split(parts[1], ", ")
		var m module
		name := ""
		switch src[0] {
		case '%':
			name = src[1:]
			m = newFlipFlop(name)
		case '&':
			name = src[1:]
			m = newConjunction(name)
			conjunctions[name] = m.(*conjunction)
		case 'b':
			name = src
			m = newBroadcast(name)
		}
		n.names = append(n.names, name)
		n.modules[name] = m
		n.dests[name] = dests
	}
	for src, dests := range n.dests {
		for _, dest := range dests {
			if c, ok := conjunctions[dest]; ok {
				c.addInput(src)
			}
		}
	}
	return n
}

func (n *network) String() string {
	s := ""
	for k, v := range n.modules {
		s += fmt.Sprintf("%s -> %v\n", v, n.dests[k])
	}
	return s
}

func (n *network) Hash() uint64 {
	h := uint64(1234)
	for _, name := range n.names {
		v := n.modules[name]
		h += v.Hash(h)
	}
	return h
}

var lastPulse = make(map[string]int)

func (n *network) processQueue(buttonPresses int) (int, int, bool) {
	counts := make(map[pulse]int)
	for len(n.queue) > 0 {
		cur := <-n.queue
		// fmt.Println(cur)
		counts[cur.pul]++
		// xmtoggles := make(map[string]int)
		if m, ok := n.modules[cur.dest]; ok {
			if evt := m.Receive(cur); evt != nil {
				for _, d := range n.dests[evt.src] {
					switch d {
					case "rx":
						if evt.pul == low {
							fmt.Println("got low to rx!")
							return counts[high], counts[low], true
						}
					case "ft", "jz", "ng", "sv":
						if evt.pul == low {
							delta := buttonPresses - lastPulse[d]
							lastPulse[d] = buttonPresses
							fmt.Printf("got %v from %s to %s at %d (%d)\n", evt.pul, evt.src, d, buttonPresses, delta)
						}
					}
					n.queue <- evt.setDest(d)
				}
			}
		}
	}
	return counts[high], counts[low], false
}

func (n *network) pressButton() {
	n.queue <- &event{src: "button", pul: low, dest: "broadcaster"}
}

type result struct {
	high     int
	low      int
	nexthash uint64
}

func part1(lines []string) int {
	net := newNetwork(lines)
	// fmt.Println(net)
	start := net.Hash()
	cycles := 0
	highTotal := 0
	lowTotal := 0
	h := start
	results := make(map[uint64]result)
	for i := 0; i < 1000; i++ {
		if r, ok := results[h]; ok {
			highTotal += r.high
			lowTotal += r.low
			h = r.nexthash
		} else {
			net.pressButton()
			high, low, _ := net.processQueue(0)
			highTotal += high
			lowTotal += low
			next := net.Hash()
			results[h] = result{high, low, next}
			h = next
			// fmt.Println(high, low, net.Hash())
		}
		cycles++
		if h == start {
			fmt.Println("cycle ", cycles, highTotal, lowTotal)
		}
	}
	fmt.Println(cycles, highTotal, lowTotal)
	return highTotal * lowTotal
}

func part2(lines []string) int {
	net := newNetwork(lines)
	// fmt.Println(net)
	start := net.Hash()
	highTotal := 0
	lowTotal := 0
	buttonPresses := 0
	for {
		buttonPresses++
		net.pressButton()
		high, low, done := net.processQueue(buttonPresses)
		highTotal += high
		lowTotal += low
		if done {
			return buttonPresses
		}
		if net.Hash() == start {
			fmt.Println("cycle! ", highTotal, lowTotal)
			return 0
		}
		if buttonPresses%1000000 == 0 {
			fmt.Printf("%d presses, %d pulses\n", buttonPresses, highTotal+lowTotal)
		}
	}
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
	fmt.Println(part2(lines))
}
