package cmdbus

import (
	"context"
	"log"
)

type Validator interface {
	Validate(s any) error
}

type validatorMiddleware struct {
	validator Validator
}

func NewValidatorMw(validator Validator) *validatorMiddleware {
	return &validatorMiddleware{validator}
}

func (m *validatorMiddleware) Handle(ctx context.Context, cmd any, next handler) (any, error) {
	if err := m.validator.Validate(cmd); err != nil {
		log.Printf("Validator error: %s", err)
		// todo: prepare text
		return nil, err
	}

	return next(ctx, cmd)
}
