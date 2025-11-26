package internal

import (
	"md/internal/domain/usecase"
	"net/http"

	"github.com/gorilla/rpc/v2/json2"
)

type FileReceiver struct {
	useCase *usecase.EditFileUseCase
}

func NewFileReceiver(useCase *usecase.EditFileUseCase) *FileReceiver {
	return &FileReceiver{useCase: useCase}
}

func (s *FileReceiver) Edit(r *http.Request, cmd *usecase.EditFileCommand, res *json2.EmptyResponse) error {
	return s.useCase.Handle(r.Context(), *cmd)
}
