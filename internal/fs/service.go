package fs

import "context"

// FsService exposes filesystem operations to the Wails frontend.
type FsService struct {
	ctx       context.Context
	mediaPort int
}

func NewFsService(mediaPort int) *FsService {
	return &FsService{mediaPort: mediaPort}
}

// SetContext stores the Wails runtime context for dialog APIs.
func (s *FsService) SetContext(ctx context.Context) {
	s.ctx = ctx
}

// MediaPort returns the port of the dedicated media HTTP server.
// The frontend uses http://127.0.0.1:{port} to stream audio and art.
func (s *FsService) MediaPort() int {
	return s.mediaPort
}
