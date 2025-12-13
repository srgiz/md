package internal

import (
	"md/internal/domain/cmdbus"
	"md/internal/domain/usecase/editfile"
	"md/internal/domain/usecase/findfile"
	"net/http"
)

type FileReceiver struct {
	bus *cmdbus.Bus
}

func NewFileReceiver(bus *cmdbus.Bus) *FileReceiver {
	return &FileReceiver{bus: bus}
}

// Edit executes "File.Edit"
func (s *FileReceiver) Edit(r *http.Request, cmd *editfile.Command, res *editfile.Result) error {
	return handle(s.bus, r, cmd, res)
}

// Find executes "File.Find"
func (s *FileReceiver) Find(r *http.Request, cmd *findfile.Command, res *findfile.Result) error {
	return handle(s.bus, r, cmd, res)
}
