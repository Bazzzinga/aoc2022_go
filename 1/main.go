package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
)

const inputFileName = "input"

func main() {
	file, err := os.Open(inputFileName)

	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	max := make([]int, 3)
	current := 0

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()

		cal, err := strconv.Atoi(line)

		if err != nil {
			//save max
			if current > max[0] {
				max[2] = max[1]
				max[1] = max[0]
				max[0] = current
			} else if current > max[1] {
				max[2] = max[1]
				max[1] = current
			} else if current > max[2] {
				max[2] = current
			}
			current = 0
		} else {
			current += cal
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Println(max[0] + max[1] + max[2])
}
