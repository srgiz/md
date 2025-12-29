package cmdbus

import (
	"context"
	"fmt"
	"md/internal/domain/usecase/editfile"
	"md/internal/domain/usecase/findfile"
	"md/internal/domain/user/usecase/createuser"
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
	createUser *createuser.UseCase,
) *Bus {
	bus := &Bus{
		validatorMw: validatorMw,
		handlers:    make(map[string]handler),
	}

	bus.add(&editfile.Command{}, func(ctx context.Context, cmd any) (any, error) {
		return editFile.Handle(ctx, cmd.(*editfile.Command))
	}, bus.validatorMw)

	bus.add(&findfile.Command{}, func(ctx context.Context, cmd any) (any, error) {
		return findFile.Handle(ctx, cmd.(*findfile.Command)), nil
	}, bus.validatorMw)

	bus.add(&createuser.Command{}, func(ctx context.Context, cmd any) (any, error) {
		return nil, createUser.Handle(ctx, cmd.(*createuser.Command))
	}, bus.validatorMw)

	return bus
}

func (bus *Bus) add(cmd any, h handler, middlewares ...Middleware) {
	key := fmt.Sprintf("%T", cmd)

	if _, ok := bus.handlers[key]; ok {
		panic(fmt.Sprintf("duplicate key: %s", key))
	}

	for i := len(middlewares) - 1; i >= 0; i-- {
		h = func(m Middleware, next handler) handler {
			return func(ctx context.Context, cmd any) (any, error) {
				return m.Handle(ctx, cmd, next)
			}
		}(middlewares[i], h)
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

type Middleware interface {
	Handle(ctx context.Context, cmd any, next handler) (any, error)
}
