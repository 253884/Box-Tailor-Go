package sciter

import (
	"log"
	"path/filepath"
	"strconv"
	"strings"

	b "../box"
	u "../utility"

	"github.com/sciter-sdk/go-sciter"
)

func ButtonPress(args ...*sciter.Value) *sciter.Value {

	fp := args[0].String() // file path

	var outputPath string
	if len(args) > 1 {
		outputPath = args[1].String()
	}

	if outputPath[len(outputPath)-1] != '/' {
		outputPath += "/"
	}

	log.Println(fp, outputPath)

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
		log.Println(p.Source, "dimensions: ", p.Size)

		tmp := b.Box{Content: p}
		tmp.CalculateShape()

		product = append(product, p)
		box = append(box, tmp)
	}

	for i, v := range box {
		v.Tailor(outputPath+"box_"+v.Content.Name+"_"+strconv.Itoa(i)+".plt", b.Point2d{})
	}

	return sciter.NullValue()
}
