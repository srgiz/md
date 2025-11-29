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
	dto, err := s.bus.Handle(r.Context(), cmd)

	if dto != nil {
		if data := dto.(*usecase.EditFileResult); data != nil {
			*res = *data
		}
	}

	return err
}
