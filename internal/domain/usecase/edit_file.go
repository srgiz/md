package usecase

import (
	"context"
	"md/internal/domain/entity"
	"md/internal/domain/repo"
)

type EditFileCommand struct {
	Path  string `json:"path" validate:"required,filepath,allowedFilepath"`
	Body  string `json:"body" validate:"required"`
	Title string `json:"title"`
}

type EditFileUseCase struct {
	fileRepo repo.FileRepository
}

func NewEditFileUseCase(fileRepo repo.FileRepository) *EditFileUseCase {
	return &EditFileUseCase{fileRepo: fileRepo}
}

func (u *EditFileUseCase) Handle(ctx context.Context, cmd *EditFileCommand) (*EditFileResult, error) {
	file := &entity.File{
		Title: cmd.Title,
		Body:  cmd.Body,
	}

	if err := u.fileRepo.Write(ctx, cmd.Path, file); err != nil {
		return nil, err
	}

	return &EditFileResult{Id: file.Id}, nil
}

type EditFileResult struct {
	Id string `json:"id"`
}
