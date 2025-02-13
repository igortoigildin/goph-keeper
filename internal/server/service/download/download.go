package download

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	models "github.com/igortoigildin/goph-keeper/internal/server/models"
	rep "github.com/igortoigildin/goph-keeper/internal/server/storage"
	"github.com/igortoigildin/goph-keeper/pkg/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc/metadata"
)

type AccessRepository interface {
	GetAccess(ctx context.Context, login string, id string) (*models.FileInfo, error)
	SaveAccess(ctx context.Context, login string, id string) error
}

const (
	login = "login"
)

type DownloadService struct {
	dataRepository   rep.DataRepository
	accessRepository AccessRepository
}

func New(ctx context.Context, dataRep rep.DataRepository, accessRep AccessRepository) *DownloadService {
	return &DownloadService{dataRepository: dataRep, accessRepository: accessRep}
}

// DownloadFile checks whether user is authorized to download file with certain id,
// if so, downloading begins from storage and file is being returned, if not - returns error.
func (d *DownloadService) DownloadFile(ctx context.Context, id string) ([]byte, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		logger.Error("metadata is not received from incoming context")

		return nil, errors.New("metada not received from md")
	} else if md.Len() == 0 {
		logger.Error("metadata is emty")

		return nil, errors.New("md is empty")
	}
	login := md[login][0]

	// get metadata about file with provided id
	fileInfo, err := d.accessRepository.GetAccess(ctx, login, id)
	if err != nil {
		logger.Error("failed to get access for file")

		return nil, fmt.Errorf("error getting access for specific file from repo: %w", err)
	}

	// check whether user is authorized to get access to this specific file
	if fileInfo.Id != login {
		logger.Info("Authorization error")

		return nil, fmt.Errorf("authorization error: %w", err)
	}

	file, err := d.dataRepository.DownloadFile(ctx, login, id)
	if err != nil {
		return nil, fmt.Errorf("error downloading file from repository: %w", err)
	}

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

	// get metadata about file with provided id
	fileInfo, err := d.accessRepository.GetAccess(ctx, login, id)
	if err != nil {
		logger.Error("failed to get access for file")

		return nil, fmt.Errorf("error getting access for specific file from repo: %w", err)
	}

	// check whether user is authorized to get access to this specific file
	if fileInfo.Id != login {
		logger.Info("Authorization error")

		return nil, fmt.Errorf("authorization error: %w", err)
	}

	data, err := d.dataRepository.DownloadTextData(ctx, login, id)
	if err != nil {
		return nil, fmt.Errorf("error downloading bank details: %s", err)
	}

	res := make(map[string]string, 3)

	err = json.Unmarshal(data, &res)
	if err != nil {
		logger.Error("error unmarshalling JSON:", zap.Error(err))

		return nil, fmt.Errorf("error unmarshalling JSON: %w", err)
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

	// get metadata about file with provided id
	fileInfo, err := d.accessRepository.GetAccess(ctx, login, id)
	if err != nil {
		logger.Error("failed to get access for file")

		return "", fmt.Errorf("error getting access for specific file from repo: %w", err)
	}

	// check whether user is authorized to get access to this specific file
	if fileInfo.Id != login {
		logger.Info("Authorization error")

		return "", fmt.Errorf("authorization error: %w", err)
	}

	data, err := d.dataRepository.DownloadTextData(ctx, login, id)
	if err != nil {
		return "", fmt.Errorf("error downloading text: %s", err)
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

	// get metadata about file with provided id
	fileInfo, err := d.accessRepository.GetAccess(ctx, login, id)
	if err != nil {
		logger.Error("failed to get access for file")

		return nil, fmt.Errorf("error getting access for specific file from repo: %w", err)
	}

	// check whether user is authorized to get access to this specific file
	if fileInfo.Id != login {
		logger.Info("Authorization error")

		return nil, fmt.Errorf("authorization error: %w", err)
	}

	data, err := d.dataRepository.DownloadTextData(ctx, login, id)
	if err != nil {
		logger.Error("error downloading login credentials: ", zap.Error(err))

		return nil, fmt.Errorf("error downloading login credentials: %s", err)
	}

	res := make(map[string]string, 3)

	err = json.Unmarshal(data, &res)
	if err != nil {
		logger.Error("error unmarshalling JSON:", zap.Error(err))

		return nil, fmt.Errorf("error unmarshalling JSON: %w", err)
	}

	return res, nil
}
