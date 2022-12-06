package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

const inputFileName = "input"

const markerLength = 14

func checkSymbol(s string, m map[string]int, p string) bool {
	res := false

	m[s]++
	res = true

	if len(p) > 0 {
		m[p]--
		if m[p] == 0 {
			delete(m, p)
		}
	}

	return res
}

func main() {

	file, err := os.Open(inputFileName)

	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	var line string
	for scanner.Scan() {
		line = scanner.Text()
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	m := make(map[string]int)

	for i, c := range line {
		var p string
		if i > markerLength-1 {
			tmp := line[i-markerLength : i-markerLength+1]
			p = tmp
		}
		ok := checkSymbol(string(c), m, p)

		if ok && len(m) == markerLength {
			fmt.Println(i + 1)
			break
		}
	}
}
