package cli

import (
	"context"
	"log"
	"md/internal/domain/cmdbus"
	"md/internal/domain/user/usecase/createuser"
	"os"

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

func (app *cliApp) Run() {
	if err := app.cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
