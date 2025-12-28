package internal

import (
	"context"
	"fmt"
	"md/internal/domain/user"

	"github.com/urfave/cli/v3"
)

type CreateUserCmd struct {
	repo user.UserRepository
}

func NewCreateUserCmd(repo user.UserRepository) *CreateUserCmd {
	return &CreateUserCmd{repo}
}

func (c *CreateUserCmd) Handle(ctx context.Context, cmd *cli.Command) error {
	if cmd.Args().Len() < 2 {
		return fmt.Errorf("invalid number of arguments")
	}

	return c.repo.Create(ctx, cmd.Args().Get(0), cmd.Args().Get(1))
}
