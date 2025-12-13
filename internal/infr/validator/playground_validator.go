package validator

import (
	"md/internal/domain/cmdbus"
	"strings"

	playground_validator "github.com/go-playground/validator/v10"
)

type PlaygroundValidator struct {
	pgv *playground_validator.Validate
}

func NewPlaygroundValidator() cmdbus.Validator {
	pgv := playground_validator.New()

	if err := pgv.RegisterValidation("allowedFilepath", isAllowedFilepath); err != nil {
		panic(err)
	}

	return &PlaygroundValidator{pgv}
}

func (v *PlaygroundValidator) Validate(s any) error {
	return v.pgv.Struct(s)
}

func isAllowedFilepath(fl playground_validator.FieldLevel) bool {
	return !strings.Contains(fl.Field().String(), ".")
}
