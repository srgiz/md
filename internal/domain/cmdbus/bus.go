package cmdbus

import (
	"context"
	"fmt"
	"log"
	"md/internal/domain/usecase/editfile"
	"md/internal/domain/usecase/findfile"
)

type handler func(ctx context.Context, cmd any) (any, error)

type Bus struct {
	handlers  map[string]handler
	validator Validator
}

func New(
	validator Validator,
	editFile *editfile.UseCase,
	findFile *findfile.UseCase,
) *Bus {
	bus := &Bus{
		validator: validator,
		handlers:  make(map[string]handler),
	}

	bus.add(&editfile.Command{}, func(ctx context.Context, cmd any) (any, error) {
		return editFile.Handle(ctx, cmd.(*editfile.Command))
	})

	bus.add(&findfile.Command{}, func(ctx context.Context, cmd any) (any, error) {
		return findFile.Handle(ctx, cmd.(*findfile.Command)), nil
	})

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

	return bus.validate(ctx, cmd, h)
}

func (bus *Bus) validate(ctx context.Context, cmd any, next handler) (any, error) {
	if err := bus.validator.Validate(cmd); err != nil {
		log.Printf("Validator error: %s", err)
		// todo: prepare text
		return nil, err
	}

	return next(ctx, cmd)
}

type Middleware interface {
	Handle(ctx context.Context, cmd any, next handler) (any, error)
}
