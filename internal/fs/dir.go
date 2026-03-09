package fs

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"harmonic/internal/tags"

	wailsruntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

// AudioFile holds metadata about a discovered audio file.
type AudioFile struct {
	Title       string `json:"title"`        // from tag, falls back to filename
	Artist      string `json:"artist"`
	Album       string `json:"album"`
	AlbumArtist string `json:"album_artist"`
	Genre       string `json:"genre"`
	Year        int    `json:"year"`
	TrackNum    int    `json:"track_num"`
	TrackTotal  int    `json:"track_total"`
	DiscNum     int    `json:"disc_num"`
	DiscTotal   int    `json:"disc_total"`
	HasArt      bool   `json:"has_art"`
	FilePath    string `json:"file_path"`
	Ext         string `json:"ext"`
}

var audioExtensions = map[string]bool{
	".mp3": true, ".flac": true, ".wav": true, ".aiff": true,
	".aif": true, ".m4a": true, ".ogg": true, ".opus": true,
	".wv": true, ".alac": true, ".aac": true, ".wma": true,
}

// SelectDirectory opens a native OS directory picker and returns the chosen path.
// Returns an empty string if the user cancels.
func (s *FsService) SelectDirectory() string {
	dir, err := wailsruntime.OpenDirectoryDialog(s.ctx, wailsruntime.OpenDialogOptions{
		Title: "Select Music Library Folder",
	})
	if err != nil || dir == "" {
		return ""
	}
	return dir
}

func (s *FsService) ScanDirectory(dir string) ([]AudioFile, error) {
	if dir == "" {
		return nil, fmt.Errorf("no directory specified")
	}
	if _, err := os.Stat(dir); err != nil {
		return nil, fmt.Errorf("cannot access %q: %w", dir, err)
	}

	// Phase 1: collect file paths via WalkDir.
	type entry struct {
		path          string
		filenameTitle string
		ext           string
	}
	var entries []entry
	walkErr := filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("reading %q: %w", path, err)
		}
		if d.IsDir() {
			return nil
		}
		rawExt := filepath.Ext(d.Name())
		ext := strings.ToLower(rawExt)
		if !audioExtensions[ext] {
			return nil
		}
		entries = append(entries, entry{
			path:          path,
			filenameTitle: strings.TrimSuffix(d.Name(), rawExt),
			ext:           ext[1:],
		})
		return nil
	})
	if walkErr != nil {
		return nil, fmt.Errorf("scan incomplete: %w", walkErr)
	}

	// Phase 2: read tags in parallel using a worker pool.
	// Results are index-addressed so no mutex is needed.
	results := make([]AudioFile, len(entries))
	workers := runtime.NumCPU()
	if workers > 8 {
		workers = 8
	}

	ch := make(chan int, len(entries))
	for i := range entries {
		ch <- i
	}
	close(ch)

	var wg sync.WaitGroup
	for w := 0; w < workers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := range ch {
				e := entries[i]
				meta, _ := tags.ReadFile(e.path) // tag errors mean less metadata, not a failure
				title := meta.Title
				if title == "" {
					title = e.filenameTitle
				}
				results[i] = AudioFile{
					Title:       title,
					Artist:      meta.Artist,
					Album:       meta.Album,
					AlbumArtist: meta.AlbumArtist,
					Genre:       meta.Genre,
					Year:        meta.Year,
					TrackNum:    meta.TrackNum,
					TrackTotal:  meta.TrackTotal,
					DiscNum:     meta.DiscNum,
					DiscTotal:   meta.DiscTotal,
					HasArt:      meta.HasArt,
					FilePath:    e.path,
					Ext:         e.ext,
				}
			}
		}()
	}
	wg.Wait()

	return results, nil
}
