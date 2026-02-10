package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	_ "github.com/lib/pq"
)

const contextKeyTx = "*sql.Tx"

// Conn mockgen -source=internal/lib/kernel/db/conn.go -destination=internal/lib/kernel/db/conn_mock.go -package=db
type Conn interface {
	Master() *sql.DB
	Query(ctx context.Context, query string, args ...any) (rows *sql.Rows, err error)
	QueryRow(ctx context.Context, query string, args ...any) (row *sql.Row)
	Exec(ctx context.Context, query string, args ...any) (res sql.Result, err error)
	Begin(ctx context.Context) (Tx, error)
}

type MasterSlaveConn struct {
	db [2]*sql.DB
}

func NewMasterSlaveConn(driverName string, masterSourceName string, slaveSourceName string) *MasterSlaveConn {
	master, err := sql.Open(driverName, masterSourceName)

	if err != nil {
		panic(err)
	}

	slave, err := sql.Open(driverName, slaveSourceName)

	if err != nil {
		master.Close()
		panic(err)
	}

	return &MasterSlaveConn{[2]*sql.DB{master, slave}}
}

func NewMasterConn(driverName string, dataSourceName string) *MasterSlaveConn {
	master, err := sql.Open(driverName, dataSourceName)

	if err != nil {
		panic(err)
	}

	return &MasterSlaveConn{[2]*sql.DB{master, master}}
}

func (c *MasterSlaveConn) Master() *sql.DB {
	return c.db[0]
}

func (c *MasterSlaveConn) Query(ctx context.Context, query string, args ...any) (rows *sql.Rows, err error) {
	defer func() {
		c.log(ctx, query, err)
	}()

	if sqlTx, ok := ctx.Value(contextKeyTx).(*sql.Tx); ok {
		return sqlTx.QueryContext(ctx, query, args...)
	}

	return c.db[1].QueryContext(ctx, query, args...)
}

func (c *MasterSlaveConn) QueryRow(ctx context.Context, query string, args ...any) (row *sql.Row) {
	defer func() {
		c.log(ctx, query, row.Err())
	}()

	if sqlTx, ok := ctx.Value(contextKeyTx).(*sql.Tx); ok {
		return sqlTx.QueryRowContext(ctx, query, args...)
	}

	return c.db[1].QueryRowContext(ctx, query, args...)
}

func (c *MasterSlaveConn) Exec(ctx context.Context, query string, args ...any) (res sql.Result, err error) {
	defer func() {
		c.log(ctx, query, err)
	}()

	if sqlTx, ok := ctx.Value(contextKeyTx).(*sql.Tx); ok {
		return sqlTx.ExecContext(ctx, query, args...)
	}

	return c.Master().ExecContext(ctx, query, args...)
}

func (c *MasterSlaveConn) log(ctx context.Context, msg string, err error) {
	slog.DebugContext(ctx, fmt.Sprintf("sql: %s", msg))

	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("sql: %s", err.Error()))
	}
}

func (c *MasterSlaveConn) Begin(ctx context.Context) (Tx, error) {
	slog.DebugContext(ctx, "sql: begin")
	sqlTx, err := c.Master().BeginTx(ctx, nil)

	if err != nil {
		return nil, err
	}

	return &TxAdapter{context.WithValue(ctx, contextKeyTx, sqlTx), sqlTx}, nil
}

type Tx interface {
	Context() context.Context
	Commit() error
	Rollback() error
}

type TxAdapter struct {
	ctx   context.Context
	sqlTx *sql.Tx
}

func (t *TxAdapter) Context() context.Context {
	return t.ctx
}

func (t *TxAdapter) Commit() error {
	slog.DebugContext(t.ctx, "sql: commit")
	return t.sqlTx.Commit()
}

func (t *TxAdapter) Rollback() error {
	slog.DebugContext(t.ctx, "sql: rollback")
	err := t.sqlTx.Rollback()

	if errors.Is(err, sql.ErrTxDone) {
		return nil
	}

	return err
}
