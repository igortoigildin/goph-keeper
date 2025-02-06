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
	SaveFile(ctx context.Context, file *fl.File, login string, id string) error
	SaveTextData(ctx context.Context, data any, login string, id string) error
}

type AccessRepository interface {
	GetAccess(ctx context.Context, login string, id string) (*models.FileInfo, error)
	SaveAccess(ctx context.Context, login string, id string) (error)
}

