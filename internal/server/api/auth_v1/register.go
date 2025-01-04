package auth

import (
	"context"
	"errors"

	descAuth "github.com/igortoigildin/goph-keeper/pkg/auth_v1"
	"github.com/igortoigildin/goph-keeper/pkg/logger"
	"go.uber.org/zap"
)

func (i *Implementation) Register(ctx context.Context, req *descAuth.RegisterRequest) (*descAuth.RegisterResponse, error) {
	id, err := i.authService.RegisterNewUser(ctx, req.Email, req.Password)
	if err != nil {
		logger.Error("registration error", zap.Error(err))
		return nil, errors.New("failed to login")
	}

	return &descAuth.RegisterResponse{
		UserId: id,
	}, nil
}
