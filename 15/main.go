package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"regexp"
	"strconv"
	"time"
)

const inputFileName = "input"

type Cave struct {
	MinX int
	MinY int
	MaxX int
	MaxY int
}

type Sensor struct {
	X int
	Y int
	D int
}

type Point struct {
	X int
	Y int
}

func NewPoint(x, y int) *Point {
	return &Point{
		X: x,
		Y: y,
	}
}

func NewSensorAndBeacon(sx, sy, bx, by int) (*Sensor, *Point) {
	d := int(math.Abs(float64(sx-bx)) + math.Abs(float64(sy-by)))
	return &Sensor{
		X: sx,
		Y: sy,
		D: d,
	}, NewPoint(bx, by)
}

func (s *Sensor) InsideRange(x, y int) bool {
	d := int(math.Abs(float64(s.X-x)) + math.Abs(float64(s.Y-y)))

	return d <= s.D
}

func NewCave() *Cave {
	c := Cave{
		MinY: 999999999,
		MinX: 999999999,
		MaxX: -999999999,
		MaxY: -999999999,
	}

	return &c
}

func (c *Cave) AddBeacon(b *Point) {
	if c.MinX > b.X {
		c.MinX = b.X
	}

	if c.MaxX < b.X {
		c.MaxX = b.X
	}

	if c.MinY > b.Y {
		c.MinY = b.Y
	}

	if c.MaxY < b.Y {
		c.MaxY = b.Y
	}
}

func (c *Cave) AddSensor(s *Sensor) {
	if c.MinX > s.X-s.D {
		c.MinX = s.X - s.D
	}

	if c.MaxX < s.X+s.D {
		c.MaxX = s.X + s.D
	}

	if c.MinY > s.Y-s.D {
		c.MinY = s.Y - s.D
	}

	if c.MaxY < s.Y+s.D {
		c.MaxY = s.Y + s.D
	}
}

func firstPart(cave *Cave, beacons []*Point, sensors []*Sensor) {
	t := time.Now()
	y := 2000000

	counter := 0
	x := cave.MinX
	for x <= cave.MaxX {
		can := true
		isBeacon := false
		for _, b := range beacons {
			if b.X == x && b.Y == y {
				isBeacon = true
				break
			}
		}

		if !isBeacon {
			for _, s := range sensors {
				if !can {
					break
				}
				if s.InsideRange(x, y) {
					can = false
				}
			}
		}

		if !can {
			counter++
		}

		x++
	}
	dt := time.Now().Sub(t).Seconds()
	fmt.Println(counter)
	fmt.Printf("calculated in: %v\n", dt)
}

func tryPoint(x, y int, beacons []*Point, sensors []*Sensor) bool {
	if x < 0 || x > 4000000 || y < 0 || y > 4000000 {
		return false
	}

	can := true
	isBeacon := false
	for _, b := range beacons {
		if b.X == x && b.Y == y {
			isBeacon = true
			break
		}
	}

	if isBeacon {
		return false
	}

	if !isBeacon {
		for _, s := range sensors {
			if !can {
				break
			}
			if s.InsideRange(x, y) {
				can = false
			}
		}
	}

	return can
}

func secondPart(beacons []*Point, sensors []*Sensor) {
	t := time.Now()

	freq := func(x, y int) int {
		return x*4000000 + y
	}

STOP:
	for _, s := range sensors {
		d := s.D + 1

		dx := d
		dy := 0

		for dx >= 0 {
			if tryPoint(s.X+dx, s.Y+dy, beacons, sensors) {
				fmt.Printf("frequency: %v\n", freq(s.X+dx, s.Y+dy))
				break STOP
			}

			if tryPoint(s.X-dx, s.Y+dy, beacons, sensors) {
				fmt.Printf("frequency: %v\n", freq(s.X+dx, s.Y+dy))
				break STOP
			}

			if tryPoint(s.X+dx, s.Y-dy, beacons, sensors) {
				fmt.Printf("frequency: %v\n", freq(s.X+dx, s.Y+dy))
				break STOP
			}

			if tryPoint(s.X-dx, s.Y-dy, beacons, sensors) {
				fmt.Printf("frequency: %v\n", freq(s.X+dx, s.Y+dy))
				break STOP
			}

			dx--
			dy++
		}
	}

	dt := time.Now().Sub(t).Seconds()
	fmt.Printf("calculated in: %v\n", dt)
}

func main() {
	file, err := os.Open(inputFileName)

	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	expr := regexp.MustCompile(`^Sensor\sat\sx=([-\d]+),\sy=([-\d]+):\sclosest\sbeacon\sis\sat\sx=([-\d]+),\sy=([-\d]+)$`)

	cave := NewCave()

	sensors := make([]*Sensor, 0)
	beacons := make([]*Point, 0)

	for scanner.Scan() {
		line := scanner.Text()
		matches := expr.FindStringSubmatch(line)

		sx, _ := strconv.Atoi(matches[1])
		sy, _ := strconv.Atoi(matches[2])
		bx, _ := strconv.Atoi(matches[3])
		by, _ := strconv.Atoi(matches[4])

		sensor, beacon := NewSensorAndBeacon(sx, sy, bx, by)

		sensors = append(sensors, sensor)
		beacons = append(beacons, beacon)

		cave.AddBeacon(beacon)
		cave.AddSensor(sensor)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	firstPart(cave, beacons, sensors)
	secondPart(beacons, sensors)
}
