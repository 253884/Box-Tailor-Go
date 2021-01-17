package box

import (
	"bufio"
	"log"
	"os"
	"path/filepath"
	"strconv"

	u "../utility"
)

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
	origin u.IntPair
	end    u.IntPair
}

type Dimensions struct {
	x, y, z int
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
		(*p).Size.x, (*p).Size.y = (ext.Max.X-ext.Min.X)/Unit, (ext.Max.Y-ext.Min.Y)/Unit
	}
	err = scanner.Err()
	u.Check(err)
}

type Box struct {
	Content   Product
	size      u.IntPair
	CutOrigin u.IntPair   // relative position of first cut
	ToCut     []u.IntPair // altering x, y with each move
	ToEngrave []Vector2d  // origin point x,y; relative endpoint x,y
}

func (b *Box) addCuts(arr ...int) {

	for i, v := range arr {
		if i%2 != 0 {
			(*b).ToCut = append((*b).ToCut, u.IntPair{X: arr[i-1], Y: v})
		}
	}
	if len(arr)%2 != 0 {
		(*b).ToCut = append((*b).ToCut, u.IntPair{X: arr[len(arr)-1]})
	}
}

func (b *Box) CalculateShape() {
	x, y, z, w := b.Content.Size.x, b.Content.Size.y, b.Content.Size.z, WallThk

	leftWallX := 3*WallThk + 2*z
	(*b).CutOrigin = u.IntPair{X: leftWallX - y/2}
	if b.CutOrigin.X < 0 {
		(*b).CutOrigin.X = 0
	} // too long

	(*b).size.X, (*b).size.Y = 2*x+4*z+6*w, y+2*z+2*w // calculate box dimensions

	(*b).addCuts(
		leftWallX-b.CutOrigin.X+WallThk+x+WallThk+y, // x
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
		-(leftWallX - b.CutOrigin.X + WallThk + x + WallThk + y),
		-z,
		z+WallThk,
		-WallThk,
		-leftWallX,
		-y,
		z+2*w+z,
		-w,
		-(z + w),
		-z)
}

func (b *Box) Tailor(outputPath string, origin Point2d) {
	file, err := os.Create(outputPath) // create output file
	u.Check(err)
	defer func() {
		err = file.Close()
		u.Check(err)
	}() // close output file, check for err

	if b.Content.Size.x <= 0 {
		log.Println("This box has no Content!")
		return
	}

	if len(b.ToCut) <= 0 {
		(*b).CalculateShape()
	}

	_, _ = file.WriteString("IN;\nLT;\n") // initialize file
	_, _ = file.WriteString("SP1;\n")     // choose pen: cut
	_, _ = file.WriteString("PU:" + strconv.Itoa(b.CutOrigin.X*Unit) + "," + strconv.Itoa(b.CutOrigin.Y*Unit) + ";\n")
	for _, v := range b.ToCut {
		if v.X != 0 {
			_, _ = file.WriteString(origin.line(v.X, 0))
		}
		if v.Y != 0 {
			_, _ = file.WriteString(origin.line(0, v.Y))
		}
	}
	_, _ = file.WriteString("SP0;\n")
}

const Unit = 40 // 40 points per mm in HPGL

func LessOrEqual(boxes []Box, target int) int {
	var (
		l = 0
		r = len(boxes) - 1
	)

	for l < r {

		m := (l + r + 1) / 2

		if boxes[m].size.Y > target {
			r = m - 1
		} else {
			l = m
		}
	}
	if boxes[l].size.Y > target {
		return -1
	}
	return l
}

func GetDimensions(path string) Dimensions { // TO BE FINISHED

	if extension := filepath.Ext(path); extension != ".plt" {
		dimensions := Dimensions{-1, -1, -1}
		return dimensions
	}

	file, err := os.Open(path)
	u.Check(err)
	defer func() {
		err := file.Close()
		u.Check(err)
	}()

	extremes := struct {
		min u.MinS
		max u.MaxS
	}{
		u.MinS{X: u.MaxInt, Y: u.MaxInt},
		u.MaxS{X: u.MinInt, Y: u.MinInt},
	}

	dimensions := Dimensions{}

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
						extremes.min.X = u.Min(v, extremes.min.X)
						extremes.max.X = u.Max(v, extremes.max.X)
					} else {
						extremes.min.Y = u.Min(v, extremes.min.Y)
						extremes.max.Y = u.Max(v, extremes.max.Y)
					}
				}
			}
		}
		dimensions.x, dimensions.y = (extremes.max.X-extremes.min.X)/Unit, (extremes.max.Y-extremes.min.Y)/Unit
	}
	err = scanner.Err()
	u.Check(err)

	return dimensions
}
