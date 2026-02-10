package db

import (
	"context"
	"md/internal/lib/kernel/cmdbus"
)

type TxMiddleware struct {
	conn Conn
}

func NewTxMiddleware(conn Conn) cmdbus.Middleware {
	return &TxMiddleware{conn}
}

func (m *TxMiddleware) Handle(ctx context.Context, cmd any, next cmdbus.Handler) (res any, err error) {
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
