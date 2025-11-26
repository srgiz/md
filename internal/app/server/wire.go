//go:build wireinject
// +build wireinject

package server

import (
	"md/internal/app/server/internal"
	"md/internal/infra"

	"github.com/google/wire"
)

func Initialize(dataPath string) *server {
	wire.Build(
		newServer,
		infra.DiWireSet,
		internal.NewPingHandler,
		internal.NewFileReceiver,
	)

	return &server{}
}
