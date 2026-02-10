package cmdbus

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandlerNotFound(t *testing.T) {
	_, err := NewBus().Handle(context.Background(), NewValidatorMiddleware())
	assert.ErrorIs(t, err, ErrHandlerNotFound)
}

func TestPanicAdd(t *testing.T) {
	tests := []struct {
		handler any
		msg     string
	}{
		{
			handler: 0,
			msg:     "handler must be a function",
		},
		{
			handler: func(context.Context) {},
			msg:     "handler must have two arguments",
		},
		{
			handler: func(int64, context.Context) {},
			msg:     "handler must have two return values",
		},
		{
			handler: func(int64, context.Context) (int64, error) {
				return 0, nil
			},
			msg: "first argument must be context.Context",
		},
		{
			handler: func(context.Context, int64) (int64, error) {
				return 0, nil
			},
			msg: "second argument must be a pointer",
		},
		{
			handler: func(context.Context, *int64) (int64, error) {
				return 0, nil
			},
			msg: "first return value must be a pointer or an interface",
		},
		{
			handler: func(context.Context, *int64) (*int64, int64) {
				return nil, 0
			},
			msg: "second return value must be an error",
		},
		{
			handler: func(context.Context, *int64) (*int64, error) {
				return nil, nil
			},
			msg: "handler already exists for *.int64",
		},
	}

	bus := NewBus()
	bus.Add(func(ctx context.Context, cmd *int64) (*int64, error) {
		return nil, errors.New("test error")
	})

	for i, test := range tests {
		t.Run(fmt.Sprintf("Example %d", i), func(t *testing.T) {
			assert.PanicsWithValue(t, test.msg, func() {
				bus.Add(test.handler)
			})
		})
	}
}
