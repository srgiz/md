package cmdbus

import (
	"context"
	"md/internal/domain/validator"
)

type validatorMiddleware struct {
	validator validator.Validator
}

func (m *validatorMiddleware) Handle(ctx context.Context, cmd any, next handler) (any, error) {
	if err := m.validator.Validate(cmd); err != nil {
		// todo: prepare text
		return nil, err
	}

	return next(ctx, cmd)
}
