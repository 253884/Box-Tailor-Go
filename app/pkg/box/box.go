package box

import (
	"bufio"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"

	"box-tailor-go/app/pkg/db"
	u "box-tailor-go/app/pkg/utility"
)

const Unit = 40 // 40 points per mm in HPGL

var (
	WallThk   = 4
	margin    = u.IntPair{X: 10, Y: 10}          // 10mm margin
	boxDist   = u.IntPair{X: 1, Y: 1}            // 1mm space between boxes
	boxAddSpc = Dimensions{X: 60, Y: 60, Z: 120} // additional space for foam ETC
	defBoard  = u.IntPair{X: 3500, Y: 2500}
) // settings

func UpdateSettingValues() {
	dataBase := db.AccessData()
	defer func() {
		err := dataBase.Close()
		if err != nil {
			log.Println("box err:", err)
		}
	}()

	arr := db.ReadSettings(dataBase)
	log.Println("settings:", arr)
	WallThk = arr[0]
	margin.X = arr[1]
	margin.Y = arr[2]
	boxDist.X = arr[3]
	boxDist.Y = arr[4]
	boxAddSpc.X = arr[5]
	boxAddSpc.Y = arr[6]
	boxAddSpc.Z = arr[7]
	defBoard.X = arr[8]
	defBoard.Y = arr[9]
}

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
	Name     string
	Source   string
	Size     Dimensions
	AddSpace Dimensions // additional space for foam ETC
	Type     rune
}

type moveCut struct {
	x, y  int
	toCut bool
}

// GetDimensions reads through source of Product to determine it's dimensions.
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
						if v < ext.Min.X {
							ext.Min.X = v
						}
						if v > ext.Max.X {
							ext.Max.X = v
						}
						//ext.Min.X = u.Min(v, ext.Min.X)
						//ext.Max.X = u.Max(v, ext.Max.X)
					} else {
						if v < ext.Min.Y {
							ext.Min.Y = v
						}
						if v > ext.Max.Y {
							ext.Max.Y = v
						}
						//ext.Min.Y = u.Min(v, ext.Min.X)
						//ext.Max.Y = u.Max(v, ext.Max.X)
					}
				}
			}
		}
	}
	(*p).Size.X, (*p).Size.Y, (*p).Size.Z = (ext.Max.X-ext.Min.X)/Unit, (ext.Max.Y-ext.Min.Y)/Unit, 20
	err = scanner.Err()
	u.Check(err)
}

type Box struct {
	Content Product
	Size    u.IntPair
}

// CalculateSize calculates the size based on it's content.
func (b *Box) CalculateSize() {
	x, y, z, w := b.Content.Size.X+b.Content.AddSpace.X, b.Content.Size.Y+b.Content.AddSpace.Y, b.Content.Size.Z+b.Content.AddSpace.Z, WallThk
	if b.Content.Type == 'm' {
		(*b).Size.X, (*b).Size.Y = 2*x+4*z+6*w, y+2*z+2*w
	} else if b.Content.Type == 'f' {
		(*b).Size = u.IntPair{X: 2*x + 2*y + 4*w - w/2 + 20, Y: y + 2*w + z}
	}
}

// DefaultAddSpace sets Box's AddSpace to default values.
func (b *Box) DefaultAddSpace() {
	(*b).Content.AddSpace = boxAddSpc
}

