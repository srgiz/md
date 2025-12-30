package cmdbus

import (
	"context"
)

type Validator interface {
	Validate(ctx context.Context, s any) error
}

type validatorMiddleware struct {
	validator Validator
}

func NewValidatorMw(validator Validator) *validatorMiddleware {
	return &validatorMiddleware{validator}
}

func (m *validatorMiddleware) Handle(ctx context.Context, cmd any, next handler) (any, error) {
	if err := m.validator.Validate(ctx, cmd); err != nil {
		// todo: prepare text
		return nil, err
	}

	return next(ctx, cmd)
}
