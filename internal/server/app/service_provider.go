package app

import (
	"context"
	"log"

	auth "github.com/igortoigildin/goph-keeper/internal/server/api/auth_v1"
	upload "github.com/igortoigildin/goph-keeper/internal/server/api/upload_v1"

	"github.com/igortoigildin/goph-keeper/internal/server/config"
	authService "github.com/igortoigildin/goph-keeper/internal/server/service/auth"
	uploadService "github.com/igortoigildin/goph-keeper/internal/server/service/upload"
)

type serviceProvider struct {
	grpcConfig config.GRPCConfig

	uploadService upload.UploadService
	uploadImpl    *upload.Implementation

	authService auth.AuthService
	authImpl    *auth.Implementation
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

func (s *serviceProvider) UploadService(ctx context.Context) upload.UploadService {
	if s.uploadService == nil {
		s.uploadService = uploadService.New()
	}

	return s.uploadService
}

func (s *serviceProvider) AuthImpl(ctx context.Context) *auth.Implementation {
	if s.authImpl == nil {
		s.authImpl = auth.NewImplementation(s.AuthService(ctx))
	}

	return s.authImpl
}

func (s *serviceProvider) AuthService(ctx context.Context) auth.AuthService {
	if s.authService == nil {
		s.authService = authService.New()
	}
	return s.authService
}
