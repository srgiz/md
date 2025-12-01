package editfile

import (
	"context"
	"md/internal/domain/entity"
	"md/internal/domain/repo"
)

type Command struct {
	Path  string `json:"path" validate:"required,filepath,allowedFilepath"`
	Body  string `json:"body" validate:"required"`
	Title string `json:"title"`
}

type UseCase struct {
	fileRepo repo.FileRepository
}

func New(fileRepo repo.FileRepository) *UseCase {
	return &UseCase{fileRepo: fileRepo}
}

func (u *UseCase) Handle(ctx context.Context, cmd *Command) (*Result, error) {
	file := &entity.File{
		Title: cmd.Title,
		Body:  cmd.Body,
	}

	if err := u.fileRepo.Write(ctx, cmd.Path, file); err != nil {
		return nil, err
	}

	return &Result{Id: file.Id}, nil
}

type Result struct {
	Id string `json:"id"`
}
