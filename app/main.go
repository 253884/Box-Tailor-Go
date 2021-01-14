package main

import (
	"./pkg/box"

	"github.com/sciter-sdk/go-sciter"
	"github.com/sciter-sdk/go-sciter/rice"
	"github.com/sciter-sdk/go-sciter/window"
)

func main() {
	winRect := sciter.NewRect(100, 100, 200, 200)

	win, err := window.New(sciter.SW_MAIN|sciter.SW_CONTROLS, winRect)
	if err != nil {
		panic(err)
	}
	rice.HandleDataLoad(win.Sciter)

	win.SetTitle("Box Tailor")

	err = win.LoadFile("./front/index.html")
	if err != nil {
		panic(err)
	}

	GetDimensions("glo_laser.plt")

	win.Show()
	win.Run()
}
