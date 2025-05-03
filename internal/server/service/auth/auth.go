package auth

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	models "github.com/igortoigildin/goph-keeper/internal/server/models"
	"github.com/igortoigildin/goph-keeper/internal/server/service"
	"github.com/igortoigildin/goph-keeper/internal/server/storage"
	"github.com/igortoigildin/goph-keeper/pkg/logger"
	utils "github.com/igortoigildin/goph-keeper/pkg/utils"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

const (
	tokenExpiration = 60 * time.Minute
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserExists         = errors.New("user already exists")
)

type UserRepository interface {
	GetUser(ctx context.Context, login string) (*models.UserInfo, error)
	SaveUser(ctx context.Context, login string, passHash []byte) (uid int64, err error)
}

type authServ struct {
	userRepo UserRepository
}

func New(userRepo UserRepository) service.AuthService {
	return &authServ{
		userRepo: userRepo,
	}
}

// Login checks if user with given credentials exists in the system and returns access token.
// If user exists, but password is incorrect, returns error.
// If user doesn't exist, returns error.
func (a *authServ) Login(ctx context.Context, login, password string) (string, error) {
	const op = "Auth.Login"
	logger.Info("attempting to login user")

	// identify user by email
	user, err := a.userRepo.GetUser(ctx, login)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			logger.Warn("user not found", zap.Error(err))

			return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}

		logger.Error("failed to get user", zap.Error(err))

		return "", fmt.Errorf("%s: %w", op, err)
	}

	// compare password hash
	if !utils.VerifyPassword(string(user.Hash), password) {
		return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	jwtSecret := os.Getenv("JWT_SECRET")

	// generate refresh token
	refreshToken, err := utils.GenerateToken(*user, []byte(jwtSecret), tokenExpiration)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	logger.Info("user logged in successfully:", zap.String("login", login))

	return refreshToken, nil
}

// RegisterNewUser registers new user in the system and returns user ID.
// If user with given username already exists, returns error.
func (a *authServ) RegisterNewUser(ctx context.Context, login string, pass string) (int64, error) {
	op := "server/service/auth"

	passHash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		logger.Error("failed to generate password hash", zap.Error(err))

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	id, err := a.userRepo.SaveUser(ctx, login, passHash)
	if err != nil {
		logger.Error("failed to save user", zap.Error(err))

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	logger.Info("user registered:", zap.String("login", login))

	return id, nil
}
