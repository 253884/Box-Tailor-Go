package main

import (
	"log"

	//b "./pkg/box"
	//u "./pkg/utility"
	s "./pkg/sciter"

	"github.com/sciter-sdk/go-sciter"
	"github.com/sciter-sdk/go-sciter/rice"
	"github.com/sciter-sdk/go-sciter/window"
)

func main() {
	// define window position and size
	winRect := sciter.NewRect(100, 100, 400, 300)

	// create new window
	win, err := window.New(
		sciter.SW_MAIN|
			sciter.SW_ENABLE_DEBUG|
			sciter.SW_CONTROLS|
			sciter.SW_RESIZEABLE|
			sciter.SW_TITLEBAR,
		winRect)
	if err != nil {
		panic(err)
	}

	// load app frontend
	win.SetTitle("Box Tailor")
	err = win.LoadFile("front/index.html")
	if err != nil {
		panic(err)
	}

	// use 'rice' to handle html 'src' import
	rice.HandleDataLoad(win.Sciter)

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

	win.DefineFunction("buttonPress", s.ButtonPress)

	log.Println("Before")
	win.Show()
	log.Println("Show")
	win.Run()
	log.Println("Run")
}
