package utility

import (
	"log"
	"regexp"
)

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

// Check checks if any errors occurred.
func Check(err error) {
	if err != nil {
		panic(err)
	}
}

/*func Max(a, b int) int {
	if a < b {
		return b
	}
	return a
}*/ // to do MAX

// Min returns sm
/*func Min(a, b int) int {
	if a > b {
		return b
	}
	return a
}*/ // to do MIN

// Area calculates area/volume by multiplying input values.
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

// DelChar removes a char from the string.
func DelChar(s string, i int) string {
	r := []rune(s)
	return string(append(r[0:i], r[i+1:]...))
}

// GetNumbers looks for numbers in provided string and returns them as string slice separately.
func GetNumbers(s string) []string {
	//re := regexp.MustCompile("[0-9]+")
	re := regexp.MustCompile(`[-]?\d[\d,]*[\.]?[\d{2}]*`)
	return re.FindAllString(s, -1)
}