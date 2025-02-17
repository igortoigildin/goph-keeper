package register

import (
	"context"
	"errors"
	"fmt"

	desc "github.com/igortoigildin/goph-keeper/pkg/auth_v1"
	"github.com/igortoigildin/goph-keeper/pkg/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

type AuthService struct {
	addr   string
	client desc.AuthV1Client
}

func New(addr string) *AuthService {
	return &AuthService{
		addr: addr,
	}
}

func (auth *AuthService) RegisterNewUser(ctx context.Context, login, pass string) error {
	conn, err := grpc.NewClient(auth.addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("error dialing client: %w", err)
	}
	defer conn.Close()

	auth.client = desc.NewAuthV1Client(conn)

	_, err = auth.client.Register(ctx, &desc.RegisterRequest{Login: login, Password: pass})
	if err != nil {
		if e, ok := status.FromError(err); ok {
			if e.Code() == codes.AlreadyExists {
				logger.Error(`User already exists`, zap.String("error", e.Message()))

				return fmt.Errorf("failed to create user: %s", err)
			} else {
				logger.Error(e.Message(), zap.Any("status", e.Code()))

				return fmt.Errorf("failed to create user: %s", err)
			}
		} else {
			logger.Error("failed to parse error", zap.Error(err))

			return fmt.Errorf("failed to create user: %s", err)
		}
	}

	return nil
}

func (auth *AuthService) Login(ctx context.Context, login, pass string) (string, error) {
	conn, err := grpc.NewClient(auth.addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return "", fmt.Errorf("error dialing client: %w", err)
	}
	defer conn.Close()

	auth.client = desc.NewAuthV1Client(conn)

	resp, err := auth.client.Login(ctx, &desc.LoginRequest{Login: login, Password: pass})

	if err != nil {
		logger.Error("login error", zap.Error(err))

		return "", errors.New("failed to login")
	}

	return resp.RefreshToken, nil
}
