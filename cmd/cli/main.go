package main

import (
	"context"
	"md/internal/lib/kernel"
	"md/internal/lib/kernel/db"
	userdi "md/internal/userctx/infr/di"
	"os"
)

func main() {
	app := newCliApp(map[string]any{
		"conn": db.NewMasterSlaveConn("postgres", os.Getenv("GOOSE_DBSTRING"), os.Getenv("GOOSE_DBSTRING_SLAVE1")),
	})

	if err := app.RunCli(context.Background(), os.Args); err != nil {
		os.Exit(1)
	}
}

func newCliApp(services map[string]any) *kernel.App {
	conn, hasConn := services["conn"].(db.Conn)

	if !hasConn {
		panic("service db.Conn not found")
	}

	app := kernel.NewApp(services)
	// common infr:
	app.Provide(kernel.NewProviderValidatorMiddleware())
	app.Provide(kernel.NewProviderTransactionMiddleware(conn))
	app.Provide(kernel.NewProviderTraceMiddleware())
	// domain contexts:
	app.Provide(userdi.NewProviderCli())

	return app
}
