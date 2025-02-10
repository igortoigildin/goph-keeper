package upload

import (
	"context"
	"io"
	"os"
	"sync"

	"github.com/igortoigildin/goph-keeper/pkg/logger"
	"github.com/igortoigildin/goph-keeper/pkg/session"
	desc "github.com/igortoigildin/goph-keeper/pkg/upload_v1"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type Sender interface {
	SendPassword(addr, loginStr, passStr string, id string) error
	SendText(addr, text string, id string) error
	SendFile(addr string, filePath string, batchSize int, id string) error
	SendBankDetails(addr, cardNumber, cvc, expDate string, id string) error
}

type ClientService struct {
	client desc.UploadV1Client
}

func New() Sender {
	return &ClientService{}
}

func (s *ClientService) SendPassword(addr, loginStr, passStr string, id string) error {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()

	s.client = desc.NewUploadV1Client(conn)

	var wg sync.WaitGroup
	ss, err := session.LoadSession()

	md := metadata.Pairs("login", ss.Login, "id", id)

	ctx := metadata.NewOutgoingContext(context.Background(), md)

	wg.Add(1)
	go func(s *ClientService) {
		if err = s.uploadPassword(ctx, loginStr, passStr, &wg); err != nil {
			logger.Fatal("error while sending file", zap.Error(err))
		}
	}(s)

	wg.Wait()

	return nil
}

func (s *ClientService) uploadPassword(ctx context.Context, loginStr, passStr string, wg *sync.WaitGroup) error {
	defer wg.Done()

	data := make(map[string]string, 2)
	data["login"] = loginStr
	data["password"] = passStr

	_, err := s.client.UploadPassword(ctx, &desc.UploadPasswordRequest{Data: data})
	if err != nil {
		logger.Error("error", zap.Error(err))
		return err
	}

	logger.Info("Login && password data sent successfully")

	return nil
}

func (s *ClientService) SendBankDetails(addr, cardNumber, cvc, expDate string, id string) error {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()

	s.client = desc.NewUploadV1Client(conn)

	ss, err := session.LoadSession()
	var wg sync.WaitGroup
	md := metadata.Pairs("login", ss.Login, "id", id)

	ctx := metadata.NewOutgoingContext(context.Background(), md)

	wg.Add(1)
	go func(s *ClientService) {
		if err = s.uploadBankDetails(ctx, cardNumber, cvc, expDate, &wg); err != nil {
			logger.Fatal("error while sending file", zap.Error(err))
		}
	}(s)

	wg.Wait()

	return nil
}

func (s *ClientService) uploadBankDetails(ctx context.Context, cardNumber, cvc, expDate string, wg *sync.WaitGroup) error {
	defer wg.Done()

	data := make(map[string]string, 3)
	data["card_number"] = cardNumber
	data["CVC"] = cvc
	data["expiration_date"] = expDate

	_, err := s.client.UploadBankData(ctx, &desc.UploadBankDataRequest{Data: data})
	if err != nil {
		logger.Error("error", zap.Error(err))
		return err
	}

	logger.Info("Login && password data sent successfully")

	return nil
}

func (s *ClientService) SendText(addr, text string, id string) error {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()

	s.client = desc.NewUploadV1Client(conn)

	var wg sync.WaitGroup
	ss, err := session.LoadSession()

	md := metadata.Pairs("login", ss.Login, "id", id)

	ctx := metadata.NewOutgoingContext(context.Background(), md)

	wg.Add(1)
	go func(s *ClientService) {
		if err = s.uploadText(ctx, text, &wg); err != nil {
			logger.Fatal("error while sending file", zap.Error(err))
		}
	}(s)

	wg.Wait()

	return nil
}

func (s *ClientService) uploadText(ctx context.Context, text string, wg *sync.WaitGroup) error {
	defer wg.Done()

	_, err := s.client.UploadText(ctx, &desc.UploadTextRequest{Text: text})
	if err != nil {
		logger.Error("error", zap.Error(err))
		return err
	}

	logger.Info("Text data sent successfully")

	return nil
}

func (s *ClientService) SendFile(addr string, filePath string, batchSize int, id string) error {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()

	s.client = desc.NewUploadV1Client(conn)
	var wg sync.WaitGroup

	ss, err := session.LoadSession()
	md := metadata.Pairs("login", ss.Login, "id", id)
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	wg.Add(1)
	go func(s *ClientService) {
		if err = s.uploadFile(ctx, filePath, batchSize, &wg); err != nil {
			logger.Fatal("error while sending file", zap.Error(err))
		}

	}(s)

	wg.Wait()

	return nil
}

func (s *ClientService) uploadFile(ctx context.Context, filepath string, batchSize int, wg *sync.WaitGroup) error {
	defer wg.Done()

	stream, err := s.client.UploadFile(ctx)
	if err != nil {
		logger.Error("error", zap.Error(err))
		return err
	}

	file, err := os.Open(filepath)
	if err != nil {
		logger.Error("error", zap.Error(err))
		return err
	}
	buf := make([]byte, batchSize)
	batchNumber := 1
	for {
		num, err := file.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		chunk := buf[:num]

		if err := stream.Send(&desc.UploadFileRequest{FileName: filepath, Chunk: chunk}); err != nil {
			logger.Error("error", zap.Error(err))
			return err
		}

		logger.Info("Sent:",
			zap.Int("batch", batchNumber),
			zap.Int("size", len(chunk)),
		)

		batchNumber += 1
	}
	res, err := stream.CloseAndRecv()
	if err != nil {
		logger.Error("error", zap.Error(err))
		return err
	}

	logger.Info("Sent:",
		zap.Int("bytes", int(res.GetSize())),
		zap.String("file name", res.GetFileName()),
	)

	return nil
}
