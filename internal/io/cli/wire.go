//go:build wireinject

package cli

import (
	"md/internal/infr"

	"github.com/google/wire"
)

func Initialize(dataPath string) *cliApp {
	wire.Build(
		newCliApp,
		infr.DiWireSet,
	)

	return &cliApp{}
}
