package main

import (
	"fmt"
	"log/slog"
	"md/internal/lib/kernel"
	"md/internal/lib/logger"
	"md/internal/lib/sqldb"
	userdi "md/internal/userctx/present/di"
	"net/http"
	"os"
)

func init() {
	slog.SetDefault(slog.New(logger.NewHandler(os.Stderr)))
}

func main() {
	app := newHttpApp(map[string]any{
		"conn": sqldb.NewMasterSlaveConn("postgres", os.Getenv("GOOSE_DBSTRING"), os.Getenv("GOOSE_DBSTRING_SLAVE1")),
	})

	// todo. del
	app.AddHttpHandler("GET /{$}", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(fmt.Sprintf("%s %s %v", r.Method, r.URL.Path, r.Context().Value("auth.Token"))))
	}))

	if err := app.RunHttp(); err != nil {
		os.Exit(1)
	}
}

func newHttpApp(services map[string]any) *kernel.App {
	conn, hasConn := services["conn"].(sqldb.Conn)

	if !hasConn {
		panic("service db.Conn not found")
	}

	app := kernel.NewApp(services)
	app.WrapHttp(kernel.JwtHttpMiddleware, kernel.TraceHttpMiddleware)
	// common infr:
	app.Provide(kernel.NewTxBusProvider("bus_tx", conn))
	app.Provide(kernel.NewNotFoundHttpProvider())
	// domain contexts:
	app.Provide(userdi.NewHttpProvider())

	return app
}
