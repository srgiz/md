package repo

import (
	"context"
	"md/internal/domain/entity"
)

type FileRepository interface {
	Write(ctx context.Context, path string, file *entity.File) error
}
