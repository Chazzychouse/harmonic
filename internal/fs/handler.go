package fs

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
)

var extMIME = map[string]string{
	".mp3":  "audio/mpeg",
	".flac": "audio/flac",
	".wav":  "audio/wav",
	".aiff": "audio/aiff",
	".aif":  "audio/aiff",
	".m4a":  "audio/mp4",
	".ogg":  "audio/ogg",
	".opus": "audio/ogg; codecs=opus",
	".aac":  "audio/aac",
}

// NewAudioHandler returns an http.Handler that streams audio files.
// GET /audio?path=<absolute-audio-file-path>
func NewAudioHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		filePath := r.URL.Query().Get("path")
		if filePath == "" {
			http.Error(w, "Missing path parameter", http.StatusBadRequest)
			return
		}
		ext := filepath.Ext(filePath)
		mime, ok := extMIME[ext]
		if !ok {
			http.Error(w, "Unsupported file type", http.StatusBadRequest)
			return
		}
		f, err := os.Open(filePath)
		if err != nil {
			if os.IsNotExist(err) {
				http.Error(w, "File not found", http.StatusNotFound)
			} else {
				http.Error(w, fmt.Sprintf("Error opening file: %v", err), http.StatusInternalServerError)
			}
			return
		}
		defer f.Close()
		info, err := f.Stat()
		if err != nil {
			http.Error(w, "Error reading file", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", mime)
		w.Header().Set("Accept-Ranges", "bytes")
		http.ServeContent(w, r, info.Name(), info.ModTime(), f)
	})
}
