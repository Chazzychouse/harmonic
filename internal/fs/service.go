package fs

import "context"

// FsService exposes filesystem operations to the Wails frontend.
type FsService struct {
	ctx context.Context
}

func NewFsService() *FsService {
	return &FsService{}
}

// SetContext stores the Wails runtime context for dialog APIs.
func (s *FsService) SetContext(ctx context.Context) {
	s.ctx = ctx
}
