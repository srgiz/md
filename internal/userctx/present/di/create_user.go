package di

import (
	"context"
	"md/internal/lib/kernel"
	"md/internal/userctx/domain/createuser"
	"md/internal/userctx/present/route"
)

func newCreateUserProvider(app *kernel.App, userRepo createuser.UserRepo) {
	app.AddCmdHandler(
		func(ctx context.Context, cmd *createuser.Command) (any, error) {
			return nil, createuser.NewUseCase(userRepo).Handle(ctx, cmd)
		},
		app.Service("bus_tx").(kernel.BusMiddleware),
		kernel.ValidatorBusMiddleware,
	)

	app.AddCliHandler(route.NewCreateUser(app))
}
