package auth

import (
	"context"
	"fmt"

	descAuth "github.com/igortoigildin/goph-keeper/pkg/auth_v1"
	"github.com/igortoigildin/goph-keeper/pkg/logger"
	"go.uber.org/zap"
)

func (i *Implementation) Login(ctx context.Context, req *descAuth.LoginRequest) (*descAuth.LoginResponse, error) {
	tkn, err := i.authService.Login(ctx, req.Email, req.Password)
	if err != nil {
		logger.Error("login error", zap.Error(err))
		return nil, fmt.Errorf("failed to login: %w", err)
	}

	return &descAuth.LoginResponse{
		RefreshToken: tkn,
	}, nil
}
