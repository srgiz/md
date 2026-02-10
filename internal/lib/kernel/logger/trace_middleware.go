package logger

import (
	"context"
	"md/internal/lib/kernel/cmdbus"

	"github.com/google/uuid"
)

type traceMiddleware struct {
}

func NewTraceMiddleware() cmdbus.Middleware {
	return &traceMiddleware{}
}

func (m *traceMiddleware) Handle(ctx context.Context, cmd any, next cmdbus.Handler) (any, error) {
	if _, ok := ctx.Value(ContextKeyTraceId).(string); !ok {
		traceId, _ := uuid.NewV7()
		ctx = context.WithValue(ctx, ContextKeyTraceId, traceId.String())
	}

	return next(ctx, cmd)
}
