package postgres

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"md/internal/domain/cmdbus"
)

type Tx struct {
	ctx   context.Context
	sqlTx *sql.Tx
}

func (t *Tx) Context() context.Context {
	return t.ctx
}

func (t *Tx) Commit() error {
	slog.DebugContext(t.ctx, "sql: commit")
	return t.sqlTx.Commit()
}

func (t *Tx) Rollback() error {
	slog.DebugContext(t.ctx, "sql: rollback")
	err := t.sqlTx.Rollback()

	if errors.Is(err, sql.ErrTxDone) {
		return nil
	}

	return err
}

func (c *Conn) Begin(ctx context.Context) (cmdbus.Transaction, error) {
	slog.DebugContext(ctx, "sql: begin")
	sqlTx, err := c.Master().BeginTx(ctx, nil)

	if err != nil {
		return nil, err
	}

	return &Tx{context.WithValue(ctx, contextTxKey, sqlTx), sqlTx}, nil
}
