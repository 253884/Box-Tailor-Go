package box

import (
	"bufio"
	"os"
	"path/filepath"
	"sort"
	"strconv"

	u "../utility"
)

const Unit = 40 // 40 points per mm in HPGL

var (
	WallThk = 5
)

type Point2d struct {
	X, Y int
}

// move pen using PU
func (p *Point2d) move(x, y int) string {
	(*p).X, (*p).Y = p.X+x, p.Y+y
	return "PU:" + strconv.Itoa(p.X*Unit) + "," + strconv.Itoa(p.Y*Unit) + ";"
}

// move pen using PD -> draws a line
func (p *Point2d) line(x, y int) string {
	(*p).X, (*p).Y = p.X+x, p.Y+y
	return "PD:" + strconv.Itoa(p.X*Unit) + "," + strconv.Itoa(p.Y*Unit) + ";"
}

// return pen to 0,0
func (p *Point2d) toOrigin() string {
	(*p).X, (*p).Y = 0, 0
	return "PU:0,0;"
}


type Vector2d struct {
	Origin u.IntPair
	End    u.IntPair
}

type Dimensions struct {
	X, Y, Z int
}

type Product struct {
	Name   string
	Source string
	Size   Dimensions
}

// uses source of Product to determine it's dimensions
func (p *Product) GetDimensions() {

	if extension := filepath.Ext(p.Source); extension != ".plt" {
		(*p).Size = Dimensions{-1, -1, -1}
		return
	} // incorrect file type

	file, err := os.Open(p.Source) // open input file
	u.Check(err)
	defer func() {
		err := file.Close()
		u.Check(err)
	}() // close input file

	ext := u.Extremes{
		Min: u.IntPair{X: u.MaxInt, Y: u.MaxInt},
		Max: u.IntPair{X: u.MinInt, Y: u.MinInt},
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line[0] == 'P' { // PD [pen down] set coordinates
			if line[1] == 'D' {
				stringSlice := u.GetNumbers(scanner.Text())

				for i, v := range stringSlice {
					v, err := strconv.Atoi(v)
					u.Check(err)

					if i%2 == 0 {
						ext.Min.X = u.Min(v, ext.Min.X)
						ext.Max.X = u.Max(v, ext.Max.X)
					} else {
						ext.Min.Y = u.Min(v, ext.Min.X)
						ext.Max.Y = u.Max(v, ext.Max.X)
					}
				}
			}
		}
		(*p).Size.X, (*p).Size.Y, (*p).Size.Z = (ext.Max.X-ext.Min.X)/Unit, (ext.Max.Y-ext.Min.Y)/Unit, 20
	}
	err = scanner.Err()
	u.Check(err)
}


type Box struct {
	Content   Product
	Size      u.IntPair
	AddSpace Dimensions // additional space for foam ETC
}

func returnPLT(origin Point2d, arr ...int) []string {
	var result []string
	for i, v := range arr {
		if i%2 == 0 {
			result = append(result, origin.line(v, 0))
		} else {
			result = append(result, origin.line(0, v))
		}
	}
	if len(arr)%2 != 0 {
		result = append(result, origin.line(arr[len(arr)-1], 0))
	}
	return result
}

func (b *Box) CalculateSize() {
	x, y, z, w := b.Content.Size.X + b.AddSpace.X, b.Content.Size.Y + b.AddSpace.Y, b.Content.Size.Z + b.AddSpace.Z, WallThk
	(*b).Size.X, (*b).Size.Y = 2*x+4*z+6*w, y+2*z+2*w
}

func (b *Box) DefaultAddSpace() {
	(*b).AddSpace.X, (*b).AddSpace.Y, (*b).AddSpace.Z = 30, 30, 50
}

func (b *Box) DrawBox(outputPath string, origin Point2d) {
	x, y, z, w := b.Content.Size.X + b.AddSpace.X, b.Content.Size.Y + b.AddSpace.Y, b.Content.Size.Z + b.AddSpace.Z, WallThk

	leftWallX := 3*WallThk + 2*z
	CutOrigin := u.IntPair{X: leftWallX - y/2}
	if CutOrigin.X < 0 {
		CutOrigin.X = 0
	} // too long

	(*b).Size.X, (*b).Size.Y = 2*x+4*z+6*w, y+2*z+2*w // calculate box dimensions

	result := returnPLT( origin,
		leftWallX-CutOrigin.X+WallThk+x+WallThk+y, // x
			z, // y
		-z,
			WallThk,
		WallThk+z+WallThk,
			-(WallThk + z),
		x,
			WallThk+z,
		WallThk+z,
			y,
		-(WallThk + z),
			WallThk+z,
		-x,
			-(WallThk + z),
		-(WallThk + z + WallThk),
			WallThk,
		WallThk+z,
			z,
		-(leftWallX - CutOrigin.X + WallThk + x + WallThk + y),
			-z,
		z+WallThk,
			-WallThk,
		-leftWallX,
			-y,
		z+2*w+z,
			-w,
		-(z + w),
			-z)

	file, err := os.Create(outputPath)
	u.Check(err)

	_, err = file.WriteString("IN;\nLT;\nSP1;\n")
	u.Check(err)
	for _, v := range result{
		_, err := file.WriteString(v+"\n")
		u.Check(err)
	}
}


func lessOrEqual(boxes []Box, target int, s rune) int {

	if s == 'x' {
		for i, v := range boxes {
			if v.Size.X <= target {
				return i
			}
		}
		return -1
	} else {
		for i, v := range boxes {
			if v.Size.Y <= target {
				return i
			}
		}
		return -1
	}
}

func removeBox(b []Box, i int) []Box {
	if i == len(b)-1 {
		return b[:i]
	}
	return append(b[:i], b[i+1:]...)
}

func ShelfPack(boxes []Box, boardSize u.IntPair) [][]Box {
	sort.SliceStable(boxes, func(i, j int) bool {
		return boxes[i].Size.Y > boxes[j].Size.Y
	})

	var (
		shelf []Box
		rack [][]Box
		currPos int
	)

	for len(boxes) > 0 {
		i := lessOrEqual(boxes, boardSize.X - currPos, 'x')

		if i == -1 {
			currPos= 0

			rack = append(rack, shelf)
			shelf = []Box{}

			i = 0
		}

		shelf = append(shelf, boxes[i])
		boxes = removeBox(boxes, i)
		currPos += shelf[len(shelf)-1].Size.X

		if len(boxes) == 0 {
			rack = append(rack, shelf)
			shelf = []Box{}
			currPos = 0
		}
	}
	return rack
}