package internal

import (
	"md/internal/domain/cmdbus"
	"md/internal/domain/usecase"
	"net/http"

	"github.com/gorilla/rpc/v2/json2"
)

type FileReceiver struct {
	bus *cmdbus.Bus
}

func NewFileReceiver(bus *cmdbus.Bus) *FileReceiver {
	return &FileReceiver{bus: bus}
}

func (s *FileReceiver) Edit(r *http.Request, cmd *usecase.EditFileCommand, res *json2.EmptyResponse) error {
	_, err := s.bus.Handle(r.Context(), cmd)
	return err
}
