package app

import (
	"log"
	"os"

	service "github.com/igortoigildin/goph-keeper/internal/client/service"
	"github.com/igortoigildin/goph-keeper/pkg/logger"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var (
	serverAddr  string
	filePath    string
	batchSize   int
	loggerLevel string
	rootCmd     = &cobra.Command{
		Use:   "transfer_client",
		Short: "Sending files via gRPC",
		Run: func(cmd *cobra.Command, args []string) {
			clientService := service.New(serverAddr, filePath, batchSize)

			if err := clientService.SendFile(); err != nil {
				log.Fatal(err)
			}
		},
	}
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		logger.Error("error while executing root cmd", zap.Error(err))
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVarP(&serverAddr, "addr", "a", "", "server address")
	rootCmd.Flags().StringVarP(&filePath, "file", "f", "", "file path")
	rootCmd.Flags().IntVarP(&batchSize, "batch", "b", 1024*1024, "batch size for sending")
	rootCmd.Flags().StringVarP(&loggerLevel, "log", "l", "info", "logger level")

	logger.Initialize(loggerLevel)

	if err := rootCmd.MarkFlagRequired("file"); err != nil {
		log.Fatal(err)
	}

	if err := rootCmd.MarkFlagRequired("addr"); err != nil {
		log.Fatal(err)
	}
}
