package syncData

import (
	"context"
	"fmt"

	"github.com/igortoigildin/goph-keeper/pkg/logger"
	"github.com/igortoigildin/goph-keeper/pkg/session"
	desc "github.com/igortoigildin/goph-keeper/pkg/sync_v1"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
)

const (
	login    = "login"
	password = "password"
)

type Syncer interface {
	SyncAllData(addr string) error
}

type ClientService struct {
	client desc.SyncV1Client
}

func New() *ClientService {
	return &ClientService{}
}

func (s *ClientService) SyncAllData(addr string) error {
	// Load TLS credentials
	creds, err := credentials.NewClientTLSFromFile("certs/server.crt", "")
	if err != nil {
		logger.Error("failed to load TLS certificates: %w", zap.Error(err))

		return fmt.Errorf("failed to load TLS certificates: %w", err)
	}

	// Create gRPC connection with TLS
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(creds))
	if err != nil {
		return fmt.Errorf("error dialing client: %w", err)
	}
	defer conn.Close()

	s.client = desc.NewSyncV1Client(conn)

	ss, err := session.LoadSession()
	if err != nil {
		return fmt.Errorf("error loading session: %w", err)
	}

	md := metadata.Pairs(login, ss.Login, "authorization", "Bearer "+ss.Token)

	ctx := metadata.NewOutgoingContext(context.Background(), md)

	resp, err := s.client.GetObjectList(ctx, &desc.SyncRequest{Login: ss.Login})
	if err != nil {
		return fmt.Errorf("error getting object list: %w", err)
	}

	logger.Info("Object list: %v", zap.Any("objects", resp.Objects))

	return nil
}
