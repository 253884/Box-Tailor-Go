package sct

import (
	b "../box"
	"../db"
	u "../utility"
	"path/filepath"

	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/sciter-sdk/go-sciter"
)

func getProducts(s string) []b.Product {
	var (
		products []b.Product
		boxCount int
		)

	rule, err := regexp.Compile(`"(.*?)"`) // look for text between quotes
	u.Check(err)

	s = u.RemoveBraces(s)

	if s[0] == '{' {
		if s[len(s)-1] == '}' {
			s = u.DelChar(s, 0)
			s = u.DelChar(s, len(s)-1)
		}
	}

	arr := strings.Split(s, ",")

	var p = b.Product{
		Name:     "",
		Source:   "",
		Size:     b.Dimensions{},
		AddSpace: b.Dimensions{},
		Type:     0,
	}
	for _, v := range arr { // start of product input
		tmp := rule.FindAllString(v, -1)

		var label, value string
		if len(tmp) >= 2 {
			label, value = tmp[0], tmp[1]
		} else if len(tmp) == 1 {
			num := u.GetNumbers(v)
			label, value = tmp[0], num[len(num)-1]
		} else {
			continue
		}

		label = u.RemoveQuotes(label)
		value = u.RemoveQuotes(value)

		//fmt.Println(label, value, p)

		switch label {
		case "name":
			if p.Name != "" {
				if p.Size.X > 0 && p.Size.Y > 0 && p.Size.Z > 0 {
					if p.Name == "<from_path>" {
						p.Name = strings.TrimSuffix(filepath.Base(p.Source), filepath.Ext(p.Source))
					} else if p.Name == "<default>" {
						p.Name = "box_" + strconv.Itoa(boxCount)
						boxCount++
					}

					products = append(products, p)
				}

				p = b.Product{
					Name:     "",
					Source:   "",
					Size:     b.Dimensions{},
					AddSpace: b.Dimensions{},
					Type:     0,
				}
			}
			p.Name = value
		case "path":
			p.Source = value
			p.GetDimensions()
		case "size_x":
			p.Size.X, err = strconv.Atoi(value)
			u.Check(err)
		case "size_y":
			p.Size.Y, err= strconv.Atoi(value)
			u.Check(err)
		case "size_z":
			p.Size.Z, err= strconv.Atoi(value)
			u.Check(err)
		case "add_spc_x":
			p.AddSpace.X, err = strconv.Atoi(value)
			u.Check(err)
		case "add_spc_y":
			p.AddSpace.Y, err = strconv.Atoi(value)
			u.Check(err)
		case "add_spc_z":
			p.AddSpace.Z, err = strconv.Atoi(value)
			u.Check(err)
		case "type":
			if value == "flap" {
				p.Type = 'f'
			} else {
				p.Type = 'm'
			}
		default:
			continue
		}
	}

	if p.Name != "" && p.Size.X > 0 && p.Size.Y > 0 && p.Size.Z > 0 {
		if p.Name == "<from_path>" {
			p.Name = strings.TrimSuffix(filepath.Base(p.Source), filepath.Ext(p.Source))
		} else if p.Name == "<default>" {
			p.Name = "box_" + strconv.Itoa(boxCount)
			boxCount++
		}
		products = append(products, p)
	}

	return products
}

func ButtonPress(args ...*sciter.Value) *sciter.Value {

	var (
		products []b.Product
		boxes []b.Box
		outputPath string
		boardSize u.IntPair
		err error
		)

	in := args[0].String()
	bs := args[1].String()

	bs = u.RemoveBraces(bs)
	sb := strings.Split(bs, ",")

	boardSize.X, err = strconv.Atoi(sb[0])
	u.Check(err)
	boardSize.Y, err = strconv.Atoi(sb[1])
	u.Check(err)
	log.Println(boardSize)

	outputPath = args[2].String()
	if outputPath == "" {
		outputPath = "./"
	} else if outputPath[len(outputPath)-1] != '/' {
		outputPath += "/"
	}


	products = getProducts(in)

	for i, v := range products {
		fmt.Println("product", i, ":", v )

		tmp := b.Box{Content: v}
		tmp.CalculateSize()

		boxes = append(boxes, tmp)
	}
	fmt.Println("output path:", outputPath)

	fmt.Println("before rack")
	rack := b.ShelfPack(boxes, boardSize)
	fmt.Println("after rack")
	b.SplitToBoards(rack, boardSize, outputPath)
	fmt.Println("after split")

	return sciter.NullValue()
}

func GetSettings(_ ...*sciter.Value) *sciter.Value {
	dataBase := db.AccessData()
	defer func() {
		err := dataBase.Close()
		u.Check(err)
	}()

	settings := db.ReadSettings(dataBase)

	var result string
	for _, v := range settings {
		result += strconv.Itoa(v) + "|"
	}

	return sciter.NewValue(result)
}

func ChangeSettings(args ...*sciter.Value) *sciter.Value {
	dataBase := db.AccessData()

	for _, v := range args {
		log.Println(v.Int())
	}

	for i, v := range args {
		db.EditSetting(dataBase, i+1, v.Int())
	}

	return sciter.NullValue()
}
