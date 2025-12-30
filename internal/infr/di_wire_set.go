//go:build wireinject

package infr

import (
	"md/internal/domain/cmdbus"
	"md/internal/domain/usecase/editfile"
	"md/internal/domain/usecase/findfile"
	"md/internal/domain/user/usecase/createuser"
	"md/internal/infr/postgres"
	infrRepo "md/internal/infr/repo"
	infrUser "md/internal/infr/user"
	infrValidator "md/internal/infr/validator"

	"github.com/google/wire"
	//"github.com/gorilla/schema"
)

var DiWireSet = wire.NewSet(
	//wire.NewSet(
	//	slog.New,
	//	logger.NewHandler,
	//	wire.Bind(new(slog.Handler), new(*logger.Handler)),
	//	wire.InterfaceValue(new(io.Writer), os.Stderr),
	//),

	//wire.NewSet(
	//	keydb.NewStream,
	//	keydb.NewClient,
	//wire.Value(keydb.MaxLen(2000)),
	//),
	//wire.Bind(new(messenger.Bus), new(*keydb.Stream)),

	//manticore.NewClient,
	//repositoryimpl.NewHttpClient,

	//schema.NewDecoder,

	postgres.NewConn,

	//infrValidator.NewPlaygroundValidator,
	//wire.NewSet(
	//	infrValidator.NewPlaygroundValidator,
	//wire.Bind(new(validator.Validator), new(*infrValidator.PlaygroundValidator)),
	//),

	// md
	wire.NewSet(
		infrRepo.NewFileRepository,
		editfile.New,
		findfile.New,
	),
	// user
	wire.NewSet(
		infrUser.NewUserRepository,
		createuser.New,
	),
	// cmdbus
	wire.NewSet(
		infrValidator.NewPlaygroundValidator,
		cmdbus.New,
		cmdbus.NewValidatorMw,
	),
)
