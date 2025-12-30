package validator

import (
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

func (v *PlaygroundValidator) Validate(s any) error {
	slog.Debug(fmt.Sprintf("validator: %v", s))
	return v.pgv.Struct(s)
}

func isAllowedFilepath(fl playgroundValidator.FieldLevel) bool {
	return !strings.Contains(fl.Field().String(), ".")
}
