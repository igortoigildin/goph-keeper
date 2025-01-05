package register

import (
	"context"

	desc "github.com/igortoigildin/goph-keeper/pkg/auth_v1"
	"github.com/igortoigildin/goph-keeper/pkg/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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

func (auth *AuthService) RegisterNewUser(ctx context.Context, email, pass string) error {
	conn, err := grpc.Dial(auth.addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	defer conn.Close()

	auth.client = desc.NewAuthV1Client(conn)

	if _, err := auth.client.Register(ctx, &desc.RegisterRequest{Email: email, Password: pass}); err != nil {
		logger.Fatal("error while registering user", zap.Error(err))
	}

	return nil
}

func (auth *AuthService) Login(ctx context.Context, email, pass string) error {
	conn, err := grpc.Dial(auth.addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	defer conn.Close()

	auth.client = desc.NewAuthV1Client(conn)

	if _, err := auth.client.Login(ctx, &desc.LoginRequest{Email: email, Password: pass}); err != nil {
		logger.Fatal("authentication error:", zap.Error(err))
	}

	return nil
}
