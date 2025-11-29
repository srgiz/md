package cmdbus

import (
	"context"
	"fmt"
	"md/internal/domain/usecase"
	"md/internal/domain/validator"
)

type handler func(ctx context.Context, cmd any) (any, error)

type Bus struct {
	handlers            map[string]handler
	validatorMiddleware *validatorMiddleware
}

func New(
	validator validator.Validator,
	editFile *usecase.EditFileUseCase,
) *Bus {
	bus := &Bus{
		handlers:            make(map[string]handler),
		validatorMiddleware: &validatorMiddleware{validator: validator},
	}

	bus.add(&usecase.EditFileCommand{}, func(ctx context.Context, cmd any) (any, error) {
		return editFile.Handle(ctx, cmd.(*usecase.EditFileCommand))
	})

	return bus
}

func (bus *Bus) add(cmd any, h handler) {
	bus.handlers[fmt.Sprintf("%T", cmd)] = h
}

func (bus *Bus) Handle(ctx context.Context, cmd any) (any, error) {
	h, ok := bus.handlers[fmt.Sprintf("%T", cmd)]

	if !ok {
		panic("handler not found")
	}

	return bus.validatorMiddleware.Handle(ctx, cmd, h)
}

type Middleware interface {
	Handle(ctx context.Context, cmd any, next handler) (any, error)
}
