package di

import (
	"md/internal/lib/kernel"
	"md/internal/lib/kernel/db"
	"md/internal/userctx/infr/repo"
)

func NewProviderCli() kernel.ServiceProvider {
	return func(di kernel.Di) {
		conn := di.Service("conn").(db.Conn)
		userRepo := repo.NewUserRepo(conn)
		newProviderCreateUser(di, userRepo)
	}
}
