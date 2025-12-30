package validator

import (
	"context"
	"fmt"
	"log/slog"
	"md/internal/domain/cmdbus"
	"strings"

	playgroundValidator "github.com/go-playground/validator/v10"
)

type PlaygroundValidator struct {
	pgv *playgroundValidator.Validate
}

func NewPlaygroundValidator() cmdbus.Validator {
	pgv := playgroundValidator.New()

	if err := pgv.RegisterValidation("allowedFilepath", isAllowedFilepath); err != nil {
		panic(err)
	}

	return &PlaygroundValidator{pgv}
}

func (v *PlaygroundValidator) Validate(ctx context.Context, s any) error {
	slog.DebugContext(ctx, fmt.Sprintf("pgv: validate %T", s), "struct", fmt.Sprintf("%#v", s))
	err := v.pgv.StructCtx(ctx, s)

	if err != nil {
		slog.DebugContext(ctx, fmt.Sprintf("pgv: %s", err.Error()))
	}

	return err
}

func isAllowedFilepath(fl playgroundValidator.FieldLevel) bool {
	return !strings.Contains(fl.Field().String(), ".")
}
