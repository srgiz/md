package sqldb

import (
	"context"
	"database/sql"
	"errors"

	"github.com/golang/mock/gomock"
)

func NewMockConnWithTx(ctrl *gomock.Controller) (*MockConn, *MockTx) {
	conn := NewMockConn(ctrl)
	tx := NewMockTx(ctrl)

	conn.EXPECT().Begin(gomock.Any()).DoAndReturn(func(ctx context.Context) (Tx, error) {
		tx.EXPECT().Context().Return(ctx)
		return tx, nil
	})

	return conn, tx
}

type mockSqlResult struct {
	lastInsertId *int64
	rowsAffected *int64
}

func NewMockRowsAffected(rowsAffected int64) sql.Result {
	return &mockSqlResult{
		lastInsertId: nil,
		rowsAffected: &rowsAffected,
	}
}

func (res *mockSqlResult) LastInsertId() (int64, error) {
	if res.lastInsertId != nil {
		return *res.lastInsertId, nil
	}

	return 0, errors.New("LastInsertId mock error")
}

func (res *mockSqlResult) RowsAffected() (int64, error) {
	if res.rowsAffected != nil {
		return *res.rowsAffected, nil
	}

	return 0, errors.New("RowsAffected mock error")
}

/*
func NewMockSqlResult[TId int64 | *int64, TRows int64 | *int64](lastInsertId TId, rowsAffected TRows) sql.Result {
	var id, rows *int64

	switch any(lastInsertId).(type) {
	case *int64:
		id = any(lastInsertId).(*int64)
	case int64:
		pointerId := any(lastInsertId).(int64)
		id = &pointerId
	}

	switch any(rowsAffected).(type) {
	case *int64:
		rows = any(rowsAffected).(*int64)
	case int64:
		pointerRows := any(rowsAffected).(int64)
		rows = &pointerRows
	}

	return &mockSqlResult{
		lastInsertId: id,
		rowsAffected: rows,
	}
}*/
