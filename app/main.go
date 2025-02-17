package main

import (
	"log"

	"box-tailor-go/app/pkg/box"
	"box-tailor-go/app/pkg/db"
	s "box-tailor-go/app/pkg/sct"

	"github.com/sciter-sdk/go-sciter"
	"github.com/sciter-sdk/go-sciter/rice"
	"github.com/sciter-sdk/go-sciter/window"
)

func main() {
	// access database
	dataBase := db.AccessData()
	defer func() {
		err := dataBase.Close()
		if err != nil {
			log.Println("main1 err:", err)
		}
	}()
	// create settings table if there is no such table
	db.Initiate(dataBase)

	box.UpdateSettingValues()

	// define window position and size
	winRect := sciter.NewRect(100, 100, 800, 600)

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

	win.SetOption(sciter.SCITER_SET_DEBUG_MODE, 1)

	// use 'rice' to handle html 'src' import
	rice.HandleDataLoad(win.Sciter)

	// load app frontend
	win.SetTitle("Box Tailor")
	err = win.LoadFile("rice://front/index.html")
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

	win.DefineFunction("buttonPress", s.ButtonPress)
	win.DefineFunction("getSettings", s.GetSettings)
	win.DefineFunction("changeSettings", s.ChangeSettings)

	win.Show()
	win.Run()
}
