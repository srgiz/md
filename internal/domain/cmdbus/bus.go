package cmdbus

import (
	"context"
	"fmt"
	"md/internal/domain/usecase/editfile"
	"md/internal/domain/usecase/findfile"
)

type handler func(ctx context.Context, cmd any) (any, error)

type Bus struct {
	validatorMw *validatorMiddleware
	handlers    map[string]handler
}

func New(
	validatorMw *validatorMiddleware,
	editFile *editfile.UseCase,
	findFile *findfile.UseCase,
) *Bus {
	bus := &Bus{
		validatorMw: validatorMw,
		handlers:    make(map[string]handler),
	}

	bus.add(&editfile.Command{}, bus.validatorMw.Handle(func(ctx context.Context, cmd any) (any, error) {
		return editFile.Handle(ctx, cmd.(*editfile.Command))
	}))

	bus.add(&findfile.Command{}, bus.validatorMw.Handle(func(ctx context.Context, cmd any) (any, error) {
		return findFile.Handle(ctx, cmd.(*findfile.Command)), nil
	}))

	return bus
}

func (bus *Bus) add(cmd any, h handler) {
	key := fmt.Sprintf("%T", cmd)

	if _, ok := bus.handlers[key]; ok {
		panic(fmt.Sprintf("duplicate key: %s", key))
	}

	bus.handlers[key] = h
}

func (bus *Bus) Handle(ctx context.Context, cmd any) (any, error) {
	h, ok := bus.handlers[fmt.Sprintf("%T", cmd)]

	if !ok {
		panic("handler not found")
	}

	return h(ctx, cmd)
}

/*
type Middleware interface {
	Handle(next handler) handler
}*/
