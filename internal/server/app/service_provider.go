package app

import (
	"context"
	"log"

	"github.com/igortoigildin/goph-keeper/internal/client/db"
	"github.com/igortoigildin/goph-keeper/internal/client/db/pg"
	auth "github.com/igortoigildin/goph-keeper/internal/server/api/auth_v1"
	api "github.com/igortoigildin/goph-keeper/internal/server/api/upload_v1"
	"github.com/igortoigildin/goph-keeper/internal/server/closer"
	upload "github.com/igortoigildin/goph-keeper/internal/server/service"
	"github.com/igortoigildin/goph-keeper/pkg/logger"
	"go.uber.org/zap"

	"github.com/igortoigildin/goph-keeper/internal/server/config"
	authService "github.com/igortoigildin/goph-keeper/internal/server/service/auth"
	uploadService "github.com/igortoigildin/goph-keeper/internal/server/service/upload"
	repository "github.com/igortoigildin/goph-keeper/internal/server/storage"
	dataRepository "github.com/igortoigildin/goph-keeper/internal/server/storage/minio"
	accessRepository "github.com/igortoigildin/goph-keeper/internal/server/storage/pg/access"
	userRepository "github.com/igortoigildin/goph-keeper/internal/server/storage/pg/user"
)

type serviceProvider struct {
	grpcConfig config.GRPCConfig
	pgConfig   config.PGConfig

	dbClient db.Client

	uploadService upload.UploadService
	uploadImpl    *api.Implementation

	authService auth.AuthService
	authImpl    *auth.Implementation

	userRepository repository.UserRepository
	dataRepository repository.DataRepository
	accessRepository repository.AccessRepository
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

func (s *serviceProvider) PGConfig() config.PGConfig {
	if s.pgConfig == nil {
		cfg, err := config.NewPGConfig()
		if err != nil {
			logger.Fatal("failed to get pg config:", zap.Error(err))
		}

		s.pgConfig = cfg
	}

	return s.pgConfig
}

func (s *serviceProvider) UploadImpl(ctx context.Context) *api.Implementation {
	if s.uploadImpl == nil {
		s.uploadImpl = api.NewImplementation(s.UploadService(ctx))
	}

	return s.uploadImpl
}

func (s *serviceProvider) UploadService(ctx context.Context) upload.UploadService {
	if s.uploadService == nil {
		s.uploadService = uploadService.New(ctx, s.DataRepository(ctx), s.AccessRepository(ctx))
	}

	return s.uploadService
}

func (s *serviceProvider) AuthService(ctx context.Context) auth.AuthService {
	if s.authService == nil {
		s.authService = authService.New(s.UserRepository(ctx))
	}
	return s.authService
}

func (s *serviceProvider) AuthImpl(ctx context.Context) *auth.Implementation {
	if s.authImpl == nil {
		s.authImpl = auth.NewImplementation(s.AuthService(ctx))
	}

	return s.authImpl
}

func (s *serviceProvider) DBClient(ctx context.Context) db.Client {
	if s.dbClient == nil {
		cl, err := pg.New(ctx, s.PGConfig().DSN())
		if err != nil {
			logger.Fatal("failed to create db client:", zap.Error(err))
		}

		err = cl.DB().Ping(ctx)
		if err != nil {
			logger.Fatal("ping error:", zap.Error(err))
		}

		closer.Add(cl.Close)

		s.dbClient = cl
	}

	return s.dbClient
}

func (s *serviceProvider) UserRepository(ctx context.Context) repository.UserRepository {
	if s.userRepository == nil {
		s.userRepository = userRepository.NewRepository(s.DBClient(ctx))
	}

	return s.userRepository
}

func (s *serviceProvider) DataRepository(ctx context.Context) repository.DataRepository {
	if s.dataRepository == nil {
		s.dataRepository = dataRepository.NewRepository()
	}

	return s.dataRepository
}

func (s *serviceProvider) AccessRepository(ctx context.Context) repository.AccessRepository {
	if s.accessRepository == nil {
		s.accessRepository = accessRepository.NewRepository(s.DBClient(ctx))
	}

	return s.accessRepository
}
