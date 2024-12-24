package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	models "github.com/igortoigildin/goph-keeper/internal/server/models"
	"github.com/igortoigildin/goph-keeper/internal/server/storage"
	"github.com/igortoigildin/goph-keeper/pkg/logger"
	utils "github.com/igortoigildin/goph-keeper/pkg/utils"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

const (
	refreshTokenSecretKey = "W4/X+LLjehdxptt4YgGFCvMpq5ewptpZZYRHY6A72g0="
	accessTokenSecretKey  = "VqvguGiffXILza1f44TWXowDT4zwf03dtXmqWW4SYyE="

	refreshTokenExpiration = 60 * time.Minute
	accessTokenExpiration  = 2 * time.Minute
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserExists         = errors.New("user already exists")
)

// type AuthService interface {
// 	Login(ctx context.Context, username, password string) (string, error)
// 	Register(ctx context.Context, email, password string) (int64, error)
// 	GetAccessToken(ctx context.Context, token string) (string, error)
// 	GetRefreshToken(ctx context.Context, token string) (string, error)
// }

type Auth struct {
	userProvider UserProvider
	userSaver    UserSaver
}

type UserProvider interface {
	User(ctx context.Context, email string) (models.UserInfo, error)
}

type UserSaver interface {
	SaveUser(
		ctx context.Context,
		email string,
		passHash []byte,
	) (uid int64, err error)
}

func New() *Auth {
	return &Auth{}
}

// Login checks if user with given credentials exists in the system and returns access token.
// If user exists, but password is incorrect, returns error.
// If user doesn't exist, returns error.
func (a *Auth) Login(ctx context.Context, email, password string) (string, error) {
	const op = "Auth.Login"
	logger.Info("attempting to log in user")

	// identify user by email
	user, err := a.userProvider.User(ctx, email)
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

	// generate refresh token
	refreshToken, err := utils.GenerateToken(user, []byte(refreshTokenSecretKey), refreshTokenExpiration)
	if err != nil {
		return "", errors.New("failed to generate token")
	}

	logger.Info("user logged in successfully")

	return refreshToken, nil
}

func (a *Auth) GetAccessToken(ctx context.Context, refreshToken string) (string, error) {
	claims, err := utils.VeryfyToken(refreshToken, []byte(refreshTokenSecretKey))
	if err != nil {
		return "", errors.New("invalid refresh token")
	}

	accessToken, err := utils.GenerateToken(models.UserInfo{
		Email: claims.Email,
	}, []byte(accessTokenSecretKey), accessTokenExpiration,
	)
	if err != nil {
		return "", errors.New("internal error")
	}

	return accessToken, nil
}

func (a *Auth) GetRefreshToken(ctx context.Context, refreshToken string) (string, error) {
	claims, err := utils.VeryfyToken(refreshToken, []byte(refreshTokenSecretKey))
	if err != nil {
		return "", errors.New("invalid token")
	}

	token, err := utils.GenerateToken(models.UserInfo{
		Email: claims.Email,
	},
		[]byte(refreshTokenSecretKey),
		refreshTokenExpiration,
	)
	if err != nil {
		return "", err
	}

	return token, nil
}

// RegisterNewUser registers new user in the system and returns user ID.
// If user with given username already exists, returns error.
func (a *Auth) RegisterNewUser(ctx context.Context, Email string, pass string) (int64, error) {
	const op = "auth.RegisterNewUser"
	logger.Info("registering user")

	passHash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		logger.Error("failed to generate password hash", zap.Error(err))

		return 0, fmt.Errorf("%s: %w", op, zap.Error(err))
	}

	id, err := a.userSaver.SaveUser(ctx, Email, passHash)
	if err != nil {
		if errors.Is(err, storage.ErrUserExists) {
			logger.Warn("user already exists", zap.Error(err))

			return 0, fmt.Errorf("%s: %w", op, ErrUserExists)
		}

		logger.Error("failed to save user", zap.Error(err))
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	logger.Info("user registered")
	return id, nil
}
