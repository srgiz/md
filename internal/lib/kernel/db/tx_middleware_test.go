package db

import (
	"context"
	"errors"
	"fmt"
	"md/internal/lib/kernel/cmdbus"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestPanicHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	conn, tx := NewMockConnWithTx(ctrl)
	wasCommit := false
	wasRollback := false

	tx.EXPECT().Commit().DoAndReturn(func() error {
		wasCommit = true
		return nil
	}).AnyTimes()

	tx.EXPECT().Rollback().DoAndReturn(func() error {
		wasRollback = true
		return nil
	})

	assert.PanicsWithValue(t, "mock panic", func() {
		_, _ = NewTxMiddleware(conn).Handle(ctx, 0, func(ctx context.Context, cmd any) (any, error) {
			panic("mock panic")
		})
	})

	assert.False(t, wasCommit)
	assert.True(t, wasRollback)
}

func TestCommitAndRollback(t *testing.T) {
	var tests = []struct {
		name        string
		handler     cmdbus.Handler
		errRes      error
		wasCommit   bool
		errCommit   error
		wasRollback bool
		errRollback error
	}{
		{
			name: "commit",
			handler: func(ctx context.Context, cmd any) (any, error) {
				return 0, nil
			},
			wasCommit:   true,
			errCommit:   nil,
			wasRollback: false,
			errRollback: nil,
			errRes:      nil,
		},
		{
			name: "rollback",
			handler: func(ctx context.Context, cmd any) (any, error) {
				return nil, errors.New("mock error")
			},
			wasCommit:   false,
			errCommit:   nil,
			wasRollback: true,
			errRollback: nil,
			errRes:      errors.New("mock error"),
		},
		{
			name: "err commit",
			handler: func(ctx context.Context, cmd any) (any, error) {
				return 0, nil
			},
			wasCommit:   true,
			errCommit:   errors.New("mock commit"),
			wasRollback: true,
			errRollback: nil,
			errRes:      errors.New("mock commit"),
		},
		{
			name: "err rollback",
			handler: func(ctx context.Context, cmd any) (any, error) {
				return nil, errors.New("mock error")
			},
			wasCommit:   false,
			errCommit:   nil,
			wasRollback: true,
			errRollback: errors.New("mock rollback"),
			errRes:      errors.New("mock rollback"),
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("Example %s", test.name), func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.Background()
			conn, tx := NewMockConnWithTx(ctrl)
			wasCommit := false
			wasRollback := false

			tx.EXPECT().Commit().DoAndReturn(func() error {
				wasCommit = true
				return test.errCommit
			}).AnyTimes()

			tx.EXPECT().Rollback().DoAndReturn(func() error {
				wasRollback = true
				return test.errRollback
			}).AnyTimes()

			_, err := NewTxMiddleware(conn).Handle(ctx, 0, test.handler)

			assert.Equal(t, test.wasCommit, wasCommit)
			assert.Equal(t, test.wasRollback, wasRollback)
			assert.Equal(t, test.errRes, err)
		})
	}
}