// DrawBox writes lines of HPGL into file to draw box based on origin and type.
func (b *Box) DrawBox(file *os.File, origin Point2d, boxType rune) { // types: m - mailer; f - flap; l - lidded
	x, y, z, w := b.Content.Size.X+b.Content.AddSpace.X, b.Content.Size.Y+b.Content.AddSpace.Y, b.Content.Size.Z+b.Content.AddSpace.Z, WallThk
	var cut, engrave []string

	//fmt.Println(boxType, 'm')

	if boxType == 'm' {
		log.Println("drawing m")
		divided := y / 5
		leftWallX := 3*WallThk + 2*z
		CutOrigin := u.IntPair{X: leftWallX - y/2}
		if CutOrigin.X < 0 {
			CutOrigin.X = 0
		} // too long

		orgOrigin := origin
		origin.X, origin.Y = origin.X+CutOrigin.X, origin.Y+CutOrigin.Y

		(*b).Size.X, (*b).Size.Y = 2*x+4*z+7*w, y+2*z+2*w // calculate box dimensions

		cut = returnPLT(origin, 1,
			leftWallX-CutOrigin.X+WallThk+x+w+z, // x
			z, // y
			-(z + w),
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
			-(leftWallX - CutOrigin.X + WallThk + x + WallThk + z),
			-z,
			leftWallX-CutOrigin.X+WallThk,
			-WallThk,
			-leftWallX,
			-divided,
			-w,
			-divided,
			w,
			-(y - 4*divided),
			-w,
			-divided,
			w,
			-divided,
			leftWallX,
			-w,
			-(leftWallX - CutOrigin.X + WallThk),
			-z)
		//cut = append(cut, "PA:" + strconv.Itoa(origin.X) + " " + strconv.Itoa(origin.Y) + ";")
		engrave = engravePLT(orgOrigin, 2,
			moveCut{leftWallX + w/2, 0, false},
			moveCut{0, z, true},
			moveCut{0, w, false},
			moveCut{0, y, true},
			moveCut{0, w, false},
			moveCut{0, z, true},
			moveCut{w / 2, -(z + w/2), false},
			moveCut{x, 0, true},
			moveCut{w / 2, z + w/2, false},
			moveCut{0, -z, true},
			moveCut{0, -w, false},
			moveCut{0, -y, true},
			moveCut{0, -w, false},
			moveCut{0, -z, true},
			moveCut{-(x + w/2), z + w/2, false},
			moveCut{x, 0, true},
			moveCut{z + 2*w, 0, false},
			moveCut{x, 0, true},
			moveCut{w / 2, w / 2, false},
			moveCut{0, y, true},
			moveCut{-w / 2, w / 2, false},
			moveCut{-x, 0, true},
			moveCut{-w / 2, -w / 2, false},
			moveCut{0, -y, true})
	} else if boxType == 'f' {
		log.Println("drawing f")
		tab, a := u.IntPair{X: 20, Y: 10}, w+y/2
		(*b).Size = u.IntPair{X: 2*x + 2*y + 4*w - w/2 + tab.X, Y: 2*a + z}

		cut = engravePLT(origin, 1,
			moveCut{0, 0, false},
			moveCut{x, 0, true},
			moveCut{0, a, true},
			moveCut{w, 0, true},
			moveCut{0, -a, true},
			moveCut{y, 0, true},
			moveCut{0, a, true},
			moveCut{w, 0, true},
			moveCut{0, -a, true},
			moveCut{x, 0, true},
			moveCut{0, a, true},
			moveCut{w, 0, true},
			moveCut{0, -a, true},
			moveCut{y - w/2, 0, true},
			moveCut{0, a, true},
			//moveCut{w, 0, true},
			moveCut{tab.X, tab.Y, true},
			moveCut{0, z - 2*tab.Y, true},
			moveCut{-tab.X, tab.Y, true},
			//moveCut{-w, 0, true},
			moveCut{0, a, true},
			moveCut{-(y - w/2), 0, true},
			moveCut{0, -a, true},
			moveCut{-w, 0, true},
			moveCut{0, a, true},
			moveCut{-x, 0, true},
			moveCut{0, -a, true},
			moveCut{-w, 0, true},
			moveCut{0, a, true},
			moveCut{-y, 0, true},
			moveCut{0, -a, true},
			moveCut{-w, 0, true},
			moveCut{0, a, true},
			moveCut{-x, 0, true},
			moveCut{0, -(2*a + z), true})

		engrave = engravePLT(origin, 2,
			moveCut{0, a - w/2, false},
			moveCut{x, 0, true},
			moveCut{w, 0, false},
			moveCut{y, 0, true},
			moveCut{w, 0, false},
			moveCut{x, 0, true},
			moveCut{w, 0, false},
			moveCut{y - w/2, 0, true},
			moveCut{0, w / 2, false},
			moveCut{0, z, true},
			moveCut{0, w / 2, false},
			moveCut{-(y - w/2), 0, true},
			moveCut{-w, 0, false},
			moveCut{-x, 0, true},
			moveCut{-w, 0, false},
			moveCut{-y, 0, true},
			moveCut{-w, 0, false},
			moveCut{-x, 0, true},
			moveCut{x + w/2, -w / 2, false},
			moveCut{0, -z, true},
			moveCut{y + w, 0, false},
			moveCut{0, z, true},
			moveCut{x + w, 0, false},
			moveCut{0, -z, true})
		// draw flap box
	} else if boxType == 'l' {
		// draw lidded box
	} else {
		log.Println("err: Invalid box type.")
		return
	}

	for _, v := range cut {
		_, err := file.WriteString(v + "\n")
		u.Check(err)
	}
	for _, v := range engrave {
		_, err := file.WriteString(v + "\n")
		u.Check(err)
	}
}

