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

type ClientService struct {
	addr      string
	filePath  string
	batchSize int
	client    desc.UploadV1Client
	email 	  string
}

func New(addr string, filePath string, batchSize int) *ClientService {
	return &ClientService{
		addr:      addr,
		filePath:  filePath,
		batchSize: batchSize,
	}
}

func (s *ClientService) SendFile() error {
	conn, err := grpc.Dial(s.addr, grpc.WithInsecure())
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
	// s.email = ss.Email

	md := metadata.Pairs("email", ss.Email)
	ctx = metadata.NewOutgoingContext(context.Background(), md)

	go func(s *ClientService) {
		if err = s.upload(ctx, cancel); err != nil {
			logger.Fatal("error while sending file", zap.Error(err))
			cancel()
		}
	}(s)

	select {
	case killSignal := <-interrupt:
		logger.Info("Got ", zap.Any("singal", killSignal))
		cancel()
	case <-ctx.Done():
	}
	return nil
}

func (s *ClientService) upload(ctx context.Context, cancel context.CancelFunc) error {
	stream, err := s.client.UploadFile(ctx)
	if err != nil {
		logger.Error("error", zap.Error(err))
		return err
	}

	file, err := os.Open(s.filePath)
	if err != nil {
		logger.Error("error", zap.Error(err))
		return err
	}
	buf := make([]byte, s.batchSize)
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

		if err := stream.Send(&desc.UploadFileRequest{FileName: s.filePath, Chunk: chunk}); err != nil {
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