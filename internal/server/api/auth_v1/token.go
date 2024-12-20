package auth

import (
	"context"
	"fmt"

	descAuth "github.com/igortoigildin/goph-keeper/pkg/auth_v1"
	"github.com/igortoigildin/goph-keeper/pkg/logger"
	"go.uber.org/zap"
)

func (i *Implementation) GetAccessToken(ctx context.Context, req *descAuth.GetAccessTokenRequest) (*descAuth.GetAccessTokenResponse, error) {
	token, err := i.authService.GetAccessToken(ctx, req.RefreshToken)
	if err != nil {
		logger.Error("error while requesting access token", zap.Error(err))
		return nil, fmt.Errorf("access token error: %w", err)
	}
	return &descAuth.GetAccessTokenResponse{
		AccessToken: token,
	}, nil
}

func (i *Implementation) GetRefreshToken(ctx context.Context, req *descAuth.GetRefreshTokenRequest) (*descAuth.GetRefreshTokenResponse, error) {
	token, err := i.authService.GetRefreshToken(ctx, req.RefreshToken)
	if err != nil {
		logger.Error("error while requesting refresh token", zap.Error(err))
		return nil, fmt.Errorf("refresh token error: %w", err)
	}
	return &descAuth.GetRefreshTokenResponse{
		RefreshToken: token,
	}, nil
}
