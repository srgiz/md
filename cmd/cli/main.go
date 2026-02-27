package main

import (
	"context"
	"fmt"
	"log/slog"
	"md/internal/lib/kernel"
	"md/internal/lib/logger"
	"md/internal/lib/sqldb"
	userdi "md/internal/userctx/present/di"
	"os"
)

func init() {
	slog.SetDefault(slog.New(logger.NewHandler(os.Stderr)))
}

func main() {
	app := newCliApp(map[string]any{
		"conn": sqldb.NewMasterSlaveConn("postgres", os.Getenv("GOOSE_DBSTRING"), os.Getenv("GOOSE_DBSTRING_SLAVE1")),
	})

	if err := app.RunCli(context.Background(), os.Args); err != nil {
		fmt.Printf("\x1b[%dm%s\x1b[0m\n", 31, err.Error())
		os.Exit(1)
	}
}

func newCliApp(services map[string]any) *kernel.App {
	conn, hasConn := services["conn"].(sqldb.Conn)

	if !hasConn {
		panic("service db.Conn not found")
	}

	app := kernel.NewApp(services)
	app.WrapCli(kernel.TraceConsoleMiddleware)
	// common infr:
	app.Provide(kernel.NewTxBusProvider("bus_tx", conn))
	// domain contexts:
	app.Provide(userdi.NewCliProvider())

	return app
}
