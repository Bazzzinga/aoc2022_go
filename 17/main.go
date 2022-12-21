package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

const inputFileName = "input"

type Cave struct {
	Wind string

	Rocks map[int]map[int]struct{}
	Top   int
	Move  int
}

func (c *Cave) HasRock(x, y int) bool {
	if y < 0 || x < 0 || x > 6 {
		return true
	}

	r, ok := c.Rocks[y]

	if ok {
		_, ok2 := r[x]

		if ok2 {
			return true
		}
	}

	return false
}

func (c *Cave) AddRock(x, y int) {
	_, ok := c.Rocks[y]
	if !ok {
		c.Rocks[y] = make(map[int]struct{})
	}

	c.Rocks[y][x] = struct{}{}

	if y > c.Top {
		c.Top = y
	}
}

func (c *Cave) StopRock(r *Rock) {
	for y, row := range r.form {
		for _, x := range row {
			c.AddRock(x+r.X, y+r.Y)
		}
	}
}

func NewCave(wind string) *Cave {
	return &Cave{
		Wind:  wind,
		Rocks: make(map[int]map[int]struct{}),
		Top:   -1,
	}
}

type Rock struct {
	//x, y - bottom left corner
	X    int
	Y    int
	form map[int][]int
	t    int
}

func (r *Rock) Move(x, y int) {
	r.X += x
	r.Y += y
}

func (r *Rock) CollidesLeft(c *Cave) bool {
	for y, row := range r.form {
		for _, x := range row {
			if c.HasRock(r.X+x-1, r.Y+y) {
				return true
			}
		}
	}

	return false
}

func (r *Rock) CollidesRight(c *Cave) bool {
	for y, row := range r.form {
		for _, x := range row {
			if c.HasRock(r.X+x+1, r.Y+y) {
				return true
			}
		}
	}

	return false
}

func (r *Rock) CollidesBottom(c *Cave) bool {
	for y, row := range r.form {
		for _, x := range row {
			if c.HasRock(r.X+x, r.Y+y-1) {
				return true
			}
		}
	}

	return false
}

func (r *Rock) ToString() string {
	switch r.t {
	case 0:
		return "####"
	case 1:
		return ".#.\n###\n.#."
	case 2:
		return "..#\n..#\n###"
	case 3:
		return "#\n#\n#\n#"
	case 4:
		return "##\n##"
	}

	return ""
}

func NewRock(x, y, t int) *Rock {
	t = t % 5

	r := Rock{
		X:    x,
		Y:    y,
		form: make(map[int][]int),
		t:    t,
	}

	switch t {
	case 0:
		r.form[0] = []int{0, 1, 2, 3}
	case 1:
		r.form[0] = []int{1}
		r.form[1] = []int{0, 1, 2}
		r.form[2] = []int{1}
	case 2:
		r.form[0] = []int{0, 1, 2}
		r.form[1] = []int{2}
		r.form[2] = []int{2}
	case 3:
		r.form[0] = []int{0}
		r.form[1] = []int{0}
		r.form[2] = []int{0}
		r.form[3] = []int{0}
	case 4:
		r.form[0] = []int{0, 1}
		r.form[1] = []int{0, 1}
	}

	return &r
}

func (c *Cave) Simulate(rocks int) {
	windCount := len(c.Wind)

	move := 0
	rockNum := 0

	cache := map[[2]int][]int{}

	for rockNum < rocks {
		rock := NewRock(2, c.Top+4, rockNum)
		inAir := true

		key := [2]int{rock.t, move}

		height := c.Top + 1
		if ch, ok := cache[key]; ok {
			if n, d := rocks-rockNum, rockNum-ch[0]; n%d == 0 {
				fmt.Printf("part2: %v\n", height+n/d*(height-ch[1]))
				return
			}
		}

		cache[key] = []int{rockNum, height}

		for inAir {
			//62 = > (right)
			//60 = < (left)
			var dx int
			canMove := true
			if c.Wind[move] == 62 {
				dx = 1
				if rock.CollidesRight(c) {
					canMove = false
				}
			} else {
				dx = -1
				if rock.CollidesLeft(c) {
					canMove = false
				}
			}

			if canMove {
				rock.Move(dx, 0)
			}

			if rock.CollidesBottom(c) {
				inAir = false
			} else {
				rock.Move(0, -1)
			}

			move++
			move = move % windCount
		}

		c.StopRock(rock)

		rockNum++
	}
}

func main() {
	file, err := os.Open(inputFileName)

	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	var wind string

	for scanner.Scan() {
		wind = scanner.Text()
	}

	cave := NewCave(wind)

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	cave.Simulate(1000000000000)
	fmt.Println(cave.Top + 1)
}
