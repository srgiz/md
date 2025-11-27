package usecase

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type EditFileCommand struct {
	Path  string `json:"path"`
	title string // todo
	Body  string `json:"body"`
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
		if mkdirErr := os.MkdirAll(filepath.Dir(filePath), 0750); mkdirErr != nil {
			return err
		}

		file, createErr := os.Create(filePath)
		if createErr != nil {
			return fmt.Errorf("create new file error: %w", createErr)
		}

		defer file.Close()

		writeErr := os.WriteFile(filePath, []byte(cmd.Body), 0644)

		if writeErr != nil {
			return fmt.Errorf("create write file error: %w", writeErr)
		}
	} else if err != nil {
		return err
	} else {
		file, openErr := os.OpenFile(filePath /*os.O_APPEND|os.O_CREATE|*/, os.O_WRONLY, 0644)
		if openErr != nil {
			return openErr
		}

		defer file.Close()

		writeErr := os.WriteFile(filePath, []byte(cmd.Body), 0644)

		if writeErr != nil {
			return writeErr
		}
	}

	return nil
}
