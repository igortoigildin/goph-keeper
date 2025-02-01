package upload

import (
	"context"
	"io"
	"os"
	"os/signal"
	"syscall"

	"github.com/igortoigildin/goph-keeper/pkg/logger"
	"github.com/igortoigildin/goph-keeper/pkg/session"
	desc "github.com/igortoigildin/goph-keeper/pkg/upload_v1"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type Sender interface {
	SendPassword(addr, loginStr, passStr string) error
	SendText()
	SendFile(addr string, filePath string, batchSize int) error
	SendBankDetails()
}

type ClientService struct {
	client    desc.UploadV1Client
}

func New() Sender {
	return ClientService{}
}


func (s *ClientService) SendPassword(addr, loginStr, passStr string) error {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()

	s.client = desc.NewUploadV1Client(conn)

	interrupt := make(chan os.Signal, 1)
	shutdownSignals := []os.Signal{
		os.Interrupt,
		syscall.SIGTERM,
		syscall.SIGINT,
		syscall.SIGQUIT,
	}
	signal.Notify(interrupt, shutdownSignals...)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ss, err := session.LoadSession()

	md := metadata.Pairs("login", ss.Login)

	ctx = metadata.NewOutgoingContext(context.Background(), md)

	go func(s *ClientService) {
		if err = s.uploadFile(ctx, filePath, batchSize, cancel); err != nil {
			logger.Fatal("error while sending file", zap.Error(err))
			cancel()
		}
	}(s)

	select {
	case killSignal := <-interrupt:
		logger.Info("Got ", zap.Any("signal", killSignal))
		cancel()
	case <-ctx.Done():
	}
	return nil
}


func (s *ClientService) uploadPassword(ctx context.Context, loginStr, passStr string, cancel context.CancelFunc) error {
	data := make(map[string]string, 1)
	data[loginStr] = passStr

	_, err := s.client.UploadPassword(ctx, &desc.UploadPasswordRequest{Data: data})
	if err != nil {
		logger.Error("error", zap.Error(err))
		return err
	}
	// TODO : to be completed

	return nil
}



// type ClientService struct {
// 	addr      string
// 	filePath  string
// 	batchSize int
// 	client    desc.UploadV1Client
// 	email     string
// }

// func New(addr string, filePath string, batchSize int) *ClientService {
// 	return &ClientService{
// 		addr:      addr,
// 		filePath:  filePath,
// 		batchSize: batchSize,
// 	}
// }

func (s *ClientService) SendFile(addr string, filePath string, batchSize int) error {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()

	s.client = desc.NewUploadV1Client(conn)

	interrupt := make(chan os.Signal, 1)
	shutdownSignals := []os.Signal{
		os.Interrupt,
		syscall.SIGTERM,
		syscall.SIGINT,
		syscall.SIGQUIT,
	}
	signal.Notify(interrupt, shutdownSignals...)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ss, err := session.LoadSession()

	md := metadata.Pairs("login", ss.Login)

	ctx = metadata.NewOutgoingContext(context.Background(), md)

	go func(s *ClientService) {
		if err = s.uploadFile(ctx, filePath, batchSize, cancel); err != nil {
			logger.Fatal("error while sending file", zap.Error(err))
			cancel()
		}
	}(s)

	select {
	case killSignal := <-interrupt:
		logger.Info("Got ", zap.Any("signal", killSignal))
		cancel()
	case <-ctx.Done():
	}
	return nil
}

func (s *ClientService) uploadFile(ctx context.Context, filepath string, batchSize int, cancel context.CancelFunc) error {
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

	cancel()

	return nil
}
