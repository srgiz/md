package kernel

import (
	"context"
	"md/internal/lib/sqldb"
)

func NewTxBusProvider(serviceName string, conn sqldb.Conn) func(app *App) {
	return func(app *App) {
		app.AddService(serviceName, NewTxBusMiddleware(conn))
	}
}

type txBusMiddleware struct {
	conn sqldb.Conn
}

func NewTxBusMiddleware(conn sqldb.Conn) BusMiddleware {
	return &txBusMiddleware{conn}
}

func (m *txBusMiddleware) Handle(ctx context.Context, cmd any, next cmdHandler) (res any, err error) {
	tx, errBegin := m.conn.Begin(ctx)

	if errBegin != nil {
		return nil, errBegin
	}

	commit := false

	defer func() {
		if !commit {
			if errRollback := tx.Rollback(); errRollback != nil {
				err = errRollback
			}
		}
	}()

	defer func() {
		if rec := recover(); rec != nil {
			panic(rec) // Rollback
		}

		if err == nil {
			if errCommit := tx.Commit(); errCommit != nil {
				err = errCommit
			} else {
				commit = true
			}
		}
	}()

	return next(tx.Context(), cmd)
}
