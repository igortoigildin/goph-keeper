package user

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	models "github.com/igortoigildin/goph-keeper/internal/server/models"

	"github.com/igortoigildin/goph-keeper/internal/client/db"
)

const (
	tableName = "users"

	loginColumn        = "login"
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

func (rep *UserRepository) SaveUser(ctx context.Context, login string, passHash []byte) (int64, error) {
	builder := sq.Insert(tableName).
		PlaceholderFormat(sq.Dollar).
		Columns(loginColumn, passwordHashColumn).
		Values(login, passHash).
		Suffix("ON CONFLICT DO NOTHING RETURNING user_id")

	query, args, err := builder.ToSql()
	if err != nil {
		return 0, fmt.Errorf("error building SQL query: %w", err)
	}

	qr := db.Query{
		Name:     "user_repository.SaveUser",
		QueryRaw: query,
	}

	var id int64
	err = rep.db.DB().QueryRowContext(ctx, qr, args...).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (rep *UserRepository) GetUser(ctx context.Context, login string) (*models.UserInfo, error) {
	builder := sq.Select(loginColumn, passwordHashColumn).
		PlaceholderFormat(sq.Dollar).
		From(tableName).
		Where(sq.Eq{loginColumn: login}).
		Limit(1)

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("error building SQL query: %w", err)
	}

	qr := db.Query{
		Name:     "user_repository.Get",
		QueryRaw: query,
	}

	var user models.UserInfo
	err = rep.db.DB().ScanOneContext(ctx, &user, qr, args...)
	if err != nil {
		return nil, fmt.Errorf("error retrieving info about sepcified user: %w", err)
	}

	return &user, nil
}
