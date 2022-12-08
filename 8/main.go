package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
)

const inputFileName = "input"

type Grid struct {
	grid    [][]int
	visible map[string]bool
}

func NewGrid(size int) *Grid {
	g := make([][]int, size)

	for i := range g {
		g[i] = make([]int, size)
	}

	return &Grid{
		grid:    g,
		visible: make(map[string]bool),
	}
}

func (g *Grid) FillRow(row int, line string) {
	for i, c := range line {
		n, _ := strconv.Atoi(string(c))

		g.grid[row][i] = n
	}
}

func (g *Grid) FindBestTreeScore() int {
	max := 0
	for row := range g.grid {
		for col := range g.grid {
			s := g.countTreeScore(row, col)
			if s > max {
				fmt.Printf("Tree %d:%d (%d) has score %d\n", row, col, g.grid[row][col], s)
				max = s
			}
		}
	}
	return max
}

func (g *Grid) countTreeScore(row, col int) int {
	if row == 0 || col == 0 {
		return 0
	}

	return g.distanceToClosestTree(row, col, 1, 0) *
		g.distanceToClosestTree(row, col, -1, 0) *
		g.distanceToClosestTree(row, col, 0, 1) *
		g.distanceToClosestTree(row, col, 0, -1)
}

func (g *Grid) distanceToClosestTree(row, col, rowDir, colDir int) int {
	size := len(g.grid)
	h := g.grid[row][col]
	i := 1
	for true {
		rowI := row + i*rowDir
		colI := col + i*colDir
		if rowI == 0 || rowI >= size-1 ||
			colI == 0 || colI >= size-1 {
			break
		}

		if g.grid[rowI][colI] >= h {
			break
		}

		i++
	}

	return i
}

func (g *Grid) CountVisibleTrees() int {
	for i := range g.grid {
		g.checkVisibleTreesForRow(i, 1)
		g.checkVisibleTreesForRow(i, -1)
		g.checkVisibleTreesForCol(i, 1)
		g.checkVisibleTreesForCol(i, -1)
	}

	return len(g.visible)
}

func (g *Grid) markVisible(row, col int) {
	key := fmt.Sprintf("%d:%d", row, col)
	g.visible[key] = true
}

func (g *Grid) checkVisibleTreesForRow(row int, dir int) {
	highest := -1
	size := len(g.grid)

	for i := 0; i < size; i++ {
		col := i
		if dir < 0 {
			col = size - i - 1
		}

		if col == 0 || row == 0 || g.grid[row][col] > highest {
			g.markVisible(row, col)

			highest = g.grid[row][col]
		}
	}
}

func (g *Grid) checkVisibleTreesForCol(col int, dir int) {
	highest := -1
	size := len(g.grid)

	for i := 0; i < size; i++ {
		row := i
		if dir < 0 {
			row = size - i - 1
		}

		if col == 0 || row == 0 || g.grid[row][col] > highest {
			g.markVisible(row, col)

			highest = g.grid[row][col]
		}
	}
}

func main() {
	file, err := os.Open(inputFileName)

	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	isFirstLine := true

	var grid *Grid

	row := 0
	for scanner.Scan() {
		line := scanner.Text()

		if isFirstLine {
			isFirstLine = false
			grid = NewGrid(len(line))
		}

		grid.FillRow(row, line)
		row++
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Println(grid.CountVisibleTrees())

	fmt.Println(grid.FindBestTreeScore())
}
