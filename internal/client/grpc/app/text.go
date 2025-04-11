package app

import (
	"fmt"

	"github.com/google/uuid"
	serviceUp "github.com/igortoigildin/goph-keeper/internal/client/grpc/service/upload"

	"github.com/igortoigildin/goph-keeper/pkg/logger"
	"github.com/igortoigildin/goph-keeper/pkg/session"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func saveTextCmd(app *App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "text",
		Short: "Save arbitrary text data in storage",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			refreshTokenSecretKey, _ := viper.Get("REFRESH_SECRET").(string)

			if !session.IsSessionValid(refreshTokenSecretKey) {
				logger.Error("Session expired or not found. Please login again")
			}

			logger.Info("Session is valid")
		},
		Run: func(cmd *cobra.Command, args []string) {
			textData, err := cmd.Flags().GetString("text")
			if err != nil {
				logger.Fatal("failed to get text to be saved:", zap.Error(err))
			}

			info, err := cmd.Flags().GetString("info")
			if err != nil {
				logger.Fatal("failed to get additional information", zap.Error(err))
			}

			// Creating new uuid for text to be saved
			id := uuid.New()

			// Initializing Upload service
			clientService := serviceUp.New()

			serverAddr, _ := viper.Get("GRPC_PORT").(string)

			err = app.Saver.SaveText(id.String(), info, textData)
			if err != nil {
				logger.Error("failed to save text locally", zap.Error(err))
			}

			if err := clientService.SendText(fmt.Sprintf(":%s", serverAddr), textData, id.String(), info); err != nil {
				logger.Fatal("failed to save text", zap.Error(err))
			}

			logger.Info("Your text saved successfully", zap.String("uuid:", id.String()))
		},
	}

	cmd.Flags().StringP("text", "t", "", "Text which need to be saved")
	cmd.Flags().StringP("info", "i", "", "Additional metadata, if necessary")

	return cmd
}
