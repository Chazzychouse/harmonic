package tags

import (
	"testing"
)

func TestReadFile_TaggedMP3(t *testing.T) {
	meta, err := ReadFile("testdata/tagged.mp3")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if meta.Title != "Test Title" {
		t.Errorf("title = %q, want %q", meta.Title, "Test Title")
	}
	if meta.Artist != "Test Artist" {
		t.Errorf("artist = %q, want %q", meta.Artist, "Test Artist")
	}
	if meta.Album != "Test Album" {
		t.Errorf("album = %q, want %q", meta.Album, "Test Album")
	}
	if !meta.HasArt {
		t.Error("expected HasArt=true for tagged.mp3")
	}
}

func TestReadFile_TaggedFLAC(t *testing.T) {
	meta, err := ReadFile("testdata/tagged.flac")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if meta.Title != "FLAC Title" {
		t.Errorf("title = %q, want %q", meta.Title, "FLAC Title")
	}
	if meta.Artist != "FLAC Artist" {
		t.Errorf("artist = %q, want %q", meta.Artist, "FLAC Artist")
	}
	if meta.Album != "FLAC Album" {
		t.Errorf("album = %q, want %q", meta.Album, "FLAC Album")
	}
	if meta.HasArt {
		t.Error("expected HasArt=false for tagged.flac")
	}
}

func TestReadFile_TaglessWAV(t *testing.T) {
	meta, err := ReadFile("testdata/tagless.wav")
	if err != nil {
		t.Fatalf("tagless file should not return error, got: %v", err)
	}
	if meta.Title != "" || meta.Artist != "" || meta.Album != "" {
		t.Errorf("expected zero-value meta for tagless file, got %+v", meta)
	}
}

func TestReadFile_NonexistentFile(t *testing.T) {
	_, err := ReadFile("testdata/does_not_exist.mp3")
	if err == nil {
		t.Error("expected error for nonexistent file, got nil")
	}
}

func TestReadArt_WithArt(t *testing.T) {
	data, mime, err := ReadArt("testdata/tagged.mp3")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(data) == 0 {
		t.Error("expected non-empty image data")
	}
	if mime == "" {
		t.Error("expected non-empty MIME type")
	}
}

func TestReadArt_WithoutArt(t *testing.T) {
	_, _, err := ReadArt("testdata/tagged.flac")
	if err == nil {
		t.Error("expected error for file without art, got nil")
	}
}
