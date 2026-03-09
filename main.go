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

	audioHandler := fs.NewAudioHandler()
	artHandler := tags.NewArtHandler()

	err := wails.Run(&options.App{
		Title:  "harmonic",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
			Middleware: func(next http.Handler) http.Handler {
				return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					switch r.URL.Path {
					case "/audio":
						audioHandler.ServeHTTP(w, r)
					case "/art":
						artHandler.ServeHTTP(w, r)
					default:
						next.ServeHTTP(w, r)
					}
				})
			},
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
