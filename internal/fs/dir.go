package fs

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// AudioFile holds metadata about a discovered audio file.
type AudioFile struct {
	Title    string `json:"title"`
	FilePath string `json:"file_path"`
	Ext      string `json:"ext"`
}

var audioExtensions = map[string]bool{
	".mp3": true, ".flac": true, ".wav": true, ".aiff": true,
	".aif": true, ".m4a": true, ".ogg": true, ".opus": true,
	".wv": true, ".alac": true, ".aac": true, ".wma": true,
}

// SelectDirectory opens a native OS directory picker and returns the chosen path.
// Returns an empty string if the user cancels.
func (s *FsService) SelectDirectory() string {
	dir, err := runtime.OpenDirectoryDialog(s.ctx, runtime.OpenDialogOptions{
		Title: "Select Music Library Folder",
	})
	if err != nil || dir == "" {
		return ""
	}
	return dir
}

// Slow af
func (s *FsService) ScanDirectory(dir string) ([]AudioFile, error) {
	if dir == "" {
		return nil, fmt.Errorf("no directory specified")
	}
	if _, err := os.Stat(dir); err != nil {
		return nil, fmt.Errorf("cannot access %q: %w", dir, err)
	}

	var files []AudioFile
	err := filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
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
		name := strings.TrimSuffix(d.Name(), rawExt)
		files = append(files, AudioFile{
			Title:    name,
			FilePath: path,
			Ext:      ext[1:],
		})
		return nil
	})
	if err != nil {
		return files, fmt.Errorf("scan incomplete: %w", err)
	}
	return files, nil
}
