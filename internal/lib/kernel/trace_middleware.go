package kernel

import (
	"context"
	"fmt"
	"log/slog"
	"md/internal/lib/logger"
	"net/http"

	"github.com/google/uuid"
)

var TraceHttpMiddleware = func(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		*r = *r.WithContext(newTraceCtx(r.Context()))
		writer := &logResponseWriter{http.StatusInternalServerError, w}

		// изначально r.Pattern не заполнен
		defer func() {
			slog.DebugContext(r.Context(), fmt.Sprintf(`http: "%s" %d`, r.Pattern, writer.StatusCode), "method", r.Method, "uri", r.RequestURI)
		}()

		next.ServeHTTP(writer, r)
	})
}

var TraceConsoleMiddleware = func(next consoleHandler) consoleHandler {
	return func(ctx context.Context, args []string) error {
		return next(newTraceCtx(ctx), args)
	}
}

/*
var TraceBusMiddleware = &traceBusMiddleware{}

type traceBusMiddleware struct {
}

func (m *traceBusMiddleware) Handle(ctx context.Context, cmd any, next cmdHandler) (any, error) {
	return next(newTraceCtx(ctx), cmd)
}*/

func newTraceCtx(ctx context.Context) context.Context {
	if _, ok := ctx.Value(logger.ContextKeyTraceId).(string); !ok {
		traceId, _ := uuid.NewV7()
		return context.WithValue(ctx, logger.ContextKeyTraceId, traceId.String())
	}

	return ctx
}

type logResponseWriter struct {
	StatusCode     int
	responseWriter http.ResponseWriter
}

func (w *logResponseWriter) Header() http.Header {
	return w.responseWriter.Header()
}

func (w *logResponseWriter) WriteHeader(statusCode int) {
	w.responseWriter.WriteHeader(statusCode)
	w.StatusCode = statusCode
}

func (w *logResponseWriter) Write(b []byte) (int, error) {
	return w.responseWriter.Write(b)
}
