package tags

import (
	"net/http"
	"path/filepath"
	"strings"
)

var artAudioExts = map[string]bool{
	".mp3": true, ".flac": true, ".wav": true, ".aiff": true,
	".aif": true, ".m4a": true, ".ogg": true, ".opus": true,
	".wv": true, ".alac": true, ".aac": true, ".wma": true,
}

// NewArtHandler returns an http.Handler that serves embedded album art.
// GET /art?path=<absolute-audio-file-path>
func NewArtHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		path := r.URL.Query().Get("path")
		if path == "" {
			http.Error(w, "missing path parameter", http.StatusBadRequest)
			return
		}

		ext := strings.ToLower(filepath.Ext(path))
		if !artAudioExts[ext] {
			http.Error(w, "unsupported file type", http.StatusBadRequest)
			return
		}

		data, mime, err := ReadArt(path)
		if err != nil {
			http.NotFound(w, r)
			return
		}

		w.Header().Set("Content-Type", mime)
		w.Header().Set("Cache-Control", "public, max-age=86400")
		w.Write(data)
	})
}
