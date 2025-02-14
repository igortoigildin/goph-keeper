package download

import (
	"context"
	"fmt"
	"io"
	"log"

	desc "github.com/igortoigildin/goph-keeper/pkg/download_v1"
	fl "github.com/igortoigildin/goph-keeper/pkg/file"

	"github.com/igortoigildin/goph-keeper/pkg/logger"
	"github.com/igortoigildin/goph-keeper/pkg/session"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
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

func New() Downloader {
	return &ClientService{}
}

func (s *ClientService) DownloadPassword(addr, id string) error {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()

	s.client = desc.NewDownloadV1Client(conn)
	ss, err := session.LoadSession()

	md := metadata.Pairs("login", ss.Login)

	ctx := metadata.NewOutgoingContext(context.Background(), md)

	resp, err := s.client.DownloadPassword(ctx, &desc.DownloadPasswordRequest{Uuid: id})
	if err != nil {
		logger.Fatal("error while receiving password", zap.Error(err))
	}

	data := resp.GetData()

	fmt.Println(data)

	logger.Info("Your data: ", zap.Any("login", data["login"]), zap.Any("password", data["password"]))

	return nil
}

func (s *ClientService) DownloadText(addr, id string) error {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()

	s.client = desc.NewDownloadV1Client(conn)
	ss, err := session.LoadSession()

	md := metadata.Pairs("login", ss.Login)

	ctx := metadata.NewOutgoingContext(context.Background(), md)

	resp, err := s.client.DownloadText(ctx, &desc.DownloadTextRequest{Uuid: id})
	if err != nil {
		logger.Fatal("error while receiving text", zap.Error(err))
	}

	data := resp.GetText()

	logger.Info("Your data: ", zap.Any("text", data))

	return nil
}

func (s *ClientService) DownloadFile(addr string, id, fileName string) error {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()

	s.client = desc.NewDownloadV1Client(conn)
	ss, err := session.LoadSession()

	md := metadata.Pairs("login", ss.Login)

	ctx := metadata.NewOutgoingContext(context.Background(), md)
	stream, err := s.client.DownloadFile(ctx, &desc.DownloadFileRequest{Uuid: id})
	if err != nil {
		logger.Fatal("error while receiving file", zap.Error(err))
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
			logger.Error("error", zap.Error(err))
			return err
		}

		chunk := req.GetChunk()
		fileSize += uint32(len(chunk))
		log.Printf("received a chunk with size: %d\n", fileSize)

		if err := file.Write(chunk); err != nil {
			log.Println(err)
			return err
		}
	}

	return nil
}

func (s *ClientService) DownloadBankDetails(addr, id string) error {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()

	s.client = desc.NewDownloadV1Client(conn)
	ss, err := session.LoadSession()

	md := metadata.Pairs("login", ss.Login)
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	resp, err := s.client.DownloadBankData(ctx, &desc.DownloadBankDataRequest{Uuid: id})
	if err != nil {
		logger.Fatal("error while receiving text", zap.Error(err))
	}

	data := resp.GetData()

	logger.Info("Your data: ", zap.Any("card_number: ", data["card_number"]),
		zap.Any("CVC: ", data["CVC"]),
		zap.Any("expiration_date: ", data["expiration_date"]),
	)

	return nil
}
