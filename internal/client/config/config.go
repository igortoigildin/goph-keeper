package config

import (
	"fmt"
	"strings"

	"github.com/igortoigildin/goph-keeper/buildinfo"
	"github.com/igortoigildin/goph-keeper/pkg/logger"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

func LoadConfig() error {
	fmt.Printf("Build Version: %s\n", buildinfo.BuildVersion)
	fmt.Printf("Build Date: %s\n", buildinfo.BuildDate)
	fmt.Printf("Build Commit: %s\n", buildinfo.BuildCommit)

	err := godotenv.Load("././.env")
	if err != nil {
		logger.Error("Failed to load .env file")
	}

	if err := initConfig(); err != nil {
		return fmt.Errorf("error initialization configs: %s", err.Error())
	}

	return nil
}

func initConfig() error {
	viper.SetConfigFile("././.env")

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("errro reading config file: %w", err)
	}

	viper.AutomaticEnv()

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	return nil
}
