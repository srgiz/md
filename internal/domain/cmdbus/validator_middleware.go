package cmdbus

import (
	"context"
	"log"
	"md/internal/domain/validator"
)

type validatorMiddleware struct {
	validator validator.Validator
}

func (m *validatorMiddleware) Handle(ctx context.Context, cmd any, next handler) (any, error) {
	if err := m.validator.Validate(cmd); err != nil {
		log.Printf("Validator error: %s", err)
		// todo: prepare text
		return nil, err
	}

	return next(ctx, cmd)
}
