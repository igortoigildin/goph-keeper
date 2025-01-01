package pg

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)


type UserSaver interface {
	SaveUser(
		ctx context.Context,
		email string,
		passHash []byte,
	) (uid int64, err error)
}

type UserRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (rep *UserRepository) SaveUser(ctx context.Context, email string, passHash []byte) (int64, error) {
	var userID int64
	query := `INSERT INTO users (email, password_hash)
	VALUES ($1, $2) ON CONFLICT DO NOTHING RETURNING user_id`

	args := []any{
		email,
		passHash,
	}

	err := rep.db.QueryRowContext(ctx, query, args...).Scan(&userID)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return 0, fmt.Errorf("user already exists: %w", err)
		default:
			return 0, fmt.Errorf("error while creating user: %w", err)
		}
	}

	return userID, nil
}

// func (rep *UserRepository) User(ctx context.Context, email string) (models.UserInfo, error)