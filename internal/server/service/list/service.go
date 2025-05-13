package list

import (
	"context"
	"errors"
	"fmt"
	"strings"

	model "github.com/igortoigildin/goph-keeper/internal/server/models"
	rep "github.com/igortoigildin/goph-keeper/internal/server/storage"
	"github.com/igortoigildin/goph-keeper/pkg/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc/metadata"
)

const (
	login = "login"
)

type AccessRepository interface {
	GetAccess(ctx context.Context, login string, id string) (*model.FileInfo, error)
	SaveAccess(ctx context.Context, login string, id string) error
}

type ListService struct {
	dataRepository   rep.DataRepository
	accessRepository AccessRepository
}

func New(ctx context.Context, dataRep rep.DataRepository, accessRep AccessRepository) *ListService {
	return &ListService{dataRepository: dataRep, accessRepository: accessRep}
}

func (l *ListService) List(ctx context.Context) ([]model.ObjectInfo, error) {

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		logger.Error("metadata is not received from incoming context")

		return nil, errors.New("metada not received from md")
	} else if md.Len() == 0 {
		logger.Error("metadata is emty")

		return nil, errors.New("md is empty")
	}

	if _, ok = md[login]; !ok {
		logger.Error("login not provided")

		return nil, errors.New("login is needed")
	}

	login := md[login][0]
	// remove @ since this charac is not allowed for Minio bucket name
	login = strings.Replace(login, "@", "", -1)

	objs, err := l.dataRepository.ListObjects(ctx, login)
	if err != nil {
		logger.Error("failed to list objects", zap.Error(err))

		return nil, fmt.Errorf("error listing objects: %w", err)
	}

	return objs, nil
}
