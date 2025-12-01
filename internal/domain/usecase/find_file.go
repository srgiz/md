package usecase

import (
	"context"
	"md/internal/domain/repo"
)

type FindFileCommand struct {
	Path string `json:"path" validate:"required,filepath,allowedFilepath"`
}

type FindFileUseCase struct {
	fileRepo repo.FileRepository
}

func NewFindFileUseCase(fileRepo repo.FileRepository) *FindFileUseCase {
	return &FindFileUseCase{fileRepo: fileRepo}
}

func (u *FindFileUseCase) Handle(ctx context.Context, cmd *FindFileCommand) *FindFileResult {
	file := u.fileRepo.Find(ctx, cmd.Path)

	if file == nil {
		return &FindFileResult{}
	}

	return &FindFileResult{
		Id:    file.Id,
		Title: file.Title,
		Body:  file.Body,
	}
}

type FindFileResult struct {
	Id    string `json:"id"`
	Title string `json:"title"`
	Body  string `json:"body"`
}
