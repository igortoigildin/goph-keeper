package download

import (
	"context"
	"encoding/json"
	"errors"

	rep "github.com/igortoigildin/goph-keeper/internal/server/storage"
	"github.com/igortoigildin/goph-keeper/pkg/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc/metadata"
)

type DownloadService struct {
	dataRepository rep.DataRepository
	accessRepository rep.AccessRepository
}

func New(ctx context.Context, dataRep rep.DataRepository, accessRep rep.AccessRepository) *DownloadService {
	return &DownloadService{dataRepository: dataRep, accessRepository: accessRep}
}

func (d *DownloadService) DownloadFile(ctx context.Context, id string) ([]byte, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		logger.Error("metada is not received from incoming context")

		return nil, errors.New("metada not received from md")
	} else if md.Len() == 0 {
		logger.Error("metada is emty")

		return nil, errors.New("md is empty")
	}
	login := md["login"][0]

	file, err := d.dataRepository.DownloadFile(ctx, login, id)
	if err != nil {
		return nil, err
	}

	// f, err := os.Create(filePath)
	// if err != nil {
	// 	logger.Error("error creating file:", zap.Error(err))
	// }
	// defer f.Close()

	// _, err = io.Copy(f, file)
	// if err != nil {
	// 	logger.Error("error writing to file:", zap.Error(err))
	// }

	logger.Info("File successfully downloaded from Minio bucket ")

	return file.Bytes(), nil
}

func (d *DownloadService) DownloadBankData(ctx context.Context, id string) (map[string]string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		logger.Error("metada is not received from incoming context")

		return nil, errors.New("metada not received from md")
	} else if md.Len() == 0 {
		logger.Error("metada is emty")

		return nil, errors.New("md is empty")
	}
	login := md["login"][0]

	data, err := d.dataRepository.DownloadTextData(ctx, login, id)
	if err != nil {
		logger.Error("error downloading bank details: ", zap.Error(err))
	}

	res := make(map[string]string, 3)

	err = json.Unmarshal(data, &res)
	if err != nil {
		logger.Error("error unmarshalling JSON:", zap.Error(err))
	}

	return res, nil
}

func (d *DownloadService) DownloadText(ctx context.Context, id string) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		logger.Error("metada is not received from incoming context")

		return "", errors.New("metada not received from md")
	} else if md.Len() == 0 {
		logger.Error("metada is emty")

		return "", errors.New("md is empty")
	}
	login := md["login"][0]

	data, err := d.dataRepository.DownloadTextData(ctx, login, id)
	if err != nil {
		logger.Error("error downloading text data: ", zap.Error(err))
	}

	return string(data), nil
}

func (d *DownloadService) DownloadLoginPassword(ctx context.Context, id string) (map[string]string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		logger.Error("metada is not received from incoming context")

		return nil, errors.New("metada not received from md")
	} else if md.Len() == 0 {
		logger.Error("metada is emty")

		return nil, errors.New("md is empty")
	}
	login := md["login"][0]

	data, err := d.dataRepository.DownloadTextData(ctx, login, id)
	if err != nil {
		logger.Error("error downloading login details: ", zap.Error(err))
	}

	res := make(map[string]string, 3)

	err = json.Unmarshal(data, &res)
	if err != nil {
		logger.Error("error unmarshalling JSON:", zap.Error(err))
	}

	return res, nil
}

