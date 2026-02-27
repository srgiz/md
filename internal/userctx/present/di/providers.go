package di

import (
	"md/internal/lib/kernel"
	"md/internal/lib/sqldb"
	"md/internal/userctx/infr/repo"
)

func NewCliProvider() func(app *kernel.App) {
	return func(app *kernel.App) {
		userRepo := repo.NewUserRepo(app.Service("conn").(sqldb.Conn))
		newCreateUserProvider(app, userRepo)
	}
}

func NewHttpProvider() func(app *kernel.App) {
	return func(app *kernel.App) {
		userRepo := repo.NewUserRepo(app.Service("conn").(sqldb.Conn))
		newLoginProvider(app, userRepo)
	}
}
