package tags

import (
	"errors"
	"os"

	"github.com/dhowden/tag"
)

// TrackMeta holds audio metadata extracted from embedded tags.
type TrackMeta struct {
	Title       string `json:"title"`
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
}

// ReadFile reads embedded tags from the audio file at path.
// If the file has no tags, a zero-value TrackMeta is returned with no error.
func ReadFile(path string) (TrackMeta, error) {
	f, err := os.Open(path)
	if err != nil {
		return TrackMeta{}, err
	}
	defer f.Close()

	m, err := tag.ReadFrom(f)
	if err != nil {
		if errors.Is(err, tag.ErrNoTagsFound) {
			return TrackMeta{}, nil
		}
		return TrackMeta{}, err
	}

	trackNum, trackTotal := m.Track()
	discNum, discTotal := m.Disc()

	return TrackMeta{
		Title:       m.Title(),
		Artist:      m.Artist(),
		Album:       m.Album(),
		AlbumArtist: m.AlbumArtist(),
		Genre:       m.Genre(),
		Year:        m.Year(),
		TrackNum:    trackNum,
		TrackTotal:  trackTotal,
		DiscNum:     discNum,
		DiscTotal:   discTotal,
		HasArt:      m.Picture() != nil,
	}, nil
}

// ReadArt extracts the embedded album art from the audio file at path.
// Returns the raw image bytes and MIME type, or an error if no art is found.
func ReadArt(path string) ([]byte, string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, "", err
	}
	defer f.Close()

	m, err := tag.ReadFrom(f)
	if err != nil {
		return nil, "", err
	}

	pic := m.Picture()
	if pic == nil {
		return nil, "", errors.New("no art found")
	}

	return pic.Data, pic.MIMEType, nil
}
