package config

import (
	"errors"
	"net"
	"os"
)

const (
	grpcHostEnvName    = "GRPC_HOST"
	grpcPortEnvName    = "GRPC_PORT"
	loggerLevelEnvName = "LOGGER_LEVEL"
)

type GRPCConfig interface {
	Address() string
}

type GrpcConfig struct {
	Host string
	Port string `yaml:"port"`
}

func NewGRPCConfig() (GRPCConfig, error) {
	host := os.Getenv(grpcHostEnvName)
	if len(host) == 0 {
		return nil, errors.New("grpc host not found")
	}

	port := os.Getenv(grpcPortEnvName)
	if len(port) == 0 {
		return nil, errors.New("grpc port not found")
	}

	return &GrpcConfig{
		Host: host,
		Port: port,
	}, nil
}

func (cfg *GrpcConfig) Address() string {
	return net.JoinHostPort(cfg.Host, cfg.Port)
}
