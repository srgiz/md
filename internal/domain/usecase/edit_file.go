package usecase

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/goccy/go-yaml"
	"github.com/google/uuid"
)

type EditFileCommand struct {
	Path  string `json:"path" validate:"required"`
	Body  string `json:"body" validate:"required"`
	Title string `json:"title"`
}

type EditFileUseCase struct {
	dataPath string // example: "./data/"
}

func NewEditFileUseCase(dataPath string) *EditFileUseCase {
	return &EditFileUseCase{dataPath: dataPath}
}

func (u *EditFileUseCase) Handle(ctx context.Context, cmd *EditFileCommand) error {
	filePath := u.dataPath + strings.TrimLeft(cmd.Path, "/")
	log.Println("filePath: " + filePath)

	if _, err := os.Stat(filePath); errors.Is(err, os.ErrNotExist) {
		if mkdirErr := os.MkdirAll(filepath.Dir(filePath), 0770); mkdirErr != nil {
			return err // https://chmod-calculator.com/
		}

		file, createErr := os.Create(filePath)
		if createErr != nil {
			return fmt.Errorf("create new file error: %w", createErr)
		}

		defer file.Close()

		//writeErr := os.WriteFile(filePath, []byte(cmd.Body), 0660)
		//_, writeErr := file.Write([]byte(cmd.Body))
		return u.writeFile(file, cmd, true)

		//if writeErr != nil {
		//	return fmt.Errorf("create write file error: %w", writeErr)
		//}
	} else if err != nil {
		return err
	} else {
		file, openErr := os.OpenFile(filePath /*os.O_APPEND|os.O_CREATE|*/, os.O_APPEND|os.O_WRONLY, 0660)
		if openErr != nil {
			return openErr
		}

		defer file.Close()

		truncateErr := file.Truncate(0)

		if truncateErr != nil {
			return truncateErr
		}

		//writeErr := os.WriteFile(filePath, []byte(cmd.Body), 0660)
		//_, writeErr := file.Write([]byte(cmd.Body))
		return u.writeFile(file, cmd, false)

		//if writeErr != nil {
		//	return writeErr
		//}
	}

	//return nil
}

func (u *EditFileUseCase) writeFile(file *os.File, cmd *EditFileCommand, isNew bool) error {
	title := cmd.Title

	if title == "" {
		title = filepath.Base(file.Name())
	}

	id := "" // todo: read id from file

	if isNew {
		uuidV7, err := uuid.NewV7()

		if err != nil {
			return err
		}

		id = uuidV7.String()
	}

	headers, err := yaml.Marshal(&fileHeaders{
		Id:    id,
		Title: title,
	})

	if err != nil {
		return err
	}

	file.WriteString("---\n")
	file.Write(headers)
	file.WriteString("---\n")
	file.WriteString(cmd.Body)

	//_, err := file.WriteString("---\n" /*+ string(headers)*/ + "---\n" + cmd.Body)
	//return err
	return nil
}

type fileHeaders struct {
	Id    string `yaml:"id"`
	Title string `yaml:"Title"`
}
