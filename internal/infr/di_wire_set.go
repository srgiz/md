package infr

import (
	"md/internal/domain/cmdbus"
	"md/internal/domain/usecase/editfile"
	"md/internal/domain/usecase/findfile"
	infrRepo "md/internal/infr/repo"
	infrValidator "md/internal/infr/validator"

	"github.com/google/wire"
	//"github.com/gorilla/schema"
)

var Version = "" // ldflags

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

	infrValidator.NewPlaygroundValidator,
	//wire.NewSet(
	//	infrValidator.NewPlaygroundValidator,
	//wire.Bind(new(validator.Validator), new(*infrValidator.PlaygroundValidator)),
	//),

	infrRepo.NewFileRepository,

	wire.NewSet(
		cmdbus.New,
		cmdbus.NewValidatorMw,
		editfile.New,
		findfile.New,
	),
)
