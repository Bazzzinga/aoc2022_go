package main

import (
	"bufio"
	"container/list"
	"fmt"
	"log"
	"os"
)

const inputFileName = "input"

type Cell struct {
	X          int
	Y          int
	Key        string
	Height     rune
	Neighbours []*Cell
}

func NewCell(height rune, x, y int) *Cell {
	ns := make([]*Cell, 0)
	return &Cell{
		X:          x,
		Y:          y,
		Key:        fmt.Sprintf("%d_%d", x, y),
		Height:     height,
		Neighbours: ns,
	}
}

func (c *Cell) AddNeighbour(n *Cell) {
	alreadyHas := false
	for _, nn := range c.Neighbours {
		if n.X == nn.X && n.Y == nn.Y {
			alreadyHas = true
		}
	}

	if !alreadyHas {
		c.Neighbours = append(c.Neighbours, n)
	}
}

type Map struct {
	lines       []string
	Start       *Cell
	Destination *Cell
	Cells       []*Cell
}

func CreateMap(lines []string) *Map {
	m := Map{
		lines: lines,
		Cells: make([]*Cell, 0),
	}

	startLabel := "S"
	startHeight := "a"

	var i, j int

FOUND:
	for jj, line := range lines {
		for ii, h := range line {
			if h == rune(startLabel[0]) {
				i = ii
				j = jj
				break FOUND
			}
		}
	}

	visited := make(map[string]*Cell)

	start := NewCell(rune(startHeight[0]), i, j)

	m.Start = start

	visited[start.Key] = start
	m.Cells = append(m.Cells, start)

	m.tryAddCell(start, i+1, j, &visited)
	m.tryAddCell(start, i-1, j, &visited)
	m.tryAddCell(start, i, j+1, &visited)
	m.tryAddCell(start, i, j-1, &visited)

	return &m
}

func (m *Map) tryAddCell(from *Cell, x, y int, visited *map[string]*Cell) {
	if x < 0 || y < 0 || y >= len(m.lines) || x >= len(m.lines[0]) {
		return
	}

	finishLabel := "E"
	finishHeight := "z"
	startLabel := "S"
	startHeight := "a"

	height := rune(m.lines[y][x])

	isFinish := false

	if height == rune(finishLabel[0]) {
		isFinish = true
		height = rune(finishHeight[0])
	}

	if height == rune(startLabel[0]) {
		height = rune(startHeight[0])
	}

	key := fmt.Sprintf("%d_%d", x, y)

	to, ok := (*visited)[key]

	if ok {
		if height <= from.Height || height-from.Height == 1 {
			from.AddNeighbour(to)
		}

		return
	}

	if from.X != x || from.Y != y {
		if height <= from.Height || height-from.Height == 1 {
			to := NewCell(height, x, y)

			m.Cells = append(m.Cells, to)

			from.AddNeighbour(to)

			if isFinish {
				m.Destination = to
			}

			(*visited)[key] = to

			m.tryAddCell(to, x+1, y, visited)
			m.tryAddCell(to, x-1, y, visited)
			m.tryAddCell(to, x, y+1, visited)
			m.tryAddCell(to, x, y-1, visited)
		}
	}
}

func (m *Map) ShortestPath(pred *map[string]*Cell, dist *map[string]int) bool {
	visited := make(map[string]*Cell)

	for _, c := range m.Cells {
		(*dist)[c.Key] = 4294967295
	}

	(*dist)[m.Start.Key] = 0

	queue := list.New()

	queue.PushBack(m.Start)
	visited[m.Start.Key] = m.Start

	for queue.Len() > 0 {
		i := queue.Front()
		c := i.Value.(*Cell)
		queue.Remove(i)
		for _, n := range c.Neighbours {
			if visited[n.Key] == nil {
				visited[n.Key] = n
				(*dist)[n.Key] = (*dist)[c.Key] + 1
				(*pred)[n.Key] = c
				queue.PushBack(n)

				if n.Key == m.Destination.Key {
					return true
				}
			}
		}
	}

	return false
}

func (m *Map) ShortestDistance() int {
	dist := make(map[string]int)
	pred := make(map[string]*Cell)

	if !m.ShortestPath(&pred, &dist) {
		return 0
	}

	path := make([]*Cell, 0)
	crawl := m.Destination

	path = append(path, crawl)

	for pred[crawl.Key] != nil {
		fmt.Println(crawl.Key)
		path = append(path, pred[crawl.Key])
		crawl = pred[crawl.Key]
	}

	return dist[m.Destination.Key]
}

func main() {
	file, err := os.Open(inputFileName)

	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	lines := make([]string, 0)

	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	m := CreateMap(lines)

	lowerLabel := "a"
	start := m.Start
	min := m.ShortestDistance()
	for _, c := range m.Cells {
		if c.Height == rune(lowerLabel[0]) && c.Key != start.Key {
			m.Start = c

			d := m.ShortestDistance()

			if d > 0 && d < min {
				min = d
			}
		}
	}

	fmt.Println(min)

}
