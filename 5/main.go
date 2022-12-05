package main

import (
	"1_1/5/stack"
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

const inputFileName = "input"

func parseStacks(stacks []string) []stack.Stack {
	numbers := stacks[len(stacks)-1]

	numberParts := strings.Split(strings.TrimSpace(numbers), "   ")

	res := make([]stack.Stack, len(numberParts))

	for i := len(stacks) - 2; i >= 0; i-- {
		for j := 0; j < len(numberParts); j++ {
			idx := 4*j + 1

			if len(stacks[i]) >= idx {
				letter := stacks[i][idx]
				if string(letter) != " " {
					res[j].Push(string(letter))
				}
			}
		}
	}

	return res
}

func doMove1(line string, stacks []stack.Stack) {
	//cut "move"
	line = line[5:]

	pts1 := strings.Split(line, " from ")

	count, _ := strconv.Atoi(pts1[0])

	pts2 := strings.Split(pts1[1], " to ")

	from, _ := strconv.Atoi(pts2[0])
	to, _ := strconv.Atoi(pts2[1])

	from = from - 1
	to = to - 1

	for i := 0; i < count; i++ {
		item := stacks[from].Pop()
		stacks[to].Push(item)
	}
}

func doMove(line string, stacks []stack.Stack) {
	//cut "move"
	line = line[5:]

	pts1 := strings.Split(line, " from ")

	count, _ := strconv.Atoi(pts1[0])

	pts2 := strings.Split(pts1[1], " to ")

	from, _ := strconv.Atoi(pts2[0])
	to, _ := strconv.Atoi(pts2[1])

	from = from - 1
	to = to - 1

	tmp := stack.New()

	for i := 0; i < count; i++ {
		item := stacks[from].Pop()
		tmp.Push(item)
	}

	for i := 0; i < count; i++ {
		item := tmp.Pop()
		stacks[to].Push(item)
	}
}

func main() {

	file, err := os.Open(inputFileName)

	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	lines := make([]string, 0)

	moves := false

	var stacks []stack.Stack
	for scanner.Scan() {
		line := scanner.Text()

		if len(line) == 0 {
			stacks = parseStacks(lines)
			moves = true
			continue
		}

		if moves {
			doMove(line, stacks)
		} else {
			lines = append(lines, line)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	for _, s := range stacks {
		i := s.Peek()
		fmt.Printf("%v", i)
	}
}
