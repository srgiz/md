package kernel

import (
	"context"
	"fmt"
	"log/slog"

	playgroundValidator "github.com/go-playground/validator/v10"
)

var ValidatorBusMiddleware = &validatorMiddleware{playgroundValidator.New(playgroundValidator.WithRequiredStructEnabled())}

type validatorMiddleware struct {
	pgv *playgroundValidator.Validate
}

func (m *validatorMiddleware) Handle(ctx context.Context, cmd any, next cmdHandler) (any, error) {
	slog.DebugContext(ctx, fmt.Sprintf("bus: validate %T", cmd), "cmd", fmt.Sprintf("%#v", cmd))
	err := m.pgv.StructCtx(ctx, cmd)

	if err != nil {
		return nil, err
	}

	return next(ctx, cmd)
}

func (m *validatorMiddleware) RegisterValidation(tag string, fn playgroundValidator.Func) {
	if err := m.pgv.RegisterValidation(tag, fn); err != nil {
		panic(err)
	}
}

/*
// todo
func isAllowedFilepath(fl playgroundValidator.FieldLevel) bool {
	return !strings.Contains(fl.Field().String(), ".")
}*/
