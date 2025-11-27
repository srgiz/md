package validator

import (
	"md/internal/domain/validator"

	playground_validator "github.com/go-playground/validator/v10"
)

type PlaygroundValidator struct {
	pgv *playground_validator.Validate
}

func NewPlaygroundValidator() validator.Validator {
	return &PlaygroundValidator{playground_validator.New()}
}

func (v *PlaygroundValidator) Validate(s any) error {
	return v.pgv.Struct(s)
}