// LessOrEqual looks for highest value in a slice of Boxes that is <= target where rune determines if it is x or y.
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

// RemoveBox removes Box[i] from the slice.
func removeBox(b []Box, i int) []Box {
	if i == len(b)-1 {
		return b[:i]
	}
	return append(b[:i], b[i+1:]...)
}

// ReturnPLT returns a slice of strings where each string is one HPGL line of code.
func returnPLT(origin Point2d, pen int, arr ...int) []string {
	var result []string
	log.Println(origin)
	result = append(result, "SP"+strconv.Itoa(pen)+";")
	result = append(result, origin.move(0, 0))
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

func engravePLT(origin Point2d, pen int, arr ...moveCut) []string {
	var result []string
	result = append(result, "SP"+strconv.Itoa(pen)+";")

	for _, v := range arr {
		if v.toCut {
			result = append(result, origin.line(v.x, v.y))
		} else {
			result = append(result, origin.move(v.x, v.y))
		}
	}
	return result
}

// ShelfPack creates rack with boxes where the rack cannot be wider than the board itself.
func ShelfPack(boxes []Box, boardSize u.IntPair) [][]Box {
	if boardSize.X <= 0 {
		boardSize.X = defBoard.X
	}
	if boardSize.Y <= 0 {
		boardSize.Y = defBoard.Y
	}

	sort.SliceStable(boxes, func(i, j int) bool {
		return boxes[i].Size.Y > boxes[j].Size.Y
	})

	var (
		shelf   []Box
		rack    [][]Box
		currPos int
	)

	for len(boxes) > 0 {
		i := lessOrEqual(boxes, boardSize.X-currPos, 'x')

		if i == -1 {
			currPos = 0

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

// CreateNewBoard creates a new board file.
func createNewBoard(outputFolder string, i int, boardSize u.IntPair) *os.File {
	outputFile, err := os.Create(outputFolder + "board_" + strconv.Itoa(i) + ".plt")
	if err != nil {
		panic(err)
	}

	_, err = outputFile.WriteString("IN;\nLT;\n" + "SP3;\nPD:" + strconv.Itoa(boardSize.X*Unit) + " 0;\nPD:" + strconv.Itoa(boardSize.X*Unit) + " " + strconv.Itoa(boardSize.Y*Unit) + ";\nPD:0 " + strconv.Itoa(boardSize.Y*Unit) + ";\nPD:0 0;\n")
	if err != nil {
		panic(err)
	}

	return outputFile
}

// SplitToBoards splits the rack into board(s).
func SplitToBoards(rack [][]Box, boardSize u.IntPair, outputFolder string) {

	var (
		//margin = u.IntPair{X: 10, Y: 10} // 10mm margin
		//boxDist = u.IntPair{X: 1, Y: 1} // 1mm space between boxes
		currPos      = Point2d{X: margin.X, Y: margin.Y}
		boardCounter = 0
	)

	if boardSize.X <= 0 {
		boardSize.X = defBoard.X
	}
	if boardSize.Y <= 0 {
		boardSize.Y = defBoard.Y
	}

	if rack[0][0].Size.Y > boardSize.Y {
		log.Println("err: Board is too small.")
		return
	}

	outputFile := createNewBoard(outputFolder, boardCounter, boardSize)
	defer func() {
		err := outputFile.Close()
		if err != nil {
			panic(err)
		}
	}()

	for _, v := range rack {
		if v[0].Size.Y <= boardSize.Y-currPos.Y-2*margin.Y {

			for _, w := range v {
				w.DrawBox(outputFile, currPos, w.Content.Type)
				currPos.X += w.Size.X + boxDist.X
			}
			currPos.X = margin.X
			currPos.Y += v[0].Size.Y + boxDist.Y
		} else {

			err := outputFile.Close()
			if err != nil {
				panic(err)
			}
			boardCounter++

			outputFile = createNewBoard(outputFolder, boardCounter, boardSize)

			currPos = Point2d{X: margin.X, Y: margin.Y}

			if v[0].Size.Y <= boardSize.Y-currPos.Y-2*margin.Y {
				for _, w := range v {
					w.DrawBox(outputFile, currPos, 'f')
					currPos.X += w.Size.X + boxDist.X
				}
				currPos.X = margin.X
				currPos.Y += v[0].Size.Y + boxDist.X
			}
		} // create new board file
	}
}
