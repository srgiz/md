package repo

import (
	"context"
	"md/internal/lib/sqldb"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestCreateToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	conn := sqldb.NewMockConn(ctrl)
	row := sqldb.NewMockRow(ctrl)
	password := "pw"
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)

	conn.EXPECT().QueryRow(gomock.Any(), gomock.Any(), gomock.Any()).Return(row)
	row.EXPECT().Scan(gomock.Any()).DoAndReturn(func(h *string) error {
		*h = string(hash)
		return nil
	})

	token, err := NewUserRepo(conn).CreateToken(context.Background(), "0", password)

	assert.Nil(t, err)
	assert.Len(t, strings.Split(token, "."), 3)
}
