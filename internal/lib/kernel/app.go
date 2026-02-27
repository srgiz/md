package kernel

import (
	"context"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/urfave/cli/v3"
)

func init() {
	godotenv.Load(".env.local")
	godotenv.Load(".env")
}

type App struct {
	di
	cmdBus
	console            *cli.Command
	consoleMiddlewares []func(consoleHandler) consoleHandler
	mux                *http.ServeMux
	muxMiddlewares     []func(http.Handler) http.Handler
}

type consoleHandler func(context.Context, []string) error

func NewApp(services map[string]any) *App {
	return &App{
		di:     di{services: services},
		cmdBus: cmdBus{handlers: make(map[string]cmdHandler)},
	}
}

func (app *App) Provide(fn func(app *App)) {
	fn(app)
}

// WrapCli Последний middleware является самым внешним
func (app *App) WrapCli(middlewares ...func(consoleHandler) consoleHandler) {
	app.consoleMiddlewares = append(app.consoleMiddlewares, middlewares...)
}

func (app *App) AddCliHandler(cmd *cli.Command) {
	if app.console == nil {
		app.console = &cli.Command{}
	}

	app.console.Commands = append(app.console.Commands, cmd)
}

func (app *App) RunCli(ctx context.Context, args []string) error {
	var handler consoleHandler = app.console.Run

	for _, m := range app.consoleMiddlewares {
		handler = m(handler)
	}

	return handler(ctx, args)
}

// WrapHttp Последний middleware является самым внешним
func (app *App) WrapHttp(middlewares ...func(http.Handler) http.Handler) {
	app.muxMiddlewares = append(app.muxMiddlewares, middlewares...)
}

func (app *App) AddHttpHandler(pattern string, handler http.Handler) {
	if app.mux == nil {
		app.mux = http.NewServeMux()
	}

	app.mux.Handle(pattern, handler)
}

func (app *App) RunHttp() error {
	var handler http.Handler = app.mux

	for _, m := range app.muxMiddlewares {
		handler = m(handler)
	}

	return http.ListenAndServe(":8080", handler)
}

func NewNotFoundHttpProvider() func(app *App) {
	return func(app *App) {
		app.AddHttpHandler("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}))
	}
}
