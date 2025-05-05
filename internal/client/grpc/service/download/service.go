package download

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/igortoigildin/goph-keeper/internal/client/grpc/models"
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

type ClientService struct {
	client desc.DownloadV1Client
}

func New() *ClientService {
	return &ClientService{}
}

func (s *ClientService) DownloadPassword(addr, id string) (models.Credential, error) {
	creds, err := credentials.NewClientTLSFromFile("certs/server.crt", "")
	if err != nil {
		logger.Error("failed to load TLS certificates: %w", zap.Error(err))
		return models.Credential{}, fmt.Errorf("failed to load TLS certificates: %w", err)
	}

	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(creds))
	if err != nil {
		return models.Credential{}, fmt.Errorf("error dialing client: %w", err)
	}
	defer conn.Close()

	s.client = desc.NewDownloadV1Client(conn)

	ss, err := session.LoadSession()
	if err != nil {
		return models.Credential{}, fmt.Errorf("error loading session: %w", err)
	}

	md := metadata.Pairs(login, ss.Login, "authorization", "Bearer "+ss.Token)

	ctx := metadata.NewOutgoingContext(context.Background(), md)

	resp, err := s.client.DownloadPassword(ctx, &desc.DownloadPasswordRequest{Uuid: id})
	if err != nil {
		return models.Credential{}, fmt.Errorf("error downloading password: %w", err)
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

	resObj := models.Credential{
		ID:       id,
		Username: decryptedLogin,
		Password: decryptedPassword,
		Service:  metadata,
	}

	return resObj, nil
}

func (s *ClientService) DownloadText(addr, id string) (models.Text, error) {
	creds, err := credentials.NewClientTLSFromFile("certs/server.crt", "")
	if err != nil {
		logger.Error("failed to load TLS certificates: %w", zap.Error(err))

		return models.Text{}, fmt.Errorf("failed to load TLS certificates: %w", err)
	}

	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(creds))
	if err != nil {
		return models.Text{}, fmt.Errorf("error dialing client: %w", err)
	}
	defer conn.Close()

	s.client = desc.NewDownloadV1Client(conn)
	ss, err := session.LoadSession()
	if err != nil {
		return models.Text{}, fmt.Errorf("error loading session: %w", err)
	}

	md := metadata.Pairs(login, ss.Login, "authorization", "Bearer "+ss.Token)

	ctx := metadata.NewOutgoingContext(context.Background(), md)

	resp, err := s.client.DownloadText(ctx, &desc.DownloadTextRequest{Uuid: id})
	if err != nil {
		return models.Text{}, fmt.Errorf("error downloading text: %w", err)
	}

	dataEncrypted := resp.GetText()
	metadata := resp.GetMetadata()

	decryptedText, err := encryption.Decrypt(dataEncrypted, []byte(viper.Get("ENCRYPTION_KEY").(string)))
	if err != nil {
		logger.Error("failed to decrypt text data", zap.Error(err))
	}

	logger.Info("Your data: ", zap.Any("text", decryptedText), zap.Any("metadata", metadata))

	resObj := models.Text{
		ID:   id,
		Text: decryptedText,
		Info: metadata,
	}

	return resObj, nil
}

func (s *ClientService) DownloadFile(addr string, id, fileName string) (models.File, error) {
	creds, err := credentials.NewClientTLSFromFile("certs/server.crt", "")
	if err != nil {
		logger.Error("failed to load TLS certificates: %w", zap.Error(err))

		return models.File{}, fmt.Errorf("failed to load TLS certificates: %w", err)
	}

	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(creds))
	if err != nil {
		return models.File{}, fmt.Errorf("error dialing client: %w", err)
	}
	defer conn.Close()

	s.client = desc.NewDownloadV1Client(conn)
	ss, err := session.LoadSession()
	if err != nil {
		return models.File{}, fmt.Errorf("error loading session: %w", err)
	}

	md := metadata.Pairs(login, ss.Login, "authorization", "Bearer "+ss.Token)

	ctx := metadata.NewOutgoingContext(context.Background(), md)
	stream, err := s.client.DownloadFile(ctx, &desc.DownloadFileRequest{Uuid: id})
	if err != nil {
		return models.File{}, fmt.Errorf("error downloading file: %w", err)
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
			return models.File{}, fmt.Errorf("error receiving byte chunk: %w", err)
		}

		chunk := req.GetChunk()
		fileSize += uint32(len(chunk))
		logger.Info("received a chunk with size:", zap.Uint32("size", fileSize))

		if err := file.Write(chunk); err != nil {
			return models.File{}, fmt.Errorf("error adding byte chunk to file: %w", err)
		}
	}

	fileData, err := os.ReadFile(file.FilePath)
	if err != nil {
		logger.Error("error while reading file", zap.Error(err))
		return models.File{}, fmt.Errorf("error while reading file: %w", err)
	}

	resObj := models.File{
		ID:       id,
		Data:     fileData,
		Filename: file.FilePath,
	}

	return resObj, nil
}

func (s *ClientService) DownloadBankDetails(addr, id string) (models.BankDetails, error) {
	creds, err := credentials.NewClientTLSFromFile("certs/server.crt", "")
	if err != nil {
		logger.Error("failed to load TLS certificates: %w", zap.Error(err))

		return models.BankDetails{}, fmt.Errorf("failed to load TLS certificates: %w", err)
	}

	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(creds))
	if err != nil {
		return models.BankDetails{}, fmt.Errorf("error dialing client: %w", err)
	}
	defer conn.Close()

	s.client = desc.NewDownloadV1Client(conn)
	ss, err := session.LoadSession()
	if err != nil {
		return models.BankDetails{}, fmt.Errorf("error loading session: %w", err)
	}

	md := metadata.Pairs(login, ss.Login, "authorization", "Bearer "+ss.Token)
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	resp, err := s.client.DownloadBankData(ctx, &desc.DownloadBankDataRequest{Uuid: id})
	if err != nil {
		return models.BankDetails{}, fmt.Errorf("erorr downloading text: %w", err)
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

	resObj := models.BankDetails{
		ID:         id,
		CardNumber: decryptedCardNumber,
		Cvc:        decryptedCVC,
		ExpDate:    decryptedExpDate,
		Info:       metadata,
	}

	return resObj, nil
}
