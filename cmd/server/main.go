package main

import (
	"context"
	"flag"
	"log"

	"github.com/igortoigildin/goph-keeper/internal/server/app"
	"github.com/igortoigildin/goph-keeper/internal/server/config"
	"github.com/igortoigildin/goph-keeper/pkg/logger"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", ".env", "path to config file")
}

func main() {
	ctx := context.Background()
	cfg := config.LoadConfig()
	logger.Initialize(cfg.FlagLogLevel)

	a, err := app.NewApp(ctx)
	if err != nil {
		log.Fatalf("failed to init app: %s", err.Error())
	}

	logger.Info("app initialized successfully")

	err = a.Run()
	if err != nil {
		log.Fatalf("failed to run app: %s", err.Error())
	}
}
