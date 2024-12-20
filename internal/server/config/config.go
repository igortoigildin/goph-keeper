package config

import (
	"flag"
	"os"

	"github.com/joho/godotenv"
)

func LoadFromFile(path string) error {
	err := godotenv.Load(path)
	if err != nil {
		return err
	}

	return nil
}

type ConfigServer struct {
	FlagLogLevel string
}

func LoadConfig() *ConfigServer {
	cfg := new(ConfigServer)

	flag.StringVar(&cfg.FlagLogLevel, "l", "info", "log level")
	flag.Parse()

	if envLogLevel := os.Getenv("LOG_LEVEL"); envLogLevel != "" {
		cfg.FlagLogLevel = envLogLevel
	}

	return cfg
}
