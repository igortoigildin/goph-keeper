package config

import (
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

func LoadFromFile(path string) error {
	err := godotenv.Load(path)
	if err != nil {
		return err
	}

	return nil
}

type Config struct {
	FlagLogLevel   string
	Env            string        `yaml:"env" env:"ENV" env-default:"local"`
	StoragePath    string        `yaml:"storage_path" env:"STORAGE_PATH" env-default:"./storage/goph.db"`
	GRPC           GrpcConfig    `yaml:"grpc"`
	Timeout        time.Duration `yaml:"timeout" env:"TIMEOUT" env-default:"15s"`
	MigrationsPath string
	TokenTTL       time.Duration `yaml:"token_ttl" env-default:"1h"`
	Key            string        `yaml:"key" env:"ENCRYPTION_KEY" env-default:"test_encryption_key"`
	PG             struct {
		DSN            string `yaml:"dsn" env:"PG_DSN"`
		MigrationsPath string `yaml:"migrations_path" env:"PG_MIGRATIONS_PATH"`
	} `yaml:"pg"`
}

func MustLoad() *Config {
	return MustLoadPath("config/local_tests.yaml")
}

func MustLoadPath(path string) *Config {
	// check if file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("config file does not exist: " + path)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		panic("failed to load config: " + err.Error())
	}

	if envLogLevel := os.Getenv("LOG_LEVEL"); envLogLevel != "" {
		cfg.FlagLogLevel = envLogLevel
	}

	return &cfg
}
