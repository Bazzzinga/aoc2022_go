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

type Interval struct {
	Start int
	End   int
}

func IsNumberBetween(n, min, max int) bool {
	return n >= min && n <= max
}

func intervalsIntersect(first, second Interval) bool {
	return IsNumberBetween(int(first.Start), int(second.Start), int(second.End)) ||
		IsNumberBetween(int(first.End), int(second.Start), int(second.End)) ||
		IsNumberBetween(int(second.Start), int(first.Start), int(first.End)) ||
		IsNumberBetween(int(second.End), int(first.Start), int(first.End))
}
func checkLine(line string) bool {
	parts := strings.Split(line, ",")

	elf1 := parts[0]
	elf2 := parts[1]

	sections1 := strings.Split(elf1, "-")
	sections2 := strings.Split(elf2, "-")

	s1, _ := strconv.Atoi(sections1[0])
	e1, _ := strconv.Atoi(sections1[1])

	s2, _ := strconv.Atoi(sections2[0])
	e2, _ := strconv.Atoi(sections2[1])

	//6-6 , 4-6

	i1 := Interval{
		Start: s1,
		End:   e1,
	}

	i2 := Interval{
		Start: s2,
		End:   e2,
	}

	/*if s1 <= s2 && e1 >= e2 {
		fmt.Printf("%d-%d contains %d-%d\n", s1, e1, s2, e2)
		return true
	}

	if s2 <= s1 && e2 >= e1 {
		fmt.Printf("%d-%d contains %d-%d\n", s2, e2, s1, e1)
		return true
	}

	if*/

	return intervalsIntersect(i1, i2)
}

func main() {

	file, err := os.Open(inputFileName)

	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	res := 0

	for scanner.Scan() {
		line := scanner.Text()

		if checkLine(line) {
			res++
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Println(res)
}
