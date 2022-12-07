package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
)

const inputFileName = "input"

const sizeThreshold = 100_000
const filesystemSize = 70_000_000
const needSpace = 30_000_000

type Node struct {
	Name     string
	Size     uint
	Children map[string]*Node
	Parent   *Node
}

func NewNode(name string, size uint, parent *Node) *Node {
	return &Node{
		Name:     name,
		Size:     size,
		Parent:   parent,
		Children: make(map[string]*Node),
	}
}

func (n *Node) GetSize() uint {
	if n.Size > 0 {
		return n.Size
	}

	var res uint

	for _, c := range n.Children {
		res += c.GetSize()
	}

	n.Size = res

	return res
}

func (n *Node) TotalOfSmallDirs() uint {
	var total uint
	if len(n.Children) == 0 {
		return total
	}

	if n.Size <= sizeThreshold {
		total += n.Size
		fmt.Printf("dir %s with size %d counts  New total: %d\n", n.Name, n.Size, total)
	}

	for _, c := range n.Children {
		total += c.TotalOfSmallDirs()
	}

	return total
}

func (n *Node) getDirSizesBiggerThan(max uint, res []uint) []uint {
	if len(n.Children) == 0 {
		return res
	}

	if n.Size >= max {
		res = append(res, n.Size)
	}

	for _, c := range n.Children {
		res = c.getDirSizesBiggerThan(max, res)
	}

	return res
}

func (n *Node) BiggestSmallerThanDirSize(max uint) uint {
	sizes := n.getDirSizesBiggerThan(max, make([]uint, 0))

	min := sizes[0]

	for _, s := range sizes {
		if s < min {
			min = s
		}
	}

	return min
}

func processLine(line string, root, cursor *Node) *Node {
	if line == "$ cd /" {
		fmt.Println("move cursor to root")
		return root
	}

	if line == "$ ls" {
		fmt.Println("peek directory")
		return cursor
	}

	if line == "$ cd .." {
		if cursor.Parent == nil {
			return cursor
		}
		fmt.Println("move cursor to parent")
		return cursor.Parent
	}

	reCd := regexp.MustCompile(`^\$\scd\s([.\w]+)$`)
	match := reCd.FindStringSubmatch(line)

	if len(match) == 2 {
		fmt.Println("move cursor to dir " + match[1])
		return cursor.Children[match[1]]
	}

	reDir := regexp.MustCompile(`^dir\s([.\w]+)$`)
	match = reDir.FindStringSubmatch(line)

	if len(match) == 2 {
		_, ok := cursor.Children[match[1]]
		if ok {
			return cursor
		}
		fmt.Println("add dir " + match[1])
		dir := NewNode(match[1], 0, cursor)
		cursor.Children[match[1]] = dir
	}

	reFile := regexp.MustCompile(`^(\d+)\s([.\w]+)$`)
	match = reFile.FindStringSubmatch(line)

	if len(match) == 3 {
		_, ok := cursor.Children[match[1]]
		if ok {
			return cursor
		}
		fmt.Println("add file " + match[2] + " with size " + match[1])
		size, _ := strconv.Atoi(match[1])
		dir := NewNode(match[2], uint(size), cursor)
		cursor.Children[match[2]] = dir
	}

	return cursor
}

func main() {
	file, err := os.Open(inputFileName)

	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	root := NewNode("/", 0, nil)
	cursor := root

	scanner := bufio.NewScanner(file)

	var line string
	for scanner.Scan() {
		line = scanner.Text()

		cursor = processLine(line, root, cursor)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("total size: %d\n", root.GetSize())
	fmt.Printf("total size of small dirs: %d\n", root.TotalOfSmallDirs())

	needToDelete := needSpace + root.GetSize() - filesystemSize

	fmt.Printf("need to delete dirs of size: %d\n", needToDelete)

	fmt.Printf("must delete: %d\n", root.BiggestSmallerThanDirSize(needToDelete))
}
