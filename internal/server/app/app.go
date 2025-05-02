package app

import (
	"context"
	"fmt"
	"net"

	"github.com/igortoigildin/goph-keeper/internal/server/closer"
	config "github.com/igortoigildin/goph-keeper/internal/server/config"
	authpb "github.com/igortoigildin/goph-keeper/pkg/auth_v1"
	downloadpb "github.com/igortoigildin/goph-keeper/pkg/download_v1"
	"github.com/igortoigildin/goph-keeper/pkg/logger"
	uploadpb "github.com/igortoigildin/goph-keeper/pkg/upload_v1"
	"go.uber.org/zap"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
)

const (
	cfgFileName = ".env"
)

type App struct {
	serviceProvider *serviceProvider
	grpcServer      *grpc.Server
}

func NewApp(ctx context.Context) (*App, error) {
	a := &App{}

	err := a.initDeps(ctx)
	if err != nil {
		return nil, fmt.Errorf("error inintializing dependecies: %w", err)
	}

	return a, nil
}

func (a *App) Run() error {
	defer func() {
		closer.CloseAll()
		closer.Wait()
	}()

	return a.runGRPCServer()
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initConfig,
		a.initServiceProvider,
		a.initGRPCServer,
	}

	for _, f := range inits {
		err := f(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *App) initConfig(_ context.Context) error {
	err := config.LoadFromFile(cfgFileName)
	if err != nil {
		return fmt.Errorf("error loading config from local file: %w", err)
	}

	return nil
}

func (a *App) initGRPCServer(ctx context.Context) error {

	creds, err := credentials.NewServerTLSFromFile("certs/server.crt", "certs/server.key")
	if err != nil {
		logger.Error("failed to load TLS certificates: %w", zap.Error(err))

		return fmt.Errorf("failed to load TLS certificates: %w", err)
	}

	a.grpcServer = grpc.NewServer(grpc.Creds(creds))

	reflection.Register(a.grpcServer)

	uploadpb.RegisterUploadV1Server(a.grpcServer, a.serviceProvider.UploadImpl(ctx))
	authpb.RegisterAuthV1Server(a.grpcServer, a.serviceProvider.AuthImpl(ctx))
	downloadpb.RegisterDownloadV1Server(a.grpcServer, a.serviceProvider.DownloadImpl(ctx))

	return nil
}

func (a *App) runGRPCServer() error {
	logger.Info("GRPC server with TLS is running on:", zap.Any("address:", a.serviceProvider.GRPCConfig().Address()))

	list, err := net.Listen("tcp", a.serviceProvider.GRPCConfig().Address())
	if err != nil {
		return fmt.Errorf("error announcing on the local network address: %w", err)
	}

	err = a.grpcServer.Serve(list)
	if err != nil {
		return fmt.Errorf("error accepting incoming connections on the listener: %w", err)
	}

	return nil
}

func (a *App) initServiceProvider(_ context.Context) error {
	a.serviceProvider = newServiceProvider()

	return nil
}
