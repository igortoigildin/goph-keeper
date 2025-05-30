package app

import (
	"context"
	"log"

	"github.com/igortoigildin/goph-keeper/internal/client/db"
	"github.com/igortoigildin/goph-keeper/internal/client/db/pg"
	auth "github.com/igortoigildin/goph-keeper/internal/server/api/auth_v1"
	download "github.com/igortoigildin/goph-keeper/internal/server/api/download_v1"
	api "github.com/igortoigildin/goph-keeper/internal/server/api/upload_v1"
	"github.com/igortoigildin/goph-keeper/internal/server/closer"
	service "github.com/igortoigildin/goph-keeper/internal/server/service"
	"github.com/igortoigildin/goph-keeper/pkg/logger"
	"go.uber.org/zap"

	downloadApi "github.com/igortoigildin/goph-keeper/internal/server/api/download_v1"
	listApi "github.com/igortoigildin/goph-keeper/internal/server/api/list_v1"
	"github.com/igortoigildin/goph-keeper/internal/server/config"
	authService "github.com/igortoigildin/goph-keeper/internal/server/service/auth"
	downloadService "github.com/igortoigildin/goph-keeper/internal/server/service/download"
	listService "github.com/igortoigildin/goph-keeper/internal/server/service/list"
	uploadService "github.com/igortoigildin/goph-keeper/internal/server/service/upload"
	repository "github.com/igortoigildin/goph-keeper/internal/server/storage"
	dataRepository "github.com/igortoigildin/goph-keeper/internal/server/storage/minio"
	accessRepository "github.com/igortoigildin/goph-keeper/internal/server/storage/pg/access"
	userRepository "github.com/igortoigildin/goph-keeper/internal/server/storage/pg/user"
)

type serviceProvider struct {
	grpcConfig config.GRPCConfig
	pgConfig   config.PGConfig
	mainConfig *config.Config

	dbClient db.Client

	uploadService service.UploadService
	uploadImpl    *api.Implementation

	authService service.AuthService
	authImpl    *auth.Implementation

	downloadService service.DownloadService
	downloadImpl    *downloadApi.Implementation

	listService service.ListService
	listImpl    *listApi.Implementation

	userRepository   repository.UserRepository
	dataRepository   repository.DataRepository
	accessRepository downloadService.AccessRepository
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{
		mainConfig: config.MustLoad(),
	}
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
		cfg, err := config.NewPGConfig(s.mainConfig)
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

func (s *serviceProvider) UploadService(ctx context.Context) service.UploadService {
	if s.uploadService == nil {
		s.uploadService = uploadService.New(ctx, s.DataRepository(ctx), s.AccessRepository(ctx))
	}

	return s.uploadService
}

func (s *serviceProvider) AuthService(ctx context.Context) service.AuthService {
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

func (s *serviceProvider) DownloadService(ctx context.Context) service.DownloadService {
	if s.downloadService == nil {
		s.downloadService = downloadService.New(ctx, s.DataRepository(ctx), s.AccessRepository(ctx))
	}

	return s.downloadService
}

func (s *serviceProvider) DownloadImpl(ctx context.Context) *download.Implementation {
	if s.downloadImpl == nil {
		s.downloadImpl = download.NewImplementation(s.DownloadService(ctx))
	}

	return s.downloadImpl
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

func (s *serviceProvider) AccessRepository(ctx context.Context) downloadService.AccessRepository {
	if s.accessRepository == nil {
		s.accessRepository = accessRepository.NewRepository(s.DBClient(ctx))
	}

	return s.accessRepository
}

func (s *serviceProvider) ListImpl(ctx context.Context) *listApi.Implementation {
	if s.listImpl == nil {
		s.listImpl = listApi.NewImplementation(s.ListService(ctx))
	}

	return s.listImpl
}

func (s *serviceProvider) ListService(ctx context.Context) service.ListService {
	if s.listService == nil {
		s.listService = listService.New(ctx, s.DataRepository(ctx), s.AccessRepository(ctx))
	}

	return s.listService
}
