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
const sourceX = 500
const sourceY = 0

type Cave struct {
	Slice map[int]map[int]bool
	MaxY  int
	MinX  int
	MaxX  int
}

func NewCave() *Cave {
	return &Cave{
		Slice: make(map[int]map[int]bool),
		MinX:  999999999,
	}
}

func (c *Cave) FillTile(x, y int) {
	if !c.isFilled(x, y) {
		c.Slice[x][y] = true
	}
}

func (c *Cave) isFilled(x, y int) bool {
	_, ok := c.Slice[x]

	if !ok {
		c.Slice[x] = make(map[int]bool)
	}

	return c.Slice[x][y]
}

func (c *Cave) isFilled2(x, y int) bool {
	if y >= c.MaxY+2 {
		return true
	}

	return c.isFilled(x, y)
}

func (c *Cave) SpawnSand() bool {
	x := sourceX
	y := sourceY
	if c.isFilled(sourceX, sourceY) {
		return true
	}

	done := false

	for !done {
		tryY := y + 1

		if tryY > c.MaxY {
			return true
		}

		if !c.isFilled(x, tryY) {
			y = tryY
			continue
		}

		tryX := x - 1
		if tryX < c.MinX {
			return true
		}

		if !c.isFilled(tryX, tryY) {
			x = tryX
			y = tryY
			continue
		}

		tryX = x + 1
		if tryX > c.MaxX {
			return true
		}

		if !c.isFilled(tryX, tryY) {
			x = tryX
			y = tryY
			continue
		}

		done = true
	}

	c.FillTile(x, y)

	return false
}

func (c *Cave) SpawnSand2() bool {
	x := sourceX
	y := sourceY
	if c.isFilled(sourceX, sourceY) {
		return true
	}

	done := false

	for !done {
		tryY := y + 1

		if !c.isFilled2(x, tryY) {
			y = tryY
			continue
		}

		tryX := x - 1

		if !c.isFilled2(tryX, tryY) {
			x = tryX
			y = tryY
			continue
		}

		tryX = x + 1

		if !c.isFilled2(tryX, tryY) {
			x = tryX
			y = tryY
			continue
		}

		done = true
	}

	c.FillTile(x, y)

	return false
}

func (c *Cave) Simulate() {
	counter := 0
	for !c.SpawnSand2() {
		counter++
	}

	fmt.Println("Sand count:", counter)
}

func (c *Cave) AddRocks(line string) {
	corners := strings.Split(line, " -> ")
	first := true
	var px, py int

	for _, corner := range corners {
		coords := strings.Split(corner, ",")

		xs := coords[0]
		ys := coords[1]

		x, _ := strconv.Atoi(xs)
		y, _ := strconv.Atoi(ys)

		c.FillTile(x, y)

		if x < c.MinX {
			c.MinX = x
		}

		if x > c.MaxX {
			c.MaxX = x
		}

		if y > c.MaxY {
			c.MaxY = y
		}

		if !first {
			done := false
			if y == py {
				tx := x
				d := int(float64(px-x) / math.Abs(float64(px-x)))
				for !done {
					tx += d
					c.FillTile(tx, y)
					if tx == px {
						done = true
					}
				}
			} else {
				ty := y
				d := int(float64(py-y) / math.Abs(float64(py-y)))
				for !done {
					ty += d
					c.FillTile(x, ty)
					if ty == py {
						done = true
					}
				}
			}
		}

		px = x
		py = y
		first = false
	}
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

		cave.AddRocks(line)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	cave.Simulate()
}
