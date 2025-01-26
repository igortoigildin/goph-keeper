package storage

import (
	"context"
	"errors"

	models "github.com/igortoigildin/goph-keeper/internal/server/models"
	fl "github.com/igortoigildin/goph-keeper/pkg/file"
)

var (
	ErrUserExists   = errors.New("user already exists")
	ErrUserNotFound = errors.New("user not found")
)

type UserRepository interface {
	GetUser(ctx context.Context, email string) (*models.UserInfo, error)
	SaveUser(ctx context.Context, email string, passHash []byte) (uid int64, err error)
}

type DataRepository interface {
	SaveData(ctx context.Context, file *fl.File, email string) error
}
