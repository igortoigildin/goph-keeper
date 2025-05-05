package upload

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/igortoigildin/goph-keeper/pkg/logger"
	"github.com/igortoigildin/goph-keeper/pkg/session"
	desc "github.com/igortoigildin/goph-keeper/pkg/upload_v1"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
)

const (
	login    = "login"
	password = "password"
)

type Sender interface {
	SendPassword(addr, loginStr, passStr string, id string, meta string) error
	SendText(addr, text string, id string, meta string) error
	SendFile(addr string, filePath string, batchSize int, id, meta string) error
	SendBankDetails(addr, cardNumber, cvc, expDate string, id, meta string) error
}

type ClientService struct {
	client desc.UploadV1Client
}

func New() *ClientService {
	return &ClientService{}
}

func (s *ClientService) SendPassword(addr, loginStr, passStr string, id string, meta string) (string, error) {
	// Load TLS credentials
	creds, err := credentials.NewClientTLSFromFile("certs/server.crt", "")
	if err != nil {
		logger.Error("failed to load TLS certificates: %w", zap.Error(err))

		return "", fmt.Errorf("failed to load TLS certificates: %w", err)
	}

	// Create gRPC connection with TLS
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(creds))
	if err != nil {
		return "", fmt.Errorf("error dialing client: %w", err)
	}
	defer conn.Close()

	s.client = desc.NewUploadV1Client(conn)

	ss, err := session.LoadSession()
	if err != nil {
		return "", fmt.Errorf("error loading session: %w", err)
	}

	md := metadata.Pairs(login, ss.Login, "id", id, "authorization", "Bearer "+ss.Token)

	ctx := metadata.NewOutgoingContext(context.Background(), md)

	etag, err := s.uploadPassword(ctx, loginStr, passStr, meta)
	if err != nil {
		return "", err
	}

	return etag, nil
}

func (s *ClientService) uploadPassword(ctx context.Context, loginStr, passStr, meta string) (string, error) {
	data := make(map[string]string, 2)
	data[login] = loginStr
	data[password] = passStr
	data["metadata"] = meta

	resp, err := s.client.UploadPassword(ctx, &desc.UploadPasswordRequest{Data: data, Metadata: meta})
	if err != nil {
		return "", fmt.Errorf("error uploading credentials: %w", err)
	}

	return resp.Etag, nil
}

func (s *ClientService) SendBankDetails(addr, cardNumber, cvc, expDate string, id, meta string) (string, error) {
	creds, err := credentials.NewClientTLSFromFile("certs/server.crt", "")
	if err != nil {
		logger.Error("failed to load TLS certificates: %w", zap.Error(err))

		return "", fmt.Errorf("failed to load TLS certificates: %w", err)
	}

	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(creds))
	if err != nil {
		return "", fmt.Errorf("error dialing client: %w", err)
	}
	defer conn.Close()

	s.client = desc.NewUploadV1Client(conn)

	ss, err := session.LoadSession()
	if err != nil {
		return "", fmt.Errorf("error loading session: %w", err)
	}

	md := metadata.Pairs(login, ss.Login, "id", id, "authorization", "Bearer "+ss.Token)

	ctx := metadata.NewOutgoingContext(context.Background(), md)

	etag, err := s.uploadBankDetails(ctx, cardNumber, cvc, expDate, meta)
	if err != nil {
		return "", err
	}

	return etag, nil
}

func (s *ClientService) uploadBankDetails(ctx context.Context, cardNumber, cvc, expDate, meta string) (string, error) {
	data := make(map[string]string, 3)
	data["card_number"] = cardNumber
	data["CVC"] = cvc
	data["expiration_date"] = expDate
	data["metadata"] = meta

	resp, err := s.client.UploadBankData(ctx, &desc.UploadBankDataRequest{Data: data, Metadata: meta})
	if err != nil {
		return "", fmt.Errorf("error uploading bank details: %w", err)
	}

	return resp.Etag, nil
}

func (s *ClientService) SendText(addr, text string, id string, info string) (string, error) {
	creds, err := credentials.NewClientTLSFromFile("certs/server.crt", "")
	if err != nil {
		logger.Error("failed to load TLS certificates: %w", zap.Error(err))

		return "", fmt.Errorf("failed to load TLS certificates: %w", err)
	}

	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(creds))
	if err != nil {
		return "", fmt.Errorf("error dialing client: %w", err)
	}
	defer conn.Close()

	s.client = desc.NewUploadV1Client(conn)

	ss, err := session.LoadSession()
	if err != nil {
		return "", fmt.Errorf("error loading session: %w", err)
	}

	md := metadata.Pairs(login, ss.Login, "id", id, "authorization", "Bearer "+ss.Token)

	ctx := metadata.NewOutgoingContext(context.Background(), md)

	etag, err := s.uploadText(ctx, text, info)
	if err != nil {
		return "", err
	}

	return etag, nil
}

func (s *ClientService) uploadText(ctx context.Context, text, info string) (string, error) {
	resp, err := s.client.UploadText(ctx, &desc.UploadTextRequest{Text: text, Metadata: info})
	if err != nil {
		return "", fmt.Errorf("error uploading text: %w", err)
	}

	return resp.Etag, nil
}

func (s *ClientService) SendFile(addr string, filePath string, batchSize int, id, info string) (string, error) {
	creds, err := credentials.NewClientTLSFromFile("certs/server.crt", "")
	if err != nil {
		logger.Error("failed to load TLS certificates: %w", zap.Error(err))

		return "", fmt.Errorf("failed to load TLS certificates: %w", err)
	}

	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(creds))
	if err != nil {
		return "", fmt.Errorf("error dialing client: %w", err)
	}
	defer conn.Close()

	s.client = desc.NewUploadV1Client(conn)

	ss, err := session.LoadSession()
	if err != nil {
		return "", fmt.Errorf("error loading session: %w", err)
	}

	md := metadata.Pairs(login, ss.Login, "id", id, "authorization", "Bearer "+ss.Token)
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	etag, err := s.uploadFile(ctx, filePath, batchSize, info)
	if err != nil {
		return "", err
	}

	return etag, nil
}

func (s *ClientService) uploadFile(ctx context.Context, filepath string, batchSize int, info string) (string, error) {
	stream, err := s.client.UploadFile(ctx)
	if err != nil {
		return "", fmt.Errorf("error uploading file: %w", err)
	}

	file, err := os.Open(filepath)
	if err != nil {
		return "", fmt.Errorf("error opening file: %w", err)
	}
	buf := make([]byte, batchSize)
	batchNumber := 1
	for {
		num, err := file.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", fmt.Errorf("error reading buf: %w", err)
		}
		chunk := buf[:num]

		if err := stream.Send(&desc.UploadFileRequest{FileName: filepath, Chunk: chunk, Metadata: info}); err != nil {
			return "", fmt.Errorf("error uploading bytes: %w", err)
		}

		batchNumber += 1
	}
	resp, err := stream.CloseAndRecv()
	if err != nil {
		return "", fmt.Errorf("error closing stream: %w", err)
	}

	return resp.Etag, nil
}
