package di

import (
	"context"
	"md/internal/lib/kernel"
	"md/internal/lib/kernel/cmdbus"
	"md/internal/userctx/domain/createuser"

	"github.com/urfave/cli/v3"
)

func newProviderCreateUser(di kernel.Di, userRepo createuser.UserRepo) {
	di.Bus().Add(
		func(ctx context.Context, cmd *createuser.Command) (any, error) {
			return nil, createuser.NewUseCase(userRepo).Handle(ctx, cmd)
		},
		di.Service(kernel.ServiceNameTransactionMiddleware).(cmdbus.Middleware),
		di.Service(kernel.ServiceNameValidatorMiddleware).(cmdbus.Middleware),
		di.Service(kernel.ServiceNameTraceMiddleware).(cmdbus.Middleware),
	)

	di.AddCliCmd(&cli.Command{
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
			_, err := di.Bus().Handle(ctx, &createuser.Command{
				Id:       command.StringArg("id"),
				Password: command.StringArg("password"),
			})

			return err
		},
	})
}
