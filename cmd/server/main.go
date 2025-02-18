package main

import (
	"context"
	"log"

	"github.com/igortoigildin/goph-keeper/internal/server/app"
	"github.com/igortoigildin/goph-keeper/internal/server/config"
	"github.com/igortoigildin/goph-keeper/pkg/logger"
)

// var configPath string

// func init() {
// 	flag.StringVar(&configPath, "config-path", ".env", "path to config file")
// }

func main() {
	cfg := config.MustLoad()
	logger.Initialize(cfg.FlagLogLevel)

	app, err := app.NewApp(context.Background())
	if err != nil {
		log.Fatalf("failed to init app: %s", err.Error())
	}

	logger.Info("app initialized successfully")

	err = app.Run()
	if err != nil {
		log.Fatalf("failed to run app: %s", err.Error())
	}
}
