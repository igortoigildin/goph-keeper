package upload

import (
	"context"
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	models "github.com/igortoigildin/goph-keeper/internal/server/models"
	rep "github.com/igortoigildin/goph-keeper/internal/server/storage"
	fl "github.com/igortoigildin/goph-keeper/pkg/file"
	"github.com/igortoigildin/goph-keeper/pkg/logger"
	desc "github.com/igortoigildin/goph-keeper/pkg/upload_v1"
	"go.uber.org/zap"
	"google.golang.org/grpc/metadata"
)

const (
	login = "login"
	id    = "id"
)

type AccessRepository interface {
	GetAccess(ctx context.Context, login string, id string) (*models.FileInfo, error)
	SaveAccess(ctx context.Context, login string, id string) error
}

type UploadService struct {
	dataRepository   rep.DataRepository
	accessRepository AccessRepository
}

func New(ctx context.Context, dataRep rep.DataRepository, accessRep AccessRepository) *UploadService {
	return &UploadService{dataRepository: dataRep, accessRepository: accessRep}
}

func (f *UploadService) SaveBankData(ctx context.Context, data map[string]string) error {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		logger.Error("metada is not received from incoming context")

		return errors.New("metada not received from md")
	} else if md.Len() == 0 {
		logger.Error("metada is emty")

		return errors.New("md is empty")
	}
	login := md[login][0]
	id := md[id][0]
	// remove @ since this charac is not allowed for Minio bucket name
	login = strings.Replace(login, "@", "", -1)

	// Save information in storage about authorized user, which has right to access data.
	err := f.accessRepository.SaveAccess(ctx, login, id)
	if err != nil {
		logger.Error("error saving access: ", zap.Error(err))

		return fmt.Errorf("error saving access: %w", err)
	}

	err = f.dataRepository.SaveTextData(ctx, data, login, id)
	if err != nil {
		logger.Error("error saving bank data:", zap.Error(err))

		return fmt.Errorf("error saving bank data: %w", err)
	}

	return nil
}

func (f *UploadService) SaveText(ctx context.Context, text string) error {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		logger.Error("metada is not received from incoming context")

		return errors.New("metada not received from md")
	} else if md.Len() == 0 {
		logger.Error("metada is emty")

		return errors.New("md is empty")
	}
	login := md[login][0]
	id := md[id][0]
	// remove @ since this charac is not allowed for Minio bucket name
	login = strings.Replace(login, "@", "", -1)

	// Save information in storage about authorized user, which has right to access this data.
	err := f.accessRepository.SaveAccess(ctx, login, id)
	if err != nil {
		logger.Error("error saving access: ", zap.Error(err))

		return fmt.Errorf("error saving access: %w", err)
	}

	err = f.dataRepository.SaveTextData(ctx, text, login, id)
	if err != nil {
		logger.Error("error saving text data: ", zap.Error(err))

		return fmt.Errorf("error saving text data: %w", err)
	}

	return nil
}

func (f *UploadService) SaveLoginPassword(ctx context.Context, data map[string]string) error {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		logger.Error("metadata is not received from incoming context")

		return errors.New("metadata not received from md")
	} else if md.Len() == 0 {
		logger.Error("metadata is emty")

		return errors.New("metadata is empty")
	}
	login := md[login][0]
	id := md[id][0]
	// remove @ since this charac is not allowed for Minio bucket name
	login = strings.Replace(login, "@", "", -1)

	// Save information in storage about authorized user, which has right to access this data.
	err := f.accessRepository.SaveAccess(ctx, login, id)
	if err != nil {
		logger.Error("error saving access: ", zap.Error(err))

		return fmt.Errorf("error saving access: %w", err)
	}

	err = f.dataRepository.SaveTextData(ctx, data, login, id)
	if err != nil {
		logger.Error("error saving credentials data", zap.Error(err))

		return fmt.Errorf("error saving credentials data: %w", err)
	}

	return nil
}

func (f *UploadService) SaveFile(stream desc.UploadV1_UploadFileServer) error {
	file := fl.NewFile()
	var fileSize uint32
	fileSize = 0
	defer func() {
		if err := file.OutputFile.Close(); err != nil {
			logger.Error("error:", zap.Error(err))
		}
	}()

	for {
		req, err := stream.Recv()
		if file.FilePath == "" {
			file.SetFile(req.GetFileName(), "client_files")
		}
		if err == io.EOF {
			break
		}

		if err != nil {
			logger.Error("error", zap.Error(err))

			return fmt.Errorf("error receiveing the next request message from the client: %w", err)
		}

		chunk := req.GetChunk()
		fileSize += uint32(len(chunk))

		logger.Info("received a chunk with", zap.Any("size: ", fileSize))

		if err := file.Write(chunk); err != nil {
			logger.Error("error writing chunk to the file:", zap.Error(err))

			return fmt.Errorf("error writing chunk to the file: %d", err)
		}
	}

	md, ok := metadata.FromIncomingContext(stream.Context())
	if !ok {
		logger.Error("metadata is not received from incoming context")

		return errors.New("metadata not received from md")
	} else if md.Len() == 0 {
		logger.Error("md is emty")

		return errors.New("md is empty")
	}

	login := md[login][0]
	id := md[id][0]

	// remove @ since this charac is not allowed for Minio bucket name
	login = strings.Replace(login, "@", "", -1)

	err := f.accessRepository.SaveAccess(stream.Context(), login, id)
	if err != nil {
		logger.Error("error saving access: ", zap.Error(err))

		return fmt.Errorf("error saving access: %w", err)
	}

	logger.Info("result:", zap.String("path", file.FilePath), zap.Any("size", fileSize))
	fileName := filepath.Base(file.FilePath)

	err = f.dataRepository.SaveFile(context.TODO(), file, login, id)
	if err != nil {
		logger.Error("error uploading file to Minio: ", zap.Error(err))

		return fmt.Errorf("error uploading file to Minio: %w", err)
	}

	// Once file successfully uploaded to Minio storage, temp file in OC will be removed.
	err = file.Remove()
	if err != nil {
		logger.Error("error deleting file: ", zap.Error(err))

		return fmt.Errorf("error deleting file: %w", err)
	}

	return stream.SendAndClose(&desc.UploadFileResponse{FileName: fileName, Size: fileSize})
}
