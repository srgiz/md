package validator

import (
	"md/internal/domain/validator"
	"strings"

	playground_validator "github.com/go-playground/validator/v10"
)

type PlaygroundValidator struct {
	pgv *playground_validator.Validate
}

func NewPlaygroundValidator() validator.Validator {
	pgv := playground_validator.New()

	if err := pgv.RegisterValidation("allowedDir", isAllowedDir); err != nil {
		panic(err)
	}

	return &PlaygroundValidator{pgv}
}

func (v *PlaygroundValidator) Validate(s any) error {
	return v.pgv.Struct(s)
}

func isAllowedDir(fl playground_validator.FieldLevel) bool {
	path := fl.Field().String()
	dirs := strings.Split(path, "/")

	for i := 0; i < len(dirs)-1; i++ {
		matches := strings.Split(dirs[i], ".")

		if len(matches) > 1 && matches[len(matches)-1] == "md" {
			return false
		}
	}

	return true
}
