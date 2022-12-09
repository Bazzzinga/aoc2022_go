package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

const inputFileName = "input"

type Simulation struct {
	head       Coords
	tails      []Coords
	visited    map[int]map[string]struct{}
	tailsCount int
}

type Coords struct {
	X int
	Y int
}

func NewSimulation(tails int) *Simulation {
	m := make(map[int]map[string]struct{})

	for i := 0; i < tails; i++ {
		m[i] = make(map[string]struct{})
	}

	s := &Simulation{
		visited:    m,
		tailsCount: tails,
		tails:      make([]Coords, tails),
	}

	s.markTailsVisited()

	return s
}

func (s *Simulation) CountVisited() int {
	return len(s.visited[s.tailsCount-1])
}

func (s *Simulation) markTailsVisited() {
	for i := 0; i < s.tailsCount; i++ {
		s.markTailVisited(i)
	}
}

func (s *Simulation) markTailVisited(n int) {
	key := fmt.Sprintf("%d_%d", s.tails[n].X, s.tails[n].Y)
	s.visited[n][key] = struct{}{}
}

func (s *Simulation) ProcessLine(line string) {
	parts := strings.Split(line, " ")

	dir := parts[0]
	moves, _ := strconv.Atoi(parts[1])

	move := Coords{}

	switch dir {
	case "U":
		move.Y = 1
	case "D":
		move.Y = -1
	case "R":
		move.X = 1
	case "L":
		move.X = -1
	}

	for i := 0; i < moves; i++ {
		s.moveHead(move)
		s.moveTails()
		s.markTailsVisited()
	}
}

func (s *Simulation) moveTails() {
	for i := 0; i < s.tailsCount; i++ {
		s.moveTail(i)
	}
}

func (s *Simulation) moveTail(n int) {
	var hX int
	var hY int

	if n == 0 {
		hX = s.head.X
		hY = s.head.Y
	} else {
		hX = s.tails[n-1].X
		hY = s.tails[n-1].Y
	}

	dX := math.Abs(float64(hX - s.tails[n].X))
	dY := math.Abs(float64(hY - s.tails[n].Y))

	if hX == s.tails[n].X {
		if dY > 1 {
			d := (hY - s.tails[n].Y) / int(dY)
			s.tails[n].Y += d
			return
		}
	}

	if hY == s.tails[n].Y {
		if dX > 1 {
			d := (hX - s.tails[n].X) / int(dX)
			s.tails[n].X += d
			return
		}
	}

	if dX > 1 || dY > 1 {
		d1 := (hX - s.tails[n].X) / int(dX)
		d2 := (hY - s.tails[n].Y) / int(dY)
		s.tails[n].X += d1
		s.tails[n].Y += d2
	}
}

func (s *Simulation) moveHead(move Coords) {
	s.head.X += move.X
	s.head.Y += move.Y
}

func main() {
	file, err := os.Open(inputFileName)

	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	sim := NewSimulation(9)

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		sim.ProcessLine(line)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Println(sim.CountVisited())
}
