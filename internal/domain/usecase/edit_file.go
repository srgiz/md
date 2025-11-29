package usecase

import (
	"context"
	"md/internal/domain/entity"
	"md/internal/domain/repo"
)

type EditFileCommand struct {
	Path  string `json:"path" validate:"required"`
	Body  string `json:"body" validate:"required"`
	Title string `json:"title"`
}

type EditFileUseCase struct {
	fileRepo repo.FileRepository
}

func NewEditFileUseCase(fileRepo repo.FileRepository) *EditFileUseCase {
	return &EditFileUseCase{fileRepo: fileRepo}
}

func (u *EditFileUseCase) Handle(ctx context.Context, cmd *EditFileCommand) error {
	return u.fileRepo.Write(ctx, cmd.Path, &entity.File{
		Title: cmd.Title,
		Body:  cmd.Body,
	})
}
