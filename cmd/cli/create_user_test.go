package main

import (
	"context"
	"errors"
	"fmt"
	"md/internal/lib/sqldb"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestValidateCreateUser(t *testing.T) {
	var tests = []struct {
		name     string
		pass     string
		countErr int
	}{
		{
			name:     "",
			pass:     "",
			countErr: 2,
		},
		{
			name:     "name",
			pass:     "",
			countErr: 1,
		},
		{
			name:     "",
			pass:     "pass",
			countErr: 2, // потому что консоль не учитывает пустое значение
		},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("Example %d", i), func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			err := newCliApp(map[string]any{
				"conn": sqldb.NewMockConn(ctrl),
			}).RunCli(context.Background(), []string{"", "user:create", test.name, test.pass})

			assert.IsType(t, validator.ValidationErrors{}, err)
			assert.Len(t, err.(validator.ValidationErrors), test.countErr)
		})
	}
}

func TestFailCreateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	conn, tx := sqldb.NewMockConnWithTx(ctrl)

	expectedErr := errors.New("mock error")

	tx.EXPECT().Rollback().Return(nil)
	conn.EXPECT().Exec(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, expectedErr)

	assert.ErrorIs(t, expectedErr, newCliApp(map[string]any{
		"conn": conn,
	}).RunCli(ctx, []string{"", "user:create", "name", "pass"}))
}

func TestSuccessCreateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	conn, tx := sqldb.NewMockConnWithTx(ctrl)

	tx.EXPECT().Commit().Return(nil)
	conn.EXPECT().Exec(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(sqldb.NewMockRowsAffected(1), nil)

	assert.Nil(t, newCliApp(map[string]any{
		"conn": conn,
	}).RunCli(ctx, []string{"", "user:create", "name", "pass"}))
}
