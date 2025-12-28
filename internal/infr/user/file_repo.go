package user

import (
	"context"
	"md/internal/domain/user"
	"md/internal/infr/postgres"

	"golang.org/x/crypto/bcrypt"
)

type UserRepository struct {
	conn *postgres.Conn
}

func NewUserRepository(conn *postgres.Conn) user.UserRepository {
	return &UserRepository{conn}
}

func (r *UserRepository) Create(ctx context.Context, id string, password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)

	if err != nil {
		return err
	}

	_, err = r.conn.ExecContext(ctx, `INSERT INTO users (id, pw_hash) VALUES ($1, $2)`, id, hash)
	return err
}
