package cli

import (
	"context"
	"fmt"
	"log/slog"
	"md/internal/domain/cmdbus"
	"md/internal/domain/user/usecase/createuser"
	"strings"

	"github.com/google/uuid"
	"github.com/urfave/cli/v3"
)

type cliApp struct {
	cmd *cli.Command
}

func newCliApp(
	bus *cmdbus.Bus,
) *cliApp {
	return &cliApp{&cli.Command{
		Commands: []*cli.Command{
			{
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
					_, err := bus.Handle(ctx, &createuser.Command{
						Id:       command.StringArg("id"),
						Password: command.StringArg("password"),
					})

					return err
				},
			},
		},
	}}
}

func (app *cliApp) Run(args []string) error {
	requestId, _ := uuid.NewV7()
	ctx := context.WithValue(context.Background(), "X-Request-ID", requestId.String())

	slog.DebugContext(ctx, fmt.Sprintf("cli: %s", strings.Join(args[1:], " ")))
	err := app.cmd.Run(ctx, args)

	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("cli: %s", err.Error()))
	}

	return err
}
