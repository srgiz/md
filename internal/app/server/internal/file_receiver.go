package internal

import (
	"md/internal/domain/cmdbus"
	"md/internal/domain/usecase"
	"net/http"
)

type FileReceiver struct {
	bus *cmdbus.Bus
}

func NewFileReceiver(bus *cmdbus.Bus) *FileReceiver {
	return &FileReceiver{bus: bus}
}

// Edit executes "File.Edit"
func (s *FileReceiver) Edit(r *http.Request, cmd *usecase.EditFileCommand, res *usecase.EditFileResult) error {
	return handle(s.bus, r, cmd, res)
}

// Find executes "File.Find"
func (s *FileReceiver) Find(r *http.Request, cmd *usecase.FindFileCommand, res *usecase.FindFileResult) error {
	return handle(s.bus, r, cmd, res)
}
