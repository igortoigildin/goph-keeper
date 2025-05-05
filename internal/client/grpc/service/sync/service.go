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
	login         = "login"
	password      = "password"
	loginPassword = "login_password_"
	bankData      = "bank_data_"
	textData      = "text_data_"
	binData       = "bin_data_"
)

type Syncer interface {
	ListAllData(addr string) ([]*desc.ObjectInfo, error)
}

type ClientService struct {
	client desc.SyncV1Client
}

func New() *ClientService {
	return &ClientService{}
}

func (s *ClientService) ListAllData(addr string) ([]*desc.ObjectInfo, error) {
	// Load TLS credentials
	creds, err := credentials.NewClientTLSFromFile("certs/server.crt", "")
	if err != nil {
		logger.Error("failed to load TLS certificates: %w", zap.Error(err))

		return nil, fmt.Errorf("failed to load TLS certificates: %w", err)
	}

	// Create gRPC connection with TLS
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(creds))
	if err != nil {
		return nil, fmt.Errorf("error dialing client: %w", err)
	}
	defer conn.Close()

	s.client = desc.NewSyncV1Client(conn)

	ss, err := session.LoadSession()
	if err != nil {
		return nil, fmt.Errorf("error loading session: %w", err)
	}

	md := metadata.Pairs(login, ss.Login, "authorization", "Bearer "+ss.Token)

	ctx := metadata.NewOutgoingContext(context.Background(), md)

	resp, err := s.client.GetObjectList(ctx, &desc.SyncRequest{Login: ss.Login})
	if err != nil {
		return nil, fmt.Errorf("error getting object list: %w", err)
	}

	objects := resp.Objects

	logger.Info("Object list: %v", zap.Any("objects", objects))

	return resp.Objects, nil
}
