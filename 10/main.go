package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

const inputFileName = "input"

func IncludesInt(haystack []int, needle int) bool {
	for _, i := range haystack {
		if i == needle {
			return true
		}
	}
	return false
}

type Simulation struct {
	cycle       int
	x           int
	input       chan int
	strength    int
	checkpoints []int
	finished    chan struct{}
}

func NewSimulation() *Simulation {
	return &Simulation{
		cycle:       1,
		x:           1,
		input:       make(chan int),
		finished:    make(chan struct{}),
		checkpoints: []int{20, 60, 100, 140, 180, 220},
	}
}

func (s *Simulation) Run() {
	go func() {
		for cmd := range s.input {
			s.render()

			if IncludesInt(s.checkpoints, s.cycle) {
				d := s.x * s.cycle
				s.strength += d
			}

			s.x += cmd
			s.cycle++
		}

		s.finished <- struct{}{}
	}()
}

func (s *Simulation) render() {
	sprite := []int{
		s.x,
		s.x + 1,
		s.x + 2,
	}

	if IncludesInt(sprite, s.cycle%40) {
		fmt.Printf("#")
	} else {
		fmt.Printf(".")
	}

	if s.cycle%40 == 0 {
		fmt.Printf("\n")
	}
}

func (s *Simulation) ProcessCommand(cmd string) {
	if cmd == "noop" {
		s.input <- 0
		return
	}

	parts := strings.Split(cmd, " ")

	if parts[0] == "addx" {
		d, _ := strconv.Atoi(parts[1])
		s.input <- 0
		s.input <- d
	}
}

func (s *Simulation) EndOfInput() {
	close(s.input)
}

func (s *Simulation) GetStrength() int {
	return s.strength
}

func (s *Simulation) WaitFinished() {
	<-s.finished
}

func main() {
	file, err := os.Open(inputFileName)

	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()
	fmt.Printf("\n\n")

	scanner := bufio.NewScanner(file)

	sim := NewSimulation()
	sim.Run()

	for scanner.Scan() {
		line := scanner.Text()
		sim.ProcessCommand(line)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	sim.EndOfInput()
	sim.WaitFinished()

	fmt.Printf("\n\n")

	fmt.Printf("Strenght: %d\n", sim.GetStrength())
}
