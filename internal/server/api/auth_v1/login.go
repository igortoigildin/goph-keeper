package auth

import (
	"context"
	"errors"

	auth "github.com/igortoigildin/goph-keeper/internal/server/service/auth"
	descAuth "github.com/igortoigildin/goph-keeper/pkg/auth_v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (i *Implementation) Login(ctx context.Context, req *descAuth.LoginRequest) (*descAuth.LoginResponse, error) {
	if req.GetLogin() == "" {
		return nil, status.Error(codes.InvalidArgument, "login is required")
	}

	if req.GetPassword() == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}

	tkn, err := i.authService.Login(ctx, req.GetLogin(), req.GetPassword())
	if err != nil {
		if errors.Is(err, auth.ErrInvalidCredentials) {
			return nil, status.Error(codes.Unauthenticated, "invalid credentials")
		} else {
			return nil, status.Error(codes.Unknown, "failed to login")
		}
	}

	return &descAuth.LoginResponse{
		Token: tkn,
	}, nil
}
