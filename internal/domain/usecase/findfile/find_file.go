package findfile

import (
	"context"
	"md/internal/domain/repo"
)

type Command struct {
	Path string `json:"path" validate:"required,filepath,allowedFilepath"`
}

type UseCase struct {
	fileRepo repo.FileRepository
}

func New(fileRepo repo.FileRepository) *UseCase {
	return &UseCase{fileRepo: fileRepo}
}

func (u *UseCase) Handle(ctx context.Context, cmd *Command) *Result {
	file := u.fileRepo.Find(ctx, cmd.Path)

	if file == nil {
		return &Result{}
	}

	return &Result{
		Id:    file.Id,
		Title: file.Title,
		Body:  file.Body,
	}
}

type Result struct {
	Id    string `json:"id"`
	Title string `json:"title"`
	Body  string `json:"body"`
}
