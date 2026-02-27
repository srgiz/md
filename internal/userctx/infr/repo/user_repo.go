package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"md/internal/lib/sqldb"
	"md/internal/userctx/domain/login"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserRepo struct {
	conn sqldb.Conn
}

func NewUserRepo(conn sqldb.Conn) *UserRepo {
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
		return errors.New("user: no rows affected")
	}

	return nil
}

func (r *UserRepo) CreateToken(ctx context.Context, id string, password string) (string, error) {
	row := r.conn.QueryRow(ctx, `SELECT pw_hash FROM users WHERE id = $1`, id)
	var pw_hash string

	err := row.Scan(&pw_hash)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", login.ErrIncorrectCredentials
		}

		panic(err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(pw_hash), []byte(password))

	if err != nil {
		slog.Debug(fmt.Sprintf("user: %s", err.Error()), "userId", id)
		return "", login.ErrIncorrectCredentials
	}

	return jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp":    time.Now().Add(time.Hour).Unix(),
		"userId": id,
	}).SignedString([]byte(os.Getenv("APP_JWT_KEY")))
}
