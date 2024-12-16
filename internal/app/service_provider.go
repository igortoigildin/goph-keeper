package app

import (
	"context"
	"log"

	"github.com/moby/buildkit/cmd/buildkitd/config"
)

type serviceProvider struct {
	grpcConfig config.GRPCConfig

	uploadService service.uploadService

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

func (s *serviceProvider) NoteImpl(ctx context.Context) *upload.Implementation {
	if s.uploadImpl == nil {
		s.uploadImpl = upload.NewImplementation(s.UploadService(ctx))
	}

	return s.uploadImpl
}