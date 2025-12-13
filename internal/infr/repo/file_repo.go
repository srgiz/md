package repo

import (
	"context"
	"errors"
	"fmt"
	"log"
	"md/internal/domain/entity"
	"md/internal/domain/repo"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/goccy/go-yaml"
	"github.com/google/uuid"
)

type FileRepository struct {
	dataPath string // example: "./data/"
}

func NewFileRepository(dataPath string) repo.FileRepository {
	return &FileRepository{dataPath: dataPath}
}

func (r *FileRepository) fullPath(path string) string {
	return r.dataPath + strings.Trim(path, "/") + ".md"
}

func (r *FileRepository) Find(ctx context.Context, path string) *entity.File {
	fullPath := r.fullPath(path)
	log.Printf("fullPath: %s\n", fullPath)

	if _, err := os.Stat(fullPath); err != nil {
		return nil
	}

	headers, body, _ := r.parse(fullPath)
	log.Printf("find headers: %v\n", headers)

	return &entity.File{
		Body:  body,
		Id:    headers.Id,
		Title: headers.Title,
	}
}

func (r *FileRepository) Write(ctx context.Context, path string, file *entity.File) error {
	fullPath := r.fullPath(path)
	info, err := os.Stat(fullPath)

	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}

	if info != nil && info.IsDir() {
		return errors.New("file path is a directory")
	}

	if err == nil { // exists
		headers, _, headersErr := r.parse(fullPath)

		if headersErr != nil {
			log.Printf("err headers matches: %s\n", headersErr.Error())
		} else {
			file.Id = headers.Id
		}
	}

	if file.Id == "" {
		uuidV7, uuidErr := uuid.NewV7()

		if uuidErr != nil {
			return uuidErr
		}

		file.Id = uuidV7.String()
	}

	if file.Title == "" {
		file.Title = filepath.Base(path)
	}

	if err != nil { // not exists
		if mkdirErr := os.MkdirAll(filepath.Dir(fullPath), 0770); mkdirErr != nil {
			return mkdirErr // https://chmod-calculator.com/
		}

		osFile, createErr := os.Create(fullPath)
		if createErr != nil {
			return fmt.Errorf("create new file error: %w", createErr)
		}

		defer osFile.Close()
		return r.osWrite(osFile, file)
	} else { // exists
		osFile, openErr := os.OpenFile(fullPath, os.O_APPEND|os.O_WRONLY, 0660)
		if openErr != nil {
			return openErr
		}

		defer osFile.Close()
		return r.osWrite(osFile, file)
	}
}

func (r *FileRepository) osWrite(osFile *os.File, file *entity.File) error {
	yamlHeaders, yamlErr := yaml.Marshal(&fileHeaders{
		Id:    file.Id,
		Title: file.Title,
	})

	if yamlErr != nil {
		return yamlErr
	}

	truncateErr := osFile.Truncate(0)

	if truncateErr != nil {
		return truncateErr
	}

	osFile.WriteString("---\n")
	osFile.Write(yamlHeaders)
	osFile.WriteString("---\n")
	osFile.WriteString(file.Body)

	return nil
}

func (r *FileRepository) parse(fullPath string) (headers *fileHeaders, body string, err error) {
	re := regexp.MustCompile(`(?Us)^---\n(.*)?\n---\n(.+)?$`)
	content, _ := os.ReadFile(fullPath)
	matches := re.FindStringSubmatch(string(content))
	headers = &fileHeaders{}

	if len(matches) > 1 {
		err = yaml.Unmarshal([]byte(matches[1]), headers)
	}

	if len(matches) > 2 {
		body = matches[2]
	}

	return
}

type fileHeaders struct {
	Id    string `yaml:"id"`
	Title string `yaml:"title"`
}
