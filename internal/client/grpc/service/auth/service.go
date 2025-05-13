package register

import (
	"context"
	"fmt"
	"os"

	desc "github.com/igortoigildin/goph-keeper/pkg/auth_v1"
	"github.com/igortoigildin/goph-keeper/pkg/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

type AuthService struct {
	addr   string
	client desc.AuthV1Client
}

func New(addr string) *AuthService {
	logger.Initialize("info")
	return &AuthService{
		addr: addr,
	}
}

func (auth *AuthService) RegisterNewUser(ctx context.Context, login, pass string) error {
	// Load TLS credentials
	creds, err := credentials.NewClientTLSFromFile("certs/server.crt", "")
	if err != nil {
		logger.Error("failed to load TLS certificates: %w", zap.Error(err))
		return fmt.Errorf("failed to load TLS certificates: %w", err)
	}

	// Create gRPC connection with TLS
	conn, err := grpc.Dial(auth.addr, grpc.WithTransportCredentials(creds))
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
	var opts []grpc.DialOption

	// Use insecure credentials in test mode
	if os.Getenv("TEST_ENV") == "true" {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	} else {
		// Load TLS credentials
		creds, err := credentials.NewClientTLSFromFile("certs/server.crt", "")
		if err != nil {
			logger.Error("failed to load TLS certificates: %w", zap.Error(err))
			return "", fmt.Errorf("failed to load TLS certificates: %w", err)
		}
		opts = append(opts, grpc.WithTransportCredentials(creds))
	}

	// Create gRPC connection
	conn, err := grpc.Dial(auth.addr, opts...)
	if err != nil {
		return "", fmt.Errorf("error dialing client: %w", err)
	}
	defer conn.Close()

	auth.client = desc.NewAuthV1Client(conn)

	resp, err := auth.client.Login(ctx, &desc.LoginRequest{Login: login, Password: pass})
	if err != nil {
		return "", fmt.Errorf("authentication error: %w", err)
	}

	return resp.GetToken(), nil
}
