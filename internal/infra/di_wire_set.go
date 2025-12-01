package infra

import (
	"md/internal/domain/cmdbus"
	"md/internal/domain/usecase"
	infra_repo "md/internal/infra/repo"
	infra_validator "md/internal/infra/validator"

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

	infra_validator.NewPlaygroundValidator,
	//wire.NewSet(
	//	infra_validator.NewPlaygroundValidator,
	//wire.Bind(new(validator.Validator), new(*infra_validator.PlaygroundValidator)),
	//),

	infra_repo.NewFileRepository,

	cmdbus.New,
	usecase.NewEditFileUseCase,
	usecase.NewFindFileUseCase,
)
