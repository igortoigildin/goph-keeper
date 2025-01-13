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
	id, err := i.authService.RegisterNewUser(ctx, req.Email, req.Password)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			logger.Warn("user already exists", zap.Error(err))

			return nil, status.Errorf(codes.AlreadyExists, `User with email %s already exists`, req.Email)
		} else {
			logger.Error("login error", zap.Error(err))

			return nil, status.Error(codes.Unknown, "failed to login")
		}
	}

	return &descAuth.RegisterResponse{
		UserId: id,
	}, nil
}
