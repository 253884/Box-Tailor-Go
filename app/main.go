package main

import (
	"github.com/sciter-sdk/go-sciter/rice"
	"log"

	"./pkg/box"
	"./pkg/db"
	s "./pkg/sct"

	"github.com/sciter-sdk/go-sciter"
	"github.com/sciter-sdk/go-sciter/window"
)

func main() {
	dataBase := db.AccessData()
	defer func() {
		err := dataBase.Close()
		if err != nil {
			log.Println("main1 err:", err)
		}
	}()
	db.Initiate(dataBase)

	box.UpdateSettingValues()

	//db.AddSetting(dataBase, "WallThk", 5)
	//db.EditSetting(dataBase, 1, 33)
	//db.DeleteSetting(dataBase, 1)

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

	//log.Println("Before")
	win.Show()
	//log.Println("Show")
	win.Run()
	//log.Println("Run")
}
