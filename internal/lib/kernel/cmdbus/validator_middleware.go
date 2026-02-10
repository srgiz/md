package cmdbus

import (
	"context"
	"fmt"
	"log/slog"

	playgroundValidator "github.com/go-playground/validator/v10"
)

type ValidatorMiddleware struct {
	pgv *playgroundValidator.Validate
}

func NewValidatorMiddleware() *ValidatorMiddleware {
	return &ValidatorMiddleware{playgroundValidator.New(playgroundValidator.WithRequiredStructEnabled())}
}

func (m *ValidatorMiddleware) Handle(ctx context.Context, cmd any, next Handler) (any, error) {
	slog.DebugContext(ctx, fmt.Sprintf("bus: validate %T", cmd), "cmd", fmt.Sprintf("%#v", cmd))
	err := m.pgv.StructCtx(ctx, cmd)

	if err != nil {
		// todo. prepare text in controllers
		return nil, err
	}

	return next(ctx, cmd)
}

func (m *ValidatorMiddleware) RegisterValidation(tag string, fn playgroundValidator.Func) {
	if err := m.pgv.RegisterValidation(tag, fn); err != nil {
		panic(err)
	}
}

/*
// todo
func isAllowedFilepath(fl playgroundValidator.FieldLevel) bool {
	return !strings.Contains(fl.Field().String(), ".")
}*/
