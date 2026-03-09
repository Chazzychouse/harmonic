package main

import (
	"embed"
	"log"
	"net"
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
	// Start a dedicated HTTP server for media (audio + art).
	// The wails:// custom scheme is not supported by GStreamer, so we must
	// serve media over a real http:// address that souphttpsrc can fetch.
	mediaMux := http.NewServeMux()
	mediaMux.Handle("/audio", fs.NewAudioHandler())
	mediaMux.Handle("/art", tags.NewArtHandler())

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		log.Fatal("media server:", err)
	}
	go http.Serve(ln, mediaMux) //nolint:errcheck

	mediaPort := ln.Addr().(*net.TCPAddr).Port
	fsService := fs.NewFsService(mediaPort)

	err = wails.Run(&options.App{
		Title:  "harmonic",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
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
