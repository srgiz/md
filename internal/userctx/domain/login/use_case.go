package login

import (
	"context"
	"errors"
)

var ErrIncorrectCredentials = errors.New("user: incorrect credentials")

type Command struct {
	Id       string `validate:"required"`
	Password string `validate:"required"`
}

type UseCase struct {
	tokenRepo TokenRepo
}

func NewUseCase(tokenRepo TokenRepo) *UseCase {
	return &UseCase{tokenRepo}
}

func (u *UseCase) Handle(ctx context.Context, cmd *Command) (*string, error) {
	token, err := u.tokenRepo.CreateToken(ctx, cmd.Id, cmd.Password)
	return &token, err
}

type TokenRepo interface {
	CreateToken(ctx context.Context, id string, password string) (string, error)
}
