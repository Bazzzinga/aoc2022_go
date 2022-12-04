package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

const inputFileName = "input"

func getPriority(r rune) int {
	p := r - 96

	if p < 0 {
		p += 58
	}

	return int(p)
}

func priorityToChar(p int) rune {
	if p < 27 {
		return rune(p + 96)
	}

	return rune(p + 38)
}

func getHalves(s string) (string, string) {
	l := len(s)
	p1 := s[:l/2]
	p2 := s[l/2:]

	return p1, p2
}

func searchSameSymbol(ss []string) rune {
	prim := make([]bool, 52)

	for i, _ := range prim {
		prim[i] = true
	}

	for _, s := range ss {
		sec := make([]bool, 52)

		for _, c := range s {
			p := getPriority(c) - 1

			if prim[p] {
				sec[p] = true
			}
		}

		prim = sec
	}

	for i, b := range prim {
		if b {
			c := priorityToChar(i + 1)
			fmt.Println(string(c))
			return c
		}
	}

	return rune(0)
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
		line1 := scanner.Text()
		scanner.Scan()
		line2 := scanner.Text()
		scanner.Scan()
		line3 := scanner.Text()

		fmt.Println(line1, line2, line3)

		res += getPriority(
			searchSameSymbol(
				[]string{
					line1,
					line2,
					line3,
				},
			),
		)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Println(res)
}
