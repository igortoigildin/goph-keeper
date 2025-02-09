package upload

import (
	"context"
	"errors"
	"io"
	"log"
	"path/filepath"
	"strings"

	rep "github.com/igortoigildin/goph-keeper/internal/server/storage"
	fl "github.com/igortoigildin/goph-keeper/pkg/file"
	"github.com/igortoigildin/goph-keeper/pkg/logger"
	desc "github.com/igortoigildin/goph-keeper/pkg/upload_v1"
	"go.uber.org/zap"
	"google.golang.org/grpc/metadata"
)

type UploadService struct {
	dataRepository rep.DataRepository
	accessRepository rep.AccessRepository
}

func New(ctx context.Context, dataRep rep.DataRepository, accessRep rep.AccessRepository) *UploadService {
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
	login := md["login"][0]
	id := md["id"][0]
	// remove @ since this charac is not allowed for Minio bucket name
	login = strings.Replace(login, "@", "", -1)

	err := f.accessRepository.SaveAccess(ctx, login, id)
	if err != nil {
		logger.Error("error while saving access: ", zap.Error(err))

		return err
	}

	err = f.dataRepository.SaveTextData(ctx, data, login, id)
	if err != nil {
		logger.Error("error while saving bank data")

		return err
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
	login := md["login"][0]
	id := md["id"][0]
	// remove @ since this charac is not allowed for Minio bucket name
	login = strings.Replace(login, "@", "", -1)

	err := f.accessRepository.SaveAccess(ctx, login, id)
	if err != nil {
		logger.Error("error while saving access: ", zap.Error(err))

		return err
	}

	err = f.dataRepository.SaveTextData(ctx, text, login, id)
	if err != nil {
		logger.Error("error while saving text data")

		return err
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
	login := md["login"][0]
	id := md["id"][0]
	// remove @ since this charac is not allowed for Minio bucket name
	login = strings.Replace(login, "@", "", -1)

	err := f.accessRepository.SaveAccess(ctx, login, id)
	if err != nil {
		logger.Error("error while saving access: ", zap.Error(err))

		return err
	}

	err = f.dataRepository.SaveTextData(ctx, data, login, id)
	if err != nil {
		logger.Error("error while saving login&password data")

		return err
	}

	return nil
}

func (f *UploadService) SaveFile(stream desc.UploadV1_UploadFileServer) error {
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
			file.SetFile(req.GetFileName(), "client_files")
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

	md, ok := metadata.FromIncomingContext(stream.Context())
	if !ok {
		logger.Error("metadata is not received from incoming context")

		return errors.New("metadata not received from md")
	} else if md.Len() == 0 {
		logger.Error("md is emty")

		return errors.New("md is empty")
	}

	login := md["login"][0]
	id := md["id"][0]

	// remove @ since this charac is not allowed for Minio bucket name
	login = strings.Replace(login, "@", "", -1)

	err := f.accessRepository.SaveAccess(stream.Context(), login, id)
	if err != nil {
		logger.Error("error while saving access: ", zap.Error(err))

		return err
	}

	logger.Info("result:", zap.String("path", file.FilePath), zap.Any("size", fileSize))
	fileName := filepath.Base(file.FilePath)

	err = f.dataRepository.SaveFile(context.TODO(), file, login, id)
	if err != nil {
		logger.Error("error while uploading file to Minio: ", zap.Error(err))

		return err
	}

	return stream.SendAndClose(&desc.UploadFileResponse{FileName: fileName, Size: fileSize})
}
