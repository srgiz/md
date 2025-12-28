//go:build wireinject
// +build wireinject

package cli

import (
	"md/internal/infr"
	"md/internal/io/cli/internal"

	"github.com/google/wire"
)

func Initialize( /*dataPath string*/ ) *cliApp {
	wire.Build(
		newCliApp,
		infr.DiWireSet,
		internal.NewCreateUserCmd,
	)

	return &cliApp{}
}
