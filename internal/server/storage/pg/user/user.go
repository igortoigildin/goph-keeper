package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	models "github.com/igortoigildin/goph-keeper/internal/server/models"

	"github.com/igortoigildin/goph-keeper/internal/client/db"
)

const (
	tableName = "user"

	emailColumn        = "email"
	passwordHashColumn = "password_hash"
)

type UserRepository struct {
	db db.Client
}

func NewRepository(db db.Client) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (rep *UserRepository) SaveUser(ctx context.Context, email string, passHash []byte) (int64, error) {
	builder := sq.Insert(tableName).
		PlaceholderFormat(sq.Dollar).
		Columns(emailColumn, passwordHashColumn).
		Values(email, passHash).
		Suffix("ON CONFLICT DO NOTHING RETURNING user_id")
	query, args, err := builder.ToSql()
	if err != nil {
		return 0, errors.New("error while building SQL query")
	}

	qr := db.Query{
		Name:     "user_repository.SaveUser",
		QueryRaw: query,
	}

	var id int64
	err = rep.db.DB().QueryRowContext(ctx, qr, args...).Scan(&id)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return 0, fmt.Errorf("user already exists: %w", err)
		default:
			return 0, fmt.Errorf("error while creating user: %w", err)
		}
	}

	return id, nil
}

func (rep *UserRepository) GetUser(ctx context.Context, email string) (*models.UserInfo, error) {
	builder := sq.Select(emailColumn, passwordHashColumn).
		PlaceholderFormat(sq.Dollar).
		From(tableName).
		Where(sq.Eq{emailColumn: email}).
		Limit(1)

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, errors.New("error while building SQL query")
	}

	qr := db.Query{
		Name:     "user_repository.Get",
		QueryRaw: query,
	}

	var user models.UserInfo
	err = rep.db.DB().ScanOneContext(ctx, &user, qr, args...)
	if err != nil {
		return nil, errors.New("cannot get new user")
	}

	return &user, nil
}
