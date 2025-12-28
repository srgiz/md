package user

import "context"

type UserRepository interface {
	Create(ctx context.Context, id string, pwHash string) error
}
