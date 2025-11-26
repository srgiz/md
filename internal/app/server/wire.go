//go:build wireinject
// +build wireinject

package server

import (
	"md/internal/app/server/internal"

	"github.com/google/wire"
)

func Initialize() *server {
	wire.Build(
		newServer,
		//di.Set,
		internal.NewPingHandler,
	)

	return &server{}
}
