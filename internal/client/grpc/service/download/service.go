package download

import (
	"context"
	"fmt"
	"io"

	desc "github.com/igortoigildin/goph-keeper/pkg/download_v1"
	"github.com/igortoigildin/goph-keeper/pkg/encryption"
	fl "github.com/igortoigildin/goph-keeper/pkg/file"
	"github.com/igortoigildin/goph-keeper/pkg/logger"
	"github.com/spf13/viper"

	"github.com/igortoigildin/goph-keeper/pkg/session"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
)

const (
	login    = "login"
	password = "password"
)

type Downloader interface {
	DownloadPassword(addr, id string) error
	DownloadText(addr, id string) error
	DownloadFile(addr, id, fileName string) error
	DownloadBankDetails(addr, id string) error
}

type ClientService struct {
	client desc.DownloadV1Client
}

func New() *ClientService {
	return &ClientService{}
}

func (s *ClientService) DownloadPassword(addr, id string) error {
	creds, err := credentials.NewClientTLSFromFile("certs/server.crt", "")
	if err != nil {
		logger.Error("failed to load TLS certificates: %w", zap.Error(err))
		return fmt.Errorf("failed to load TLS certificates: %w", err)
	}

	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(creds))
	if err != nil {
		return fmt.Errorf("error dialing client: %w", err)
	}
	defer conn.Close()

	s.client = desc.NewDownloadV1Client(conn)
	ss, err := session.LoadSession()
	if err != nil {
		return fmt.Errorf("error loading session: %w", err)
	}

	md := metadata.Pairs(login, ss.Login)

	ctx := metadata.NewOutgoingContext(context.Background(), md)

	resp, err := s.client.DownloadPassword(ctx, &desc.DownloadPasswordRequest{Uuid: id})
	if err != nil {
		return fmt.Errorf("error downloading password: %w", err)
	}

	data := resp.GetData()
	metadata := resp.GetMetadata()

	decryptedLogin, err := encryption.Decrypt(data["login"], []byte(viper.Get("ENCRYPTION_KEY").(string)))
	if err != nil {
		logger.Error("failed to decrypt login", zap.Error(err))
	}

	decryptedPassword, err := encryption.Decrypt(data["password"], []byte(viper.Get("ENCRYPTION_KEY").(string)))
	if err != nil {
		logger.Error("failed to decrypt password", zap.Error(err))
	}

	logger.Info("Your data: ", zap.Any("login", decryptedLogin), zap.Any("password", decryptedPassword),
		zap.Any("info: ", metadata),
	)

	return nil
}

func (s *ClientService) DownloadText(addr, id string) error {
	creds, err := credentials.NewClientTLSFromFile("certs/server.crt", "")
	if err != nil {
		logger.Error("failed to load TLS certificates: %w", zap.Error(err))

		return fmt.Errorf("failed to load TLS certificates: %w", err)
	}

	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(creds))
	if err != nil {
		return fmt.Errorf("error dialing client: %w", err)
	}
	defer conn.Close()

	s.client = desc.NewDownloadV1Client(conn)
	ss, err := session.LoadSession()
	if err != nil {
		return fmt.Errorf("error loading session: %w", err)
	}

	md := metadata.Pairs(login, ss.Login)

	ctx := metadata.NewOutgoingContext(context.Background(), md)

	resp, err := s.client.DownloadText(ctx, &desc.DownloadTextRequest{Uuid: id})
	if err != nil {
		return fmt.Errorf("error downloading text: %w", err)
	}

	dataEncrypted := resp.GetText()
	metadata := resp.GetMetadata()

	decryptedText, err := encryption.Decrypt(dataEncrypted, []byte(viper.Get("ENCRYPTION_KEY").(string)))
	if err != nil {
		logger.Error("failed to decrypt text data", zap.Error(err))
	}

	logger.Info("Your data: ", zap.Any("text", decryptedText), zap.Any("metadata", metadata))

	return nil
}

func (s *ClientService) DownloadFile(addr string, id, fileName string) error {
	creds, err := credentials.NewClientTLSFromFile("certs/server.crt", "")
	if err != nil {
		logger.Error("failed to load TLS certificates: %w", zap.Error(err))

		return fmt.Errorf("failed to load TLS certificates: %w", err)
	}

	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(creds))
	if err != nil {
		return fmt.Errorf("error dialing client: %w", err)
	}
	defer conn.Close()

	s.client = desc.NewDownloadV1Client(conn)
	ss, err := session.LoadSession()
	if err != nil {
		return fmt.Errorf("error loading session: %w", err)
	}

	md := metadata.Pairs(login, ss.Login)

	ctx := metadata.NewOutgoingContext(context.Background(), md)
	stream, err := s.client.DownloadFile(ctx, &desc.DownloadFileRequest{Uuid: id})
	if err != nil {
		return fmt.Errorf("error downloading file: %w", err)
	}

	file := fl.NewFile()
	var fileSize uint32
	fileSize = 0
	defer func() {
		if err := file.OutputFile.Close(); err != nil {
			logger.Error("error", zap.Error(err))
		}
	}()

	for {
		req, err := stream.Recv()
		if file.FilePath == "" {
			file.SetFile(fileName, "client_files")
		}
		if err == io.EOF {
			break
		}

		if err != nil {
			return fmt.Errorf("error receiving byte chunk: %w", err)
		}

		chunk := req.GetChunk()
		fileSize += uint32(len(chunk))
		logger.Info("received a chunk with size:", zap.Uint32("size", fileSize))

		if err := file.Write(chunk); err != nil {
			return fmt.Errorf("error adding byte chunk to file: %w", err)
		}
	}

	return nil
}

func (s *ClientService) DownloadBankDetails(addr, id string) error {
	creds, err := credentials.NewClientTLSFromFile("certs/server.crt", "")
	if err != nil {
		logger.Error("failed to load TLS certificates: %w", zap.Error(err))

		return fmt.Errorf("failed to load TLS certificates: %w", err)
	}

	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(creds))
	if err != nil {
		return fmt.Errorf("error dialing client: %w", err)
	}
	defer conn.Close()

	s.client = desc.NewDownloadV1Client(conn)
	ss, err := session.LoadSession()
	if err != nil {
		return fmt.Errorf("error loading session: %w", err)
	}

	md := metadata.Pairs(login, ss.Login)
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	resp, err := s.client.DownloadBankData(ctx, &desc.DownloadBankDataRequest{Uuid: id})
	if err != nil {
		return fmt.Errorf("erorr downloading text: %w", err)
	}

	data := resp.GetData()
	metadata := resp.GetMetadata()

	decryptedCardNumber, err := encryption.Decrypt(data["card_number"], []byte(viper.Get("ENCRYPTION_KEY").(string)))
	if err != nil {
		logger.Error("failed to decrypt card number", zap.Error(err))
	}

	decryptedCVC, err := encryption.Decrypt(data["CVC"], []byte(viper.Get("ENCRYPTION_KEY").(string)))
	if err != nil {
		logger.Error("failed to decrypt cvc", zap.Error(err))
	}

	decryptedExpDate, err := encryption.Decrypt(data["expiration_date"], []byte(viper.Get("ENCRYPTION_KEY").(string)))
	if err != nil {
		logger.Error("failed to decrypt expiration date", zap.Error(err))
	}

	logger.Info("Your data: ", zap.Any("card_number: ", decryptedCardNumber),
		zap.Any("CVC: ", decryptedCVC),
		zap.Any("expiration_date: ", decryptedExpDate),
		zap.Any("metadata: ", metadata),
	)

	return nil
}
