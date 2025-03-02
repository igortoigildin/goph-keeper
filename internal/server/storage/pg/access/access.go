package access

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	models "github.com/igortoigildin/goph-keeper/internal/server/models"

	"github.com/igortoigildin/goph-keeper/internal/client/db"
)

const (
	tableName = "access"

	loginColumn  = "login"
	fileIdColumn = "data_id"
)

type AccessRepository struct {
	db db.Client
}

func NewRepository(db db.Client) *AccessRepository {
	return &AccessRepository{
		db: db,
	}
}

func (rep *AccessRepository) GetAccess(ctx context.Context, login string, id string) (*models.FileInfo, error) {
	builder := sq.Select(loginColumn, fileIdColumn).
		PlaceholderFormat(sq.Dollar).
		From(tableName).
		Where(sq.Eq{fileIdColumn: id}).
		Limit(1)

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("error building SQL query: %w", err)
	}

	qr := db.Query{
		Name:     "access_repository.Get",
		QueryRaw: query,
	}

	var file models.FileInfo
	err = rep.db.DB().ScanOneContext(ctx, &file, qr, args...)
	if err != nil {
		return nil, fmt.Errorf("error retrieving info about sepcified user: %w", err)
	}

	return &file, nil
}

func (rep *AccessRepository) SaveAccess(ctx context.Context, login string, id string) error {
	builder := sq.Insert(tableName).
		PlaceholderFormat(sq.Dollar).
		Columns(loginColumn, fileIdColumn).
		Values(login, id)

	query, args, err := builder.ToSql()
	if err != nil {
		return fmt.Errorf("error building SQL query: %w", err)
	}

	qr := db.Query{
		Name:     "access_repository.SaveAccess",
		QueryRaw: query,
	}

	_, err = rep.db.DB().QueryContext(ctx, qr, args...)
	if err != nil {
		return err
	}

	return nil
}
