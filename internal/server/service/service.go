package service

import (
	"context"

	desc "github.com/igortoigildin/goph-keeper/pkg/upload_v1"
)

type AuthService interface {
	Login(ctx context.Context, email, password string) (string, error)
	GetAccessToken(ctx context.Context, refreshToken string) (string, error)
	GetRefreshToken(ctx context.Context, refreshToken string) (string, error)
	RegisterNewUser(ctx context.Context, Email string, pass string) (int64, error)
}


type UploadService interface {
	SaveFile(stream desc.UploadV1_UploadFileServer) error
	SaveBankData(ctx context.Context, data map[string]string) error
	SaveText(ctx context.Context, text string) error
	SaveLoginPassword(ctx context.Context, data map[string]string) error
}