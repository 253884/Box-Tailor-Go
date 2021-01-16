package main

import (
	"flag"
	"github.com/sciter-sdk/go-sciter"
	"github.com/sciter-sdk/go-sciter/rice"
	"github.com/sciter-sdk/go-sciter/window"
	"log"
	"path/filepath"
	"strings"

	b "./pkg/box"
)


func main() {
	flag.Parse()

// define window position and size
	winRect := sciter.NewRect(100, 100, 400, 300)

// create new window
	win, err := window.New(
		sciter.DefaultWindowCreateFlag,
		winRect)
	if err != nil {
		panic(err)
	}

// use 'rice' to handle html 'src' import
	rice.HandleDataLoad(win.Sciter)

// load app frontend
	win.SetTitle("Box Tailor")
	err = win.LoadFile("front/index.html")
	if err != nil {
		panic(err)
	}

// enable features
	ok := win.SetOption(
		sciter.SCITER_SET_SCRIPT_RUNTIME_FEATURES,
		sciter.ALLOW_FILE_IO|
		sciter.ALLOW_SOCKET_IO|
		sciter.ALLOW_EVAL|
		sciter.ALLOW_SYSINFO)
	if !ok {
		log.Println("failed to enable features")
	}

	win.DefineFunction("onFileSelected", onFileSelected)

	log.Println("Before")
	win.Show()
	log.Println("Show")
	win.Run()
	log.Println("Run")
}

func delChar(s string, i int) string {
	r := []rune(s)
	return string(append(r[0:i], r[i+1:]...))
}

func onFileSelected(args ...*sciter.Value) *sciter.Value {

	fp := args[0].String()
	log.Println(fp)
	if fp[0] == '[' {
		fp = delChar(fp, 0)
		if fp[len(fp)-1] == ']' {
			fp = delChar(fp, len(fp)-1)
		}
	}

	fp = strings.ReplaceAll(fp,string('"'), "")
	paths := strings.Split(fp, ",")

	for _, v := range paths {
		if filepath.Ext(v) == ".plt" {
			log.Println(filepath.Base(v), "dimensions: ", b.GetDimensions(v))
		} else {
			log.Println("err:" , v, "is not a *.plt file")
		}
	}

	return sciter.NullValue()
}