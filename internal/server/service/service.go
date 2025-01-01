package service

import "context"

type AuthService interface {
	Login(ctx context.Context, email, password string) (string, error)
	GetAccessToken(ctx context.Context, refreshToken string) (string, error)
	GetRefreshToken(ctx context.Context, refreshToken string) (string, error)
	RegisterNewUser(ctx context.Context, Email string, pass string) (int64, error)
}