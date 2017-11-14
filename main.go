package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"time"

	sciter "github.com/sciter-sdk/go-sciter"
	sciterWindow "github.com/sciter-sdk/go-sciter/window"
)

const (
	assetPrefix   = "asset://"
	filePrefix    = "file://"
	frameLoadTick = 3 * time.Second
	// hardcoded to centre for my screen size (is there a builtin for this?)
	screenHeight = 1080
	screenWidth  = 1920
	appHeight    = 600
	appWidth     = 800
	// files
	templateHTMLPath = "ui/template.html"
)

var rscPrefix = assetPrefix

// if you set skip-loader to true then sciter will load from file correctly, the issue only presents when using custom resorce loader
var debugLoadFromFile = flag.Bool("skip-loader", false, "pass --skip-loader=t to bypass the custom asset loader and make everything work as expected")

func main() {

	if *debugLoadFromFile {
		rscPrefix = filePrefix
	}

	w, err := sciterWindow.New(sciter.SW_MAIN|sciter.SW_ALPHA|sciter.SW_ENABLE_DEBUG, bounds())
	if err != nil {
		log.Fatalf("error creating window: %s", err)
	}

	registerRoutes(w)

	w.SetCallback(&sciter.CallbackHandler{
		OnLoadData: func(ld *sciter.ScnLoadData) int {
			log.Println("OnLoadData call for: ", ld.Uri())
			return LoadData(ld)
		},
	})

	err = w.LoadFile(rscPrefix + templateHTMLPath)
	if err != nil {
		log.Fatalf("error loading template html: %s", err)
	}

	go loadFrameLoop(w)

	w.Show()

	w.Run()

}

func loadFrameLoop(w *sciterWindow.Window) {
	t := time.NewTicker(frameLoadTick)
	loadFile := 2
	for {
		select {
		case <-t.C:
			loadFile = 3 - loadFile
			_, err := w.Call("Template.LoadFrame", sciter.NewValue(fmt.Sprintf("%d", loadFile)))
			if err != nil {
				log.Fatalf("fatal call to Template.LoadFrame: %s", err)
			}
		}
	}

}

func bounds() *sciter.Rect {
	return sciter.NewRect(
		(int(screenHeight/2) - int(appHeight/2)),
		(int(screenWidth/2) - int(appWidth/2)),
		appWidth,
		appHeight,
	)
}

func registerRoutes(w *sciterWindow.Window) {
	w.DefineFunction("Log", func(args ...*sciter.Value) *sciter.Value {
		if len(args) == 1 {
			log.Println(fmt.Sprintf("frontend log: %s", args[0]))
		}
		return sciter.NullValue()
	})
}

func LoadData(ld *sciter.ScnLoadData) int {
	var (
		data []byte
		err  error
	)

	uri := ld.Uri()

	switch {
	case strings.HasPrefix(uri, assetPrefix):
		log.Println("asset requested: " + string(uri[len(assetPrefix):]))
		data, err = ioutil.ReadFile(uri[len(assetPrefix):])

	default:
		log.Println("unrecognised asset, passing through to sciter: " + uri)
		return sciter.LOAD_OK
	}

	if err != nil {
		log.Println("error loading asset: " + err.Error())
		return sciter.LOAD_DISCARD
	}

	if len(data) < 1 {
		log.Println("empty data retuned from asset call")
		return sciter.LOAD_DISCARD
	}

	ld.SetData(data)

	return sciter.LOAD_OK
}
