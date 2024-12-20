package auth

import (
	"context"

	desc "github.com/igortoigildin/goph-keeper/pkg/auth_v1"
)

type AuthService interface {
	Login(ctx context.Context, username, password string) (string, error)
	Register(ctx context.Context, email, password string) (int64, error)
	GetAccessToken(ctx context.Context, token string) (string, error)
	GetRefreshToken(ctx context.Context, token string) (string, error)
}

type Implementation struct {
	desc.UnimplementedAuthV1Server
	authService AuthService
}

func NewImplementation(authService AuthService) *Implementation {
	return &Implementation{
		authService: authService,
	}
}
