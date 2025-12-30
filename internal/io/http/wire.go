//go:build wireinject

package http

import (
	"md/internal/infr"
	"md/internal/io/http/internal"

	"github.com/google/wire"
)

func Initialize(dataPath string) *server {
	wire.Build(
		newServer,
		infr.DiWireSet,
		internal.NewFileReceiver,
	)

	return &server{}
}
