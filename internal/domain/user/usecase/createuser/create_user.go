package createuser

import (
	"context"
	"md/internal/domain/user/repo"
)

type Command struct {
	Id       string `validate:"required"`
	Password string `validate:"required"`
}

type UseCase struct {
	userRepo repo.UserRepository
}

func New(userRepo repo.UserRepository) *UseCase {
	return &UseCase{userRepo}
}

func (u *UseCase) Handle(ctx context.Context, cmd *Command) error {
	return u.userRepo.Create(ctx, cmd.Id, cmd.Password)
}
