package kernel

import (
	"context"
	"errors"
	"fmt"
	"reflect"
)

var ErrCommandHandlerNotFound = errors.New("app: command handler not found")

type cmdHandler func(ctx context.Context, cmd any) (any, error)

type cmdBus struct {
	handlers map[string]cmdHandler
}

// AddCmdHandler Последний middleware является самым внешним
// Пример handler: func(ctx context.Context, console *Command) (*Reply, error)
func (bus *cmdBus) AddCmdHandler(handler any, middlewares ...BusMiddleware) {
	t := reflect.TypeOf(handler)

	if t.Kind() != reflect.Func {
		panic("handler must be a function")
	}

	if t.NumIn() != 2 {
		panic("handler must have two arguments")
	}

	if t.NumOut() != 2 {
		panic("handler must have two return values")
	}

	if !t.In(0).Implements(reflect.TypeOf((*context.Context)(nil)).Elem()) {
		panic("first argument must be context.Context")
	}

	if t.In(1).Kind() != reflect.Pointer {
		panic("second argument must be a pointer")
	}

	if t.Out(0).Kind() != reflect.Pointer && t.Out(0).Kind() != reflect.Interface {
		panic("first return value must be a pointer or an interface")
	}

	if !t.Out(1).Implements(reflect.TypeOf((*error)(nil)).Elem()) {
		panic("second return value must be an error")
	}

	key := bus.keyCmd(t.In(1))

	if _, ok := bus.handlers[key]; ok {
		panic(fmt.Sprintf("handler already exists for %s", key))
	}

	h := func(ctx context.Context, cmd any) (any, error) {
		vals := reflect.ValueOf(handler).Call([]reflect.Value{reflect.ValueOf(ctx), reflect.ValueOf(cmd)})

		var err error

		if vals[1].Interface() != nil {
			err = vals[1].Interface().(error)
		}

		return vals[0].Interface(), err
	}

	for i := 0; i < len(middlewares); i++ {
		h = func(m BusMiddleware, next cmdHandler) cmdHandler {
			return func(ctx context.Context, cmd any) (any, error) {
				return m.Handle(ctx, cmd, next)
			}
		}(middlewares[i], h)
	}

	bus.handlers[key] = h
}

func (bus *cmdBus) Handle(ctx context.Context, cmd any) (any, error) {
	h, ok := bus.handlers[bus.keyCmd(reflect.TypeOf(cmd))]

	if !ok {
		return nil, ErrCommandHandlerNotFound
	}

	return h(ctx, cmd)
}

func (bus *cmdBus) keyCmd(cmdType reflect.Type) string {
	return fmt.Sprintf("*%s.%s", cmdType.Elem().PkgPath(), cmdType.Elem().Name())
}

type BusMiddleware interface {
	Handle(ctx context.Context, cmd any, next cmdHandler) (any, error)
}
