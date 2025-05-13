package download

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	models "github.com/igortoigildin/goph-keeper/internal/server/models"
	rep "github.com/igortoigildin/goph-keeper/internal/server/storage"
	"github.com/igortoigildin/goph-keeper/pkg/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc/metadata"
)

const (
	login         = "login"
	loginPassword = "login_password"
	bankData      = "bank_data"
	textData      = "text_data"
	binData       = "bin_data"
)

type AccessRepository interface {
	GetAccess(ctx context.Context, login string, id string) (*models.FileInfo, error)
	SaveAccess(ctx context.Context, login string, id string) error
}

type DownloadService struct {
	dataRepository   rep.DataRepository
	accessRepository AccessRepository
}

func New(ctx context.Context, dataRep rep.DataRepository, accessRep AccessRepository) *DownloadService {
	return &DownloadService{dataRepository: dataRep, accessRepository: accessRep}
}

// DownloadFile checks whether user is authorized to download file with certain id,
// if so, downloading begins from storage and file is being returned, if not - returns error.
func (d *DownloadService) DownloadFile(ctx context.Context, id string) ([]byte, string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		logger.Error("metadata is not received from incoming context")

		return nil, "", errors.New("metada not received from md")
	} else if md.Len() == 0 {
		logger.Error("metadata is emty")

		return nil, "", errors.New("md is empty")
	}

	if _, ok = md[login]; !ok {
		logger.Error("login not provided")

		return nil, "", errors.New("login is needed")
	}

	login := md[login][0]
	// remove @ since this charac is not allowed for Minio bucket name
	login = strings.Replace(login, "@", "", -1)

	// get metadata about file with provided id
	fileInfo, err := d.accessRepository.GetAccess(ctx, login, id)
	if err != nil {
		logger.Error("failed to get access for file", zap.Error(err))

		return nil, "", fmt.Errorf("error getting access for specific file from repo: %w", err)
	}

	// check whether user is authorized to get access to this specific file
	if fileInfo.Login != login {
		logger.Info("Authorization error")

		return nil, "", errors.New("authorization error")
	}

	file, metadata, err := d.dataRepository.DownloadFile(ctx, login, id)
	if err != nil {
		return nil, "", fmt.Errorf("error downloading file from repository: %w", err)
	}

	logger.Info("File successfully downloaded from Minio bucket ")

	return file.Bytes(), metadata, nil
}

func (d *DownloadService) DownloadBankData(ctx context.Context, id string) (map[string]string, string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		logger.Error("metada is not received from incoming context")

		return nil, "", errors.New("metada not received from md")
	} else if md.Len() == 0 {
		logger.Error("metada is emty")

		return nil, "", errors.New("md is empty")
	}

	if _, ok = md[login]; !ok {
		logger.Error("login not provided")

		return nil, "", errors.New("login is needed")
	}

	login := md["login"][0]
	// remove @ since this charac is not allowed for Minio bucket name
	login = strings.Replace(login, "@", "", -1)

	// get metadata about file with provided id
	fileInfo, err := d.accessRepository.GetAccess(ctx, login, id)
	if err != nil {
		logger.Error("failed to get access to file", zap.Error(err))

		return nil, "", fmt.Errorf("error getting access for specific file from repo: %w", err)
	}

	// check whether user is authorized to get access to this specific file
	if fileInfo.Login != login {
		logger.Info("Authorization error")

		return nil, "", errors.New("authorization error")
	}

	data, metadata, err := d.dataRepository.DownloadTextData(ctx, login, id, bankData)
	if err != nil {
		return nil, "", fmt.Errorf("error downloading bank details: %s", err)
	}

	res := make(map[string]string, 3)

	err = json.Unmarshal(data, &res)
	if err != nil {
		logger.Error("error unmarshalling JSON:", zap.Error(err))

		return nil, "", fmt.Errorf("error unmarshalling JSON: %w", err)
	}

	return res, metadata, nil
}

func (d *DownloadService) DownloadText(ctx context.Context, id string) (string, string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		logger.Error("metada is not received from incoming context")

		return "", "", errors.New("metada not received from md")
	} else if md.Len() == 0 {
		logger.Error("metada is emty")

		return "", "", errors.New("md is empty")
	}

	if _, ok = md[login]; !ok {
		logger.Error("login not provided")

		return "", "", errors.New("login is needed")
	}

	login := md["login"][0]
	// remove @ since this charac is not allowed for Minio bucket name
	login = strings.Replace(login, "@", "", -1)

	// get metadata about file with provided id
	fileInfo, err := d.accessRepository.GetAccess(ctx, login, id)
	if err != nil {
		logger.Error("failed to get access for file")

		return "", "", fmt.Errorf("error getting access for specific file from repo: %w", err)
	}

	// check whether user is authorized to get access to this specific file
	if fileInfo.Login != login {
		logger.Info("Authorization error")

		return "", "", errors.New("authorization error")
	}

	data, metadata, err := d.dataRepository.DownloadTextData(ctx, login, id, textData)
	if err != nil {
		return "", "", fmt.Errorf("error downloading text: %s", err)
	}

	return string(data), metadata, nil
}

func (d *DownloadService) DownloadLoginPassword(ctx context.Context, id string) (map[string]string, string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		logger.Error("metada is not received from incoming context")

		return nil, "", errors.New("metada not received from md")
	} else if md.Len() == 0 {
		logger.Error("metada is emty")

		return nil, "", errors.New("md is empty")
	}

	if _, ok = md[login]; !ok {
		logger.Error("login not provided")

		return nil, "", errors.New("login is needed")
	}

	login := md["login"][0]
	// remove @ since this charac is not allowed for Minio bucket name
	login = strings.Replace(login, "@", "", -1)

	// get metadata about file with provided id
	fileInfo, err := d.accessRepository.GetAccess(ctx, login, id)
	if err != nil {
		logger.Error("failed to get access for file")

		return nil, "", fmt.Errorf("error getting access for specific file from repo: %w", err)
	}

	// check whether user is authorized to get access to this specific file
	if fileInfo.Login != login {
		logger.Info("Authorization error")

		return nil, "", errors.New("authorization error")
	}

	data, metadata, err := d.dataRepository.DownloadTextData(ctx, login, id, loginPassword)
	if err != nil {
		logger.Error("error downloading login credentials: ", zap.Error(err))

		return nil, "", fmt.Errorf("error downloading login credentials: %s", err)
	}

	res := make(map[string]string, 3)

	err = json.Unmarshal(data, &res)
	if err != nil {
		logger.Error("error unmarshalling JSON:", zap.Error(err))

		return nil, "", fmt.Errorf("error unmarshalling JSON: %w", err)
	}

	return res, metadata, nil
}
