package cmdbus

import (
	"context"
	"fmt"
	"md/internal/domain/usecase"
)

type handler func(ctx context.Context, cmd any) (any, error)

type Bus struct {
	handlers map[string]handler
}

func New(
	editFile *usecase.EditFileUseCase,
) *Bus {
	bus := &Bus{handlers: make(map[string]handler)}

	bus.add(&usecase.EditFileCommand{}, func(ctx context.Context, cmd any) (any, error) {
		return nil, editFile.Handle(ctx, cmd.(*usecase.EditFileCommand))
	})

	return bus
}

func (bus *Bus) add(cmd any, h handler) {
	bus.handlers[fmt.Sprintf("%T", cmd)] = h
}

/*
func NewBus(args ...any) *Bus {
	if len(args)%2 != 0 {
		panic("args must be an even number")
	}

	handlers := make(map[string]handler)

	for i := 0; i < len(args); i += 2 {
		cmd, f := args[i], args[i+1]

		if reflect.ValueOf(cmd).Kind() != reflect.Ptr {
			panic("cmd must be a pointer")
		}

		h, ok := f.(handler)

		if !ok {
			panic("handler must be a handler func")
		}

		handlers[fmt.Sprintf("%T", cmd)] = h
	}

	return &Bus{handlers: handlers}
}*/

func (bus *Bus) Handle(ctx context.Context, cmd any) (any, error) {
	h, ok := bus.handlers[fmt.Sprintf("%T", cmd)]

	if !ok {
		panic("handler not found")
	}

	return h(ctx, cmd)
}
