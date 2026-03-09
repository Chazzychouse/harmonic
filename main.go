package main

import (
	"embed"
	"net/http"

	"harmonic/internal/fs"
	"harmonic/internal/tags"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	fsService := fs.NewFsService()

	mux := http.NewServeMux()
	mux.Handle("/audio", fs.NewAudioHandler())
	mux.Handle("/art", tags.NewArtHandler())

	err := wails.Run(&options.App{
		Title:  "harmonic",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets:  assets,
			Handler: mux,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        fsService.SetContext,
		Bind: []interface{}{
			fsService,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
