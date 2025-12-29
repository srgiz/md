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
				Action: func(ctx context.Context, command *cli.Command) error {
					_, err := bus.Handle(ctx, &createuser.Command{Id: command.Args().Get(0), Password: command.Args().Get(1)})
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
