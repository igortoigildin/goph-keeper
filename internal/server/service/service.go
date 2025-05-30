package service

import (
	"context"

	model "github.com/igortoigildin/goph-keeper/internal/server/models"
	desc "github.com/igortoigildin/goph-keeper/pkg/upload_v1"
)

type AuthService interface {
	Login(ctx context.Context, email, password string) (string, error)
	RegisterNewUser(ctx context.Context, Email string, pass string) (int64, error)
}

type UploadService interface {
	SaveFile(stream desc.UploadV1_UploadFileServer) error
	SaveBankData(ctx context.Context, data map[string]string, info string) (string, error)
	SaveText(ctx context.Context, text string, info string) (string, error)
	SaveLoginPassword(ctx context.Context, data map[string]string, info string) (string, error)
}

type DownloadService interface {
	DownloadFile(ctx context.Context, id string) ([]byte, string, error)
	DownloadBankData(ctx context.Context, id string) (map[string]string, string, error)
	DownloadText(ctx context.Context, id string) (string, string, error)
	DownloadLoginPassword(ctx context.Context, id string) (map[string]string, string, error)
}

type ListService interface {
	List(ctx context.Context) ([]model.ObjectInfo, error)
}
