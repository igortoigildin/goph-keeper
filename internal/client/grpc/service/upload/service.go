package upload

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/igortoigildin/goph-keeper/pkg/session"
	desc "github.com/igortoigildin/goph-keeper/pkg/upload_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

const (
	login    = "login"
	password = "password"
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

func New() *ClientService {
	return &ClientService{}
}

func (s *ClientService) SendPassword(addr, loginStr, passStr string, id string, meta string) error {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("error dialing client: %w", err)
	}
	defer conn.Close()

	s.client = desc.NewUploadV1Client(conn)

	ss, err := session.LoadSession()
	if err != nil {
		return fmt.Errorf("error loading session: %w", err)
	}

	md := metadata.Pairs(login, ss.Login, "id", id)

	ctx := metadata.NewOutgoingContext(context.Background(), md)

	if err = s.uploadPassword(ctx, loginStr, passStr, meta); err != nil {

		return err
	}

	return nil
}

func (s *ClientService) uploadPassword(ctx context.Context, loginStr, passStr, meta string) error {
	data := make(map[string]string, 2)
	data[login] = loginStr
	data[password] = passStr

	_, err := s.client.UploadPassword(ctx, &desc.UploadPasswordRequest{Data: data, Metadata: meta})
	if err != nil {
		return fmt.Errorf("error uploading credentials: %w", err)
	}

	return nil
}

func (s *ClientService) SendBankDetails(addr, cardNumber, cvc, expDate string, id, meta string) error {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("error dialing client: %w", err)
	}
	defer conn.Close()

	s.client = desc.NewUploadV1Client(conn)

	ss, err := session.LoadSession()
	if err != nil {
		return fmt.Errorf("error loading session: %w", err)
	}

	md := metadata.Pairs(login, ss.Login, "id", id)

	ctx := metadata.NewOutgoingContext(context.Background(), md)

	if err = s.uploadBankDetails(ctx, cardNumber, cvc, expDate, meta); err != nil {
		return err
	}

	return nil
}

func (s *ClientService) uploadBankDetails(ctx context.Context, cardNumber, cvc, expDate, meta string) error {
	data := make(map[string]string, 3)
	data["card_number"] = cardNumber
	data["CVC"] = cvc
	data["expiration_date"] = expDate

	_, err := s.client.UploadBankData(ctx, &desc.UploadBankDataRequest{Data: data, Metadata: meta})
	if err != nil {
		return fmt.Errorf("error uploading bank details: %w", err)
	}

	return nil
}

func (s *ClientService) SendText(addr, text string, id string, meta string) error {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("error dialing client: %w", err)
	}
	defer conn.Close()

	s.client = desc.NewUploadV1Client(conn)

	ss, err := session.LoadSession()
	if err != nil {
		return fmt.Errorf("error loading session: %w", err)
	}

	md := metadata.Pairs(login, ss.Login, "id", id)
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	if err = s.uploadText(ctx, text, meta); err != nil {
		return err
	}

	return nil
}

func (s *ClientService) uploadText(ctx context.Context, text, meta string) error {
	_, err := s.client.UploadText(ctx, &desc.UploadTextRequest{Text: text, Metadata: meta})
	if err != nil {
		return fmt.Errorf("error uploading text: %w", err)
	}

	return nil
}

func (s *ClientService) SendFile(addr string, filePath string, batchSize int, id, meta string) error {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("error dialing client: %w", err)
	}
	defer conn.Close()

	s.client = desc.NewUploadV1Client(conn)

	ss, err := session.LoadSession()
	if err != nil {
		return fmt.Errorf("error loading session: %w", err)
	}

	md := metadata.Pairs(login, ss.Login, "id", id)
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	if err = s.uploadFile(ctx, filePath, batchSize, meta); err != nil {
		return err
	}

	return nil
}

func (s *ClientService) uploadFile(ctx context.Context, filepath string, batchSize int, meta string) error {
	stream, err := s.client.UploadFile(ctx)
	if err != nil {
		return fmt.Errorf("error uploading file: %w", err)
	}

	file, err := os.Open(filepath)
	if err != nil {
		return fmt.Errorf("error opening file: %w", err)
	}
	buf := make([]byte, batchSize)
	batchNumber := 1
	for {
		num, err := file.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("error reading buf: %w", err)
		}
		chunk := buf[:num]

		if err := stream.Send(&desc.UploadFileRequest{FileName: filepath, Chunk: chunk, Metadata: meta}); err != nil {
			return fmt.Errorf("error uploading bytes: %w", err)
		}

		batchNumber += 1
	}
	_, err = stream.CloseAndRecv()
	if err != nil {
		return fmt.Errorf("error closing stream: %w", err)
	}

	return nil
}
