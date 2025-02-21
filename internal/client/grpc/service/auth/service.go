package register

import (
	"context"
	"fmt"

	desc "github.com/igortoigildin/goph-keeper/pkg/auth_v1"
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
				return fmt.Errorf("failed to create user: %w", err) 
			} else if e.Code() == codes.InvalidArgument {
				return fmt.Errorf("invalid argument: %w", err)
			} else {
				return fmt.Errorf("failed to create user: %w", err)
			}
		} else {
			return fmt.Errorf("failed to create user: %w", err)
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
		return "", fmt.Errorf("authentication error: %w", err)
	}

	return resp.RefreshToken, nil
}
