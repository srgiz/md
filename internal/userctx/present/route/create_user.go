package route

import (
	"context"
	"md/internal/lib/kernel"
	"md/internal/userctx/domain/createuser"

	"github.com/urfave/cli/v3"
)

func NewCreateUser(app *kernel.App) *cli.Command {
	return &cli.Command{
		Name: "user:create",
		Arguments: []cli.Argument{
			&cli.StringArg{
				Name: "id",
			},
			&cli.StringArg{
				Name: "password",
			},
		},
		Action: func(ctx context.Context, command *cli.Command) error {
			_, err := app.Handle(ctx, &createuser.Command{
				Id:       command.StringArg("id"),
				Password: command.StringArg("password"),
			})

			return err
		},
	}
}
