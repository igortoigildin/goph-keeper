package download

import (
	"context"
	"fmt"
	"io"

	desc "github.com/igortoigildin/goph-keeper/pkg/download_v1"
	fl "github.com/igortoigildin/goph-keeper/pkg/file"
	"github.com/igortoigildin/goph-keeper/pkg/logger"

	"github.com/igortoigildin/goph-keeper/pkg/session"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
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

	logger.Info("Your data: ", zap.Any("login", data["login"]), zap.Any("password", data["password"]))

	return nil
}

func (s *ClientService) DownloadText(addr, id string) error {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
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

	data := resp.GetText()
	metadata := resp.GetMetadata()

	logger.Info("Your data: ", zap.Any("text", data), zap.Any("metadata", metadata))

	return nil
}

func (s *ClientService) DownloadFile(addr string, id, fileName string) error {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
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
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
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

	logger.Info("Your data: ", zap.Any("card_number: ", data["card_number"]),
		zap.Any("CVC: ", data["CVC"]),
		zap.Any("expiration_date: ", data["expiration_date"]),
	)

	return nil
}
