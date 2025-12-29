package repo

import "context"

type UserRepository interface {
	Create(ctx context.Context, id string, password string) error
}
