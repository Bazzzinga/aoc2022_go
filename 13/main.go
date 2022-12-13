package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

const inputFileName = "input"

type Comparable interface {
	isLower(than Comparable) *bool
	ToString() string
}

type Int struct {
	Value int
}

func (i *Int) ToString() string {
	return fmt.Sprintf("%d", i.Value)
}

func (i *Int) isLower(than Comparable) *bool {
	switch t := than.(type) {
	case *List:
		tmp := &List{
			Value: []Comparable{
				&Int{Value: i.Value},
			},
		}

		return tmp.isLower(than)
	case *Int:
		ok := false
		if i.Value > t.Value {
			return &ok
		}
		if i.Value < t.Value {
			ok = true
			return &ok
		}
	}

	return nil
}

type List struct {
	Value []Comparable
}

func (l *List) ToString() string {
	res := make([]string, len(l.Value))

	for i, e := range l.Value {
		res[i] = e.ToString()
	}

	return "[" + strings.Join(res, ",") + "]"
}

func (l *List) isLower(than Comparable) *bool {
	switch t := than.(type) {
	case *Int:
		tmp := &List{
			Value: []Comparable{
				&Int{Value: t.Value},
			},
		}

		ok := l.isLower(tmp)

		if ok != nil {
			return ok
		}
	case *List:
		ln := len(t.Value)
		for i, e := range l.Value {
			if i < ln {
				ok := e.isLower(t.Value[i])

				if ok != nil {
					return ok
				}
			}

			if i > ln-1 {
				ok := false
				return &ok
			}
		}

		if len(l.Value) > len(t.Value) {
			ok := false
			return &ok
		}

		if len(l.Value) < len(t.Value) {
			ok := true
			return &ok
		}
	}

	return nil
}

type Pair struct {
	Left  *List
	Right *List
}

func NewPair(leftLine, rightLine string) *Pair {
	p := Pair{}

	p.parseLine(leftLine, true)
	p.parseLine(rightLine, false)

	return &p
}

func (p *Pair) parseLine(line string, left bool) {
	var data []interface{}

	_ = json.Unmarshal([]byte(line), &data)

	list := List{
		Value: make([]Comparable, 0),
	}

	list.addData(data)

	if left {
		p.Left = &list
	} else {
		p.Right = &list
	}

	if line != list.ToString() {
		fmt.Println("Error parsing " + line)
	}
}

func (l *List) addData(data []interface{}) {
	for _, t := range data {
		switch reflect.TypeOf(t).Kind() {
		case reflect.Slice:
			s := reflect.ValueOf(t)
			subData := make([]interface{}, s.Len())
			for i := 0; i < s.Len(); i++ {
				subData[i] = s.Index(i)
			}

			list := List{
				Value: make([]Comparable, 0),
			}

			list.addData(subData)
			l.Value = append(l.Value, &list)
		case reflect.Struct:
			st := fmt.Sprintf("%v", t)

			if st[0:1] == "[" {
				var subData []interface{}

				line := strings.Replace(st, " ", ",", -1)

				_ = json.Unmarshal([]byte(line), &subData)

				list := List{
					Value: make([]Comparable, 0),
				}

				list.addData(subData)
				l.Value = append(l.Value, &list)
			} else {
				v, _ := strconv.Atoi(st)
				i := Int{Value: v}
				l.Value = append(l.Value, &i)
			}
		default:
			s := reflect.ValueOf(t)
			v := s.Float()
			i := Int{Value: int(v)}
			l.Value = append(l.Value, &i)
		}
	}
}

func (p *Pair) isOrderCorrect() bool {
	ok := p.Left.isLower(p.Right)

	if ok == nil {
		panic("oops")
	}

	return *ok
}

func main() {
	file, err := os.Open(inputFileName)

	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	pairs := make([]*Pair, 0)

	for scanner.Scan() {
		line1 := scanner.Text()
		if len(line1) == 0 {
			continue
		}
		scanner.Scan()
		line2 := scanner.Text()

		p := NewPair(line1, line2)

		pairs = append(pairs, p)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	res := 0

	for i, p := range pairs {
		if p.isOrderCorrect() {
			res += i + 1
		}
	}

	fmt.Println(res)

	dividersPair := NewPair("[[2]]", "[[6]]")

	pairs = append(pairs, dividersPair)

	allPackets := make([]*List, len(pairs)*2)

	for i, pair := range pairs {
		allPackets[2*i] = pair.Left
		allPackets[2*i+1] = pair.Right
	}

	sort.SliceStable(allPackets, func(i, j int) bool {
		less := allPackets[i].isLower(allPackets[j])

		if less == nil {
			return false
		}

		return *less
	})

	i1 := 0
	i2 := 0

	for i, p := range allPackets {
		if p.ToString() == "[[2]]" {
			i1 = i + 1
		} else if p.ToString() == "[[6]]" {
			i2 = i + 1
		}
	}

	fmt.Println(i1 * i2)
}
