package utility

import (
	"log"
	"regexp"
)

type MinS struct {
	X int
	Y int
}

type MaxS struct {
	X int
	Y int
}

const (
	MaxUint = ^uint(0)
	MinUint = 0
	MaxInt  = int(MaxUint >> 1)
	MinInt  = -MaxInt - 1
)

type IntPair struct {
	X, Y int
}

type Extremes struct {
	Min IntPair
	Max IntPair
}

func Check(err error) {
	if err != nil {
		panic(err)
	}
}

func Max(a, b int) int {
	if a < b {
		return b
	}
	return a
}

func Min(a, b int) int {
	if a > b {
		return b
	}
	return a
}

func Area(arr ...int) int {
	if len(arr) == 0 {
		log.Panic("err: no arguments provided")
	}
	result := 1
	for _, v := range arr {
		result *= v
	}
	return result
}

func DelChar(s string, i int) string {
	r := []rune(s)
	return string(append(r[0:i], r[i+1:]...))
}

func GetNumbers(s string) []string {
	re := regexp.MustCompile("[0-9]+")
	return re.FindAllString(s, -1)
}