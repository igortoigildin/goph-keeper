package auth

import (
	"context"
	"errors"
	"fmt"

	"github.com/igortoigildin/goph-keeper/internal/server/domain/models"
	"github.com/igortoigildin/goph-keeper/internal/server/storage"
	"github.com/igortoigildin/goph-keeper/pkg/logger"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
)

// type AuthService interface {
// 	Login(ctx context.Context, username, password string) (string, error)
// 	Register(ctx context.Context, email, password string) (int64, error)
// 	GetAccessToken(ctx context.Context, token string) (string, error)
// 	GetRefreshToken(ctx context.Context, token string) (string, error)
// }

type UserProvider interface {
	User(ctx context.Context, email string) (models.User, error)
}

type Auth struct {
	userProvider UserProvider
}

func New() *Auth {
	return &Auth{}
}

// Login checks if user with given credentials exists in the system and returns access token.
// If user exists, but password is incorrect, returns error.
// If user doesn't exist, returns error.
func (a *Auth) Login(ctx context.Context, email, password string) (string, error) {
	const op = "Auth.Login"
	logger.Info("attempting to login user")

	user, err := a.userProvider.User(ctx, email)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			logger.Warn("user not found", zap.Error(err))

			return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}

		logger.Error("failed to get user", zap.Error(err))

		return "", fmt.Errorf("%s: %w", op, err)
	}

	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		logger.Info("invalid credentials", zap.Error(err))

		return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	logger.Info("user logged in successfully")

	// TODO: add jwt token
}
