package main

import (
	"bufio"
	"container/list"
	"fmt"
	"log"
	"math/big"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

const inputFileName = "input"

const rounds = 10000

type Monkey struct {
	items               *list.List
	OperationMultiply   bool
	OperationValue      int
	OperationWithItself bool
	Test                int
	ResultTrue          int
	ResultFalse         int
}

func NewMonkey(
	operationMultiply, OperationWithItself bool,
	operationValue, test, resultTrue, resultFalse int,
) *Monkey {
	return &Monkey{
		items:               list.New(),
		OperationMultiply:   operationMultiply,
		OperationValue:      operationValue,
		OperationWithItself: OperationWithItself,
		Test:                test,
		ResultTrue:          resultTrue,
		ResultFalse:         resultFalse,
	}
}

func (m *Monkey) AddItem(item *big.Int) {
	m.items.PushBack(item)
}

type Game struct {
	monkeys  []*Monkey
	business map[int]int
	divider  *big.Int
}

func NewGame() *Game {
	return &Game{
		monkeys:  make([]*Monkey, 0),
		business: make(map[int]int),
	}
}

func (g *Game) UpdateUnitedDivider() {
	res := 1

	for _, m := range g.monkeys {
		res *= m.Test
	}

	g.divider = big.NewInt(int64(res))
}

func (g *Game) GetMonkeyBusiness() int {
	values := make([]int, len(g.business))

	for i, v := range g.business {
		values[i] = v
	}

	sort.SliceStable(values, func(i, j int) bool {
		return values[i] > values[j]
	})

	return values[0] * values[1]
}

func (g *Game) PlayRound() {
	for i, monkey := range g.monkeys {
		for monkey.items.Len() > 0 {
			item := monkey.items.Front()
			monkey.items.Remove(item)

			newVal := item.Value.(*big.Int)

			if monkey.OperationMultiply {
				if monkey.OperationWithItself {
					newVal.Mul(newVal, newVal)
				} else {
					tmp := big.NewInt(int64(monkey.OperationValue))
					newVal.Mul(newVal, tmp)
				}
			} else {
				if monkey.OperationWithItself {
					newVal.Add(newVal, newVal)
				} else {
					tmp := big.NewInt(int64(monkey.OperationValue))
					newVal.Add(newVal, tmp)
				}
			}

			//newVal /= 3
			tmpQ := big.NewInt(0)
			modulus0 := big.NewInt(0)
			tmpQ.DivMod(newVal, g.divider, modulus0)

			newVal = modulus0

			modulus := big.NewInt(0)
			tmp := big.NewInt(int64(monkey.Test))

			tmpQ.DivMod(newVal, tmp, modulus)

			if modulus.Int64() == 0 {
				g.monkeys[monkey.ResultTrue].AddItem(newVal)
			} else {
				g.monkeys[monkey.ResultFalse].AddItem(newVal)
			}

			g.business[i]++
		}
	}
}

func (g *Game) ParseInput(
	itemsLine string,
	operationLine string,
	testLine string,
	resultTrueLine string,
	resultFalseLine string,
) {
	itemsParts := strings.Split(itemsLine, ":")
	itemsPart := strings.TrimSpace(itemsParts[1])
	items := strings.Split(itemsPart, ", ")

	operationExpr := regexp.MustCompile(`^\s+Operation:\snew\s=\sold\s([\+\*])\s(old|[\d]+)$`)
	operationMatches := operationExpr.FindStringSubmatch(operationLine)

	multiply := operationMatches[1] == "*"
	itself := operationMatches[2] == "old"
	value := 0
	if !itself {
		value, _ = strconv.Atoi(operationMatches[2])
	}

	testExpr := regexp.MustCompile(`^\s+Test:\sdivisible\sby\s([\d]+)$`)
	testMatches := testExpr.FindStringSubmatch(testLine)

	test, _ := strconv.Atoi(testMatches[1])

	trueResultExpr := regexp.MustCompile(`^\s+If\strue:\sthrow\sto\smonkey\s([\d]+)$`)
	falseResultExpr := regexp.MustCompile(`^\s+If\sfalse:\sthrow\sto\smonkey\s([\d]+)$`)

	trueResultMatches := trueResultExpr.FindStringSubmatch(resultTrueLine)
	falseResultMatches := falseResultExpr.FindStringSubmatch(resultFalseLine)

	trueResult, _ := strconv.Atoi(trueResultMatches[1])
	falseResult, _ := strconv.Atoi(falseResultMatches[1])

	monkey := NewMonkey(
		multiply,
		itself,
		value,
		test,
		trueResult,
		falseResult,
	)

	for _, item := range items {
		worryLevel, _ := strconv.Atoi(item)
		monkey.AddItem(big.NewInt(int64(worryLevel)))
	}

	g.monkeys = append(g.monkeys, monkey)
}

func main() {
	file, err := os.Open(inputFileName)

	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	game := NewGame()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 {
			continue
		}
		scanner.Scan()
		itemsLine := scanner.Text()
		scanner.Scan()
		operationLine := scanner.Text()
		scanner.Scan()
		testLine := scanner.Text()
		scanner.Scan()
		resultTrueLine := scanner.Text()
		scanner.Scan()
		resultFalseLine := scanner.Text()

		game.ParseInput(
			itemsLine,
			operationLine,
			testLine,
			resultTrueLine,
			resultFalseLine,
		)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	game.UpdateUnitedDivider()

	for i := 0; i < rounds; i++ {
		fmt.Printf("round %d\r", i)
		game.PlayRound()
	}

	fmt.Printf("\n%v\n", game.business)
	fmt.Printf("%v\n", game.GetMonkeyBusiness())
}
