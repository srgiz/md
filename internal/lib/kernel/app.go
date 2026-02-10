package kernel

import (
	"context"
	"fmt"
	"log/slog"
	"md/internal/lib/kernel/cmdbus"
	"md/internal/lib/kernel/logger"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/urfave/cli/v3"
)

func init() {
	godotenv.Load(".env.local")
	godotenv.Load(".env")
}

func init() {
	slog.SetDefault(slog.New(logger.NewHandler(os.Stderr)))
}

type ServiceProvider func(di Di)

type Di interface {
	Bus() *cmdbus.Bus
	Service(name string) any
	AddService(name string, service any)
	AddCliCmd(cmd *cli.Command)
}

type App struct {
	services map[string]any
	bus      *cmdbus.Bus
	cmd      *cli.Command
}

func NewApp(services map[string]any) *App {
	return &App{
		services: services,
		bus:      cmdbus.NewBus(),
	}
}

func (app *App) Bus() *cmdbus.Bus {
	return app.bus
}

func (app *App) Service(name string) any {
	return app.services[name]
}

func (app *App) AddService(name string, service any) {
	if _, ok := app.services[name]; ok {
		panic(fmt.Sprintf("service %s already exists", name))
	}

	app.services[name] = service
}

func (app *App) Provide(fn ServiceProvider) {
	fn(app)
}

func (app *App) AddCliCmd(cmd *cli.Command) {
	if app.cmd == nil {
		app.cmd = &cli.Command{}
	}

	app.cmd.Commands = append(app.cmd.Commands, cmd)
}

func (app *App) RunCli(ctx context.Context, args []string) error {
	//traceId, _ := uuid.NewV7()
	//ctx := context.WithValue(ctx, logger.ContextKeyTraceId, traceId.String())
	slog.DebugContext(ctx, fmt.Sprintf("cli: %s", strings.Join(args[1:], " ")))
	err := app.cmd.Run(ctx, args)

	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("cli: %s", err.Error()))
	}

	return err
}
