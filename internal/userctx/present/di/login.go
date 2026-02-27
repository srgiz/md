package di

import (
	"md/internal/lib/kernel"
	"md/internal/userctx/domain/login"
	"md/internal/userctx/present/route"
)

func newLoginProvider(app *kernel.App, tokenRepo login.TokenRepo) {
	app.AddCmdHandler(
		login.NewUseCase(tokenRepo).Handle,
		kernel.ValidatorBusMiddleware,
	)

	app.AddHttpHandler(route.NewPostLogin(app))
}
