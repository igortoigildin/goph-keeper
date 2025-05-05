package config

import (
	"errors"
	"os"
)

const (
	dsnEnvName            = "PG_DSN"
	migrationsPathEnvName = "PG_MIGRATIONS_PATH"
)

type PGConfig interface {
	DSN() string
	MigrationsPath() string
}

type pgConfig struct {
	dsn            string
	migrationsPath string
}

func NewPGConfig(cfg *Config) (PGConfig, error) {
	dsn := cfg.PG.DSN
	if len(dsn) == 0 {
		dsn = os.Getenv(dsnEnvName)
		if len(dsn) == 0 {
			return nil, errors.New("pg dsn not found")
		}
	}

	migrationsPath := cfg.PG.MigrationsPath
	if len(migrationsPath) == 0 {
		migrationsPath = os.Getenv(migrationsPathEnvName)
		if len(migrationsPath) == 0 {
			return nil, errors.New("pg migrations path not found")
		}
	}

	return &pgConfig{
		dsn:            dsn,
		migrationsPath: migrationsPath,
	}, nil
}

func (cfg *pgConfig) DSN() string {
	return cfg.dsn
}

func (cfg *pgConfig) MigrationsPath() string {
	return cfg.migrationsPath
}
