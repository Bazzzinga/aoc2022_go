package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

const inputFileName = "input"

func main() {
	/*rules1 := map[string]map[string]int{
		"A": {
			"X": 4,
			"Y": 8,
			"Z": 3,
		},
		"B": {
			"X": 1,
			"Y": 5,
			"Z": 9,
		},
		"C": {
			"X": 7,
			"Y": 2,
			"Z": 6,
		},
	}*/

	// 1 for камень, 2 for бумага, and 3 for ножницы
	// 3 if the round was a draw, and 6 if you won
	rules := map[string]map[string]int{
		"A": { //камень
			"X": 3, //проиграть - ножницы
			"Y": 4, //ничья - камень
			"Z": 8, //победить - бумага
		},
		"B": { //бумага
			"X": 1, //проиграть - камень
			"Y": 5, //ничья - бумага
			"Z": 9, //победить - ножницы
		},
		"C": { //ножницы
			"X": 2, //проиграть - бумага
			"Y": 6, //ничья - ножницы
			"Z": 7, //победить - камень
		},
	}

	file, err := os.Open(inputFileName)

	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	points := 0

	for scanner.Scan() {
		line := scanner.Text()

		parts := strings.Split(line, " ")

		points += rules[parts[0]][parts[1]]
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Println(points)
}
