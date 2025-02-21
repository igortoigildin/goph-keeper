package auth

import (
	"context"
	"errors"

	descAuth "github.com/igortoigildin/goph-keeper/pkg/auth_v1"
	"github.com/igortoigildin/goph-keeper/pkg/logger"
	pgx "github.com/jackc/pgx/v4"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (i *Implementation) Register(ctx context.Context, req *descAuth.RegisterRequest) (*descAuth.RegisterResponse, error) {
	if req.GetLogin() == "" {
		return nil, status.Error(codes.InvalidArgument, "login is requeired")
	}

	if req.GetPassword() == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}
	
	id, err := i.authService.RegisterNewUser(ctx, req.Login, req.Password)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			logger.Warn("User with such login already exists", zap.Error(err))

			return nil, status.Errorf(codes.AlreadyExists, `User with email %s already exists`, req.Login)
		} else {
			logger.Error("login error", zap.Error(err))

			return nil, status.Error(codes.Unknown, "failed to login")
		}
	}

	return &descAuth.RegisterResponse{
		UserId: id,
	}, nil
}
