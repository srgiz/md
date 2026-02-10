package repo

import (
	"context"
	"errors"
	"md/internal/lib/kernel/db"

	"golang.org/x/crypto/bcrypt"
)

type UserRepo struct {
	conn db.Conn
}

func NewUserRepo(conn db.Conn) *UserRepo {
	return &UserRepo{conn}
}

func (r *UserRepo) Create(ctx context.Context, id string, password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)

	if err != nil {
		return err
	}

	res, err := r.conn.Exec(ctx, `INSERT INTO users (id, pw_hash) VALUES ($1, $2)`, id, hash)

	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()

	if err != nil {
		return err
	}

	if rowsAffected != 1 {
		return errors.New("no rows affected")
	}

	return nil
}
