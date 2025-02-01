package service

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

type UploadService struct{
	dataRepository rep.DataRepository
}

func New(ctx context.Context, rep rep.DataRepository) *UploadService {
	return &UploadService{dataRepository: rep}
}

func (f *UploadService) Upload(stream desc.UploadV1_UploadFileServer) error {
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

	var login string
	md, ok := metadata.FromIncomingContext(stream.Context())
	if !ok {
		logger.Error("login is not recieved from incoming context")
		return errors.New("login not recieved from md")
	}

	login = md["login"][0]
	if len(login) == 0 {
		logger.Error("md-login is emty")
		return errors.New("md is empty")
	} 

	// remove @ since this charac is not allowed for Minio bucket name
	login = strings.Replace(login, "@", "", -1)

	logger.Info("result:", zap.String("path", file.FilePath), zap.Any("size", fileSize))
	fileName := filepath.Base(file.FilePath)

	err := f.dataRepository.SaveData(context.TODO(), file, login)
	if err != nil {
		logger.Error("error while uploading file to Minio: ", zap.Error(err))

		return err
	}

	return stream.SendAndClose(&desc.UploadFileResponse{FileName: fileName, Size: fileSize})
}
