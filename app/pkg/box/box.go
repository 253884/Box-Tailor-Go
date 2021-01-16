package box

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
)

type minS struct {
	x int
	y int
}

type maxS struct {
	x int
	y int
}

type Dimensions struct {
	x int
	y int
}

type Box struct {
	source string
	size   Dimensions
}

const MaxUint = ^uint(0)
//const MinUint = 0
const MaxInt = int(MaxUint >> 1)
const MinInt = -MaxInt - 1

const Unit = 40 // 40 points per mm in HPGL

func Check(err error) {
	if err != nil {
		panic(err)
	}
}

func LessOrEqual(boxes []Box, target int) int {
	var (
		l = 0
		r = len(boxes) - 1
	)

	for l < r {

		m := (l + r + 1) / 2

		if boxes[m].size.x > target {
			r = m - 1
		} else {
			l = m
		}
	}
	if boxes[l].size.x > target {
		return -1
	}
	return l
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

func GetNumbers(s string) []string {
	re := regexp.MustCompile("[0-9]+")
	return re.FindAllString(s, -1)
}

func GetDimensions(path string) Dimensions { // TO BE FINISHED

	if extension := filepath.Ext(path); extension != ".plt" {
		dimensions := Dimensions{-1, -1}
		return dimensions
	}

	file, err := os.Open(path)
	Check(err)
	defer func() {
		err := file.Close()
		Check(err)
	}()

	extremes := struct {
		min minS
		max maxS
	}{
		minS{MaxInt, MaxInt},
		maxS{MinInt, MinInt},
	}
	
	dimensions := Dimensions{}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line[0] == 'P' { // PD [pen down] set coordinates
			if line[1] == 'D' {
				stringSlice := GetNumbers(scanner.Text())

				for i, v := range stringSlice {
					v, err := strconv.Atoi(v)
					Check(err)

					if i%2 == 0 {
						extremes.min.x = Min(v, extremes.min.x)
						extremes.max.x = Max(v, extremes.max.x)
					} else {
						extremes.min.y = Min(v, extremes.min.y)
						extremes.max.y = Max(v, extremes.max.y)
					}
				}
			}
		}
		dimensions.x, dimensions.y = (extremes.max.x-extremes.min.x)/Unit, (extremes.max.y-extremes.min.y)/Unit
	}
	err = scanner.Err()
	Check(err)

	return dimensions
}
