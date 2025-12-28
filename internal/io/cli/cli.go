package cli

import (
	"context"
	"log"
	"md/internal/io/cli/internal"
	"os"

	"github.com/urfave/cli/v3"
)

type cliApp struct {
	cmd *cli.Command
}

func newCliApp(
	createUserCmd *internal.CreateUserCmd,
) *cliApp {
	app := &cliApp{&cli.Command{
		Commands: []*cli.Command{
			{
				Name:   "user:create",
				Action: createUserCmd.Handle,
			},
		},
	}}

	return app
}

func (app *cliApp) Run() {
	if err := app.cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
