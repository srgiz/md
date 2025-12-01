package repo

import (
	"context"
	"md/internal/domain/entity"
)

type FileRepository interface {
	Find(ctx context.Context, path string) *entity.File
	Write(ctx context.Context, path string, file *entity.File) error
}
