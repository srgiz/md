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

func (s *FileReceiver) Edit(r *http.Request, cmd *usecase.EditFileCommand, res *usecase.EditFileResult) error {
	data, err := s.bus.Handle(r.Context(), cmd)
	*res = *data.(*usecase.EditFileResult)
	return err
}
