package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

const inputFileName = "input"

type Cave struct {
	Valves        []*Valve
	history       map[string]int
	maxOpenValves int
}

func NewCave() *Cave {
	return &Cave{
		Valves:  make([]*Valve, 0),
		history: make(map[string]int),
	}
}

type Valve struct {
	Name    string
	Rate    uint
	Tunnels string
}

func NewValve(name string, rate uint, tunnels string) Valve {
	return Valve{
		Name:    name,
		Rate:    rate,
		Tunnels: tunnels,
	}
}

func (c *Cave) getMaxOpenValves() int {
	if c.maxOpenValves > 0 {
		return c.maxOpenValves
	}

	counter := 0
	for _, v := range c.Valves {
		if v.Rate > 0 {
			counter++
		}
	}

	c.maxOpenValves = counter

	return counter
}

func (c *Cave) ParseLine(line string) {
	expr := regexp.MustCompile(`^Valve\s([A-Z]+)\shas\sflow\srate=([\d]+);\stunnels?\sleads?\sto\svalves?\s([A-Z\s,]+)$`)

	matches := expr.FindStringSubmatch(line)

	name := matches[1]
	rate, _ := strconv.Atoi(matches[2])
	tunnels := matches[3]

	valve := NewValve(name, uint(rate), tunnels)

	c.Valves = append(c.Valves, &valve)
}

func (c *Cave) FindBestRate() {
	withRate := make([]*Valve, 0)

	graph := c.floydWarshall()

	for _, v := range c.Valves {
		if v.Rate > 0 || v.Name == "AA" {
			withRate = append(withRate, v)
		}
	}

	bitField := make(map[*Valve]uint)
	for i, v := range withRate {
		bitField[v] = 1 << i
	}

	var start uint

	for _, v := range withRate {
		if v.Name == "AA" {
			start = bitField[v]
			break
		}
	}

	bitGraphSlice := make([]uint, 0xffff)
	for _, v1 := range withRate {
		for _, v2 := range withRate {
			bitGraphSlice[bitField[v1]|bitField[v2]] = graph[v1][v2]
		}
	}

	withRateSlice := make([][2]uint, len(withRate))
	for i, v := range withRate {
		withRateSlice[i] = [2]uint{bitField[v], v.Rate}
	}

	var dfs func(target, pressure, minute, on, node uint) uint
	dfs = func(target, pressure, minute, on, node uint) uint {
		max := pressure
		for _, w := range withRateSlice {
			if node == w[0] || w[0] == start || w[0]&on != 0 {
				continue
			}
			l := bitGraphSlice[node|w[0]] + 1
			if minute+l > target {
				continue
			}

			next := dfs(target, pressure+(target-minute-l)*w[1], minute+l, on|w[0], w[0])

			if next > max {
				max = next
			}
		}
		return max
	}

	part1 := dfs(30, 0, 0, 0, start)
	fmt.Println("Answer 1:", part1)

	var dfsPaths func(target, pressure, minute, on, node, path uint) [][2]uint
	dfsPaths = func(target, pressure, minute, on, node, path uint) [][2]uint {
		paths := [][2]uint{{pressure, path}}
		for _, w := range withRateSlice {
			if w[0] == node || w[0] == start || w[0]&on != 0 {
				continue
			}
			l := bitGraphSlice[node|w[0]] + 1
			if minute+l > target {
				continue
			}
			paths = append(paths, dfsPaths(target, pressure+(target-minute-l)*w[1], minute+l, on|w[0], w[0], path|w[0])...)
		}
		return paths
	}

	allpaths := dfsPaths(26, 0, 0, 0, start, 0)

	var trimpaths [][2]uint
	for _, p := range allpaths {
		if p[0] > part1/2 {
			trimpaths = append(trimpaths, p)
		}
	}

	var max uint = 0
	for idx := 0; idx < len(trimpaths); idx += 1 {
		for jdx := idx + 1; jdx < len(trimpaths); jdx += 1 {
			if trimpaths[idx][1]&trimpaths[jdx][1] != 0 {
				continue
			}
			if m := trimpaths[idx][0] + trimpaths[jdx][0]; m > max {
				max = m
			}
		}
	}

	fmt.Println("Answer 2:", max)
}

func (c *Cave) floydWarshall() map[*Valve]map[*Valve]uint {
	graph := make(map[*Valve]map[*Valve]uint)
	for _, v1 := range c.Valves {
		graph[v1] = make(map[*Valve]uint)
		for _, v2 := range c.Valves {
			if v1 == v2 {
				graph[v1][v2] = 0
			} else if strings.Contains(v1.Tunnels, v2.Name) {
				graph[v1][v2] = 1
			} else {
				graph[v1][v2] = 0xff
			}
		}
	}

	for _, k := range c.Valves {
		for _, i := range c.Valves {
			for _, j := range c.Valves {
				if graph[i][j] > graph[i][k]+graph[k][j] {
					graph[i][j] = graph[i][k] + graph[k][j]
				}
			}
		}
	}

	return graph
}

func main() {
	file, err := os.Open(inputFileName)

	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	cave := NewCave()

	for scanner.Scan() {
		line := scanner.Text()
		cave.ParseLine(line)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	cave.FindBestRate()
}
