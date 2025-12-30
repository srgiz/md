package cmdbus

import (
	"context"
)

type Transaction interface {
	Context() context.Context
	Commit() error
	Rollback() error
}

type TransactionManager interface {
	Begin(ctx context.Context) (Transaction, error)
}

func NewTransactionMw(manager TransactionManager) *transactionMiddleware {
	return &transactionMiddleware{manager}
}

type transactionMiddleware struct {
	manager TransactionManager
}

func (m *transactionMiddleware) Handle(ctx context.Context, cmd any, next handler) (res any, err error) {
	tx, errBegin := m.manager.Begin(ctx)

	if errBegin != nil {
		return nil, errBegin
	}

	commit := false

	defer func() {
		if !commit {
			tx.Rollback()
		}
	}()

	defer func() {
		if rec := recover(); rec != nil {
			panic(rec)
		}

		if err == nil { // только если нет ошибки после next()
			if errCommit := tx.Commit(); errCommit != nil {
				panic(errCommit) // todo: ???
			}

			commit = true
		}
	}()

	return next(tx.Context(), cmd)
}
