package usecase

import (
	"context"
	"errors"
	"log"
	"os"
	"strings"
)

type EditFileCommand struct {
	Path  string `json:"path"`
	title string // todo
	body  string // todo
}

type EditFileUseCase struct {
	dataPath string // example: "./data/"
}

func NewEditFileUseCase(dataPath string) *EditFileUseCase {
	return &EditFileUseCase{dataPath: dataPath}
}

func (u *EditFileUseCase) Handle(ctx context.Context, cmd EditFileCommand) error {
	filePath := u.dataPath + strings.TrimLeft(cmd.Path, "/")
	log.Println("filePath: " + filePath)

	if _, err := os.Stat(filePath); errors.Is(err, os.ErrNotExist) {
		file, createErr := os.Create(filePath)
		if createErr != nil {
			return createErr
		}

		defer file.Close()

		writeErr := os.WriteFile(filePath, []byte(cmd.body), 0644)

		if writeErr != nil {
			return writeErr
		}
	} else if err != nil {
		return err
	} else {
		return errors.New("file already exists: todo")
	}

	return nil
}
