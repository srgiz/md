package createuser

import (
	"context"
)

type Command struct {
	Id       string `validate:"required"`
	Password string `validate:"required"`
}

type UseCase struct {
	userRepo UserRepo
}

func NewUseCase(userRepo UserRepo) *UseCase {
	return &UseCase{userRepo}
}

func (u *UseCase) Handle(ctx context.Context, cmd *Command) error {
	return u.userRepo.Create(ctx, cmd.Id, cmd.Password)
}

type UserRepo interface {
	Create(ctx context.Context, id string, password string) error
}
