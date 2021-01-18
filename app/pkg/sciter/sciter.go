package sciter

import (
	"fmt"
	"log"
	"path/filepath"
	"strings"

	b "../box"
	u "../utility"

	"github.com/sciter-sdk/go-sciter"
)

func ButtonPress(args ...*sciter.Value) *sciter.Value {

	fp, outputPath := args[0].String(), args[1].String() // file path

	if  outputPath == "" {
		outputPath = "./"
	}

	if outputPath[len(outputPath)-1] != '/' {
		outputPath += "/"
	}

	log.Println("file path:", fp)
	log.Println("output path:", outputPath)

	if fp[0] == '[' {
		fp = u.DelChar(fp, 0)
		if fp[len(fp)-1] == ']' {
			fp = u.DelChar(fp, len(fp)-1)
		}
	}

	fp = strings.ReplaceAll(fp, string('"'), "")
	paths := strings.Split(fp, ",")

	var product []b.Product
	var box []b.Box

	for _, v := range paths {
		if filepath.Ext(v) != ".plt" {
			log.Println("err:", v, "is not a *.plt file")
			continue
		}

		p := b.Product{Source: v}
		p.GetDimensions()
		p.Name = strings.TrimSuffix(filepath.Base(p.Source), filepath.Ext(p.Source))
		fmt.Println(p.Source, "dimensions: ", p.Size)

		tmp := b.Box{Content: p}
		tmp.DefaultAddSpace()
		tmp.CalculateSize()

		product = append(product, p)
		box = append(box, tmp)
	}

	rack := b.ShelfPack(box, u.IntPair{X: 1500, Y: 2500})

	originPoint := b.Point2d{X: 0, Y: 0}
	for _, v := range rack {
		for _, w := range v {
			w.DrawBox(outputPath+"box_"+w.Content.Name+".plt", originPoint)

			originPoint.X += w.Size.X
		}
		originPoint.Y += v[0].Size.Y
	}

	return sciter.NullValue()
}
