package app

import (
	"context"
	"log"

	upload "github.com/igortoigildin/goph-keeper/internal/server/api/note"
	"github.com/igortoigildin/goph-keeper/internal/server/config"
	"github.com/igortoigildin/goph-keeper/internal/server/service"
	uploadService "github.com/igortoigildin/goph-keeper/internal/server/service/upload"
)

type serviceProvider struct {
	grpcConfig config.GRPCConfig

	uploadService service.UploadService

	uploadImpl	*upload.Implementation
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

func (s *serviceProvider) GRPCConfig() config.GRPCConfig {
	if s.grpcConfig == nil {
		cfg, err := config.NewGRPCConfig()
		if err != nil {
			log.Fatalf("failed to get grpc config: %s", err.Error())
		}

		s.grpcConfig = cfg
	}

	return s.grpcConfig
}

func (s *serviceProvider) UploadImpl(ctx context.Context) *upload.Implementation {
	if s.uploadImpl == nil {
		s.uploadImpl = upload.NewImplementation(s.UploadService(ctx))
	}

	return s.uploadImpl
}

func (s *serviceProvider) UploadService(ctx context.Context) service.UploadService {
	if s.uploadService == nil {
		s.uploadService = uploadService.New()
	}

	return s.uploadService
}