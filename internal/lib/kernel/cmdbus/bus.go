package cmdbus

import (
	"context"
	"errors"
	"fmt"
	"reflect"
)

var ErrHandlerNotFound = errors.New("bus: handler not found")

type Handler func(ctx context.Context, cmd any) (any, error)

type Bus struct {
	handlers map[string]Handler
}

func NewBus() *Bus {
	return &Bus{
		handlers: make(map[string]Handler),
	}
}

// Add Последний middleware является самым внешним
// Пример handler: func(ctx context.Context, cmd *Command) (*Reply, error)
func (bus *Bus) Add(handler any, middlewares ...Middleware) {
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

	key := bus.key(t.In(1))

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
		h = func(m Middleware, next Handler) Handler {
			return func(ctx context.Context, cmd any) (any, error) {
				return m.Handle(ctx, cmd, next)
			}
		}(middlewares[i], h)
	}

	bus.handlers[key] = h
}

func (bus *Bus) Handle(ctx context.Context, cmd any) (any, error) {
	h, ok := bus.handlers[bus.key(reflect.TypeOf(cmd))]

	if !ok {
		return nil, ErrHandlerNotFound
	}

	return h(ctx, cmd)
}

func (bus *Bus) key(cmdType reflect.Type) string {
	return fmt.Sprintf("*%s.%s", cmdType.Elem().PkgPath(), cmdType.Elem().Name())
}

type Middleware interface {
	Handle(ctx context.Context, cmd any, next Handler) (any, error)
}
