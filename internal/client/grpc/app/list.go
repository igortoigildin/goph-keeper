package app

import (
	"github.com/igortoigildin/goph-keeper/pkg/logger"
	"github.com/igortoigildin/goph-keeper/pkg/session"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// list all
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all secrets command",
}

func listAllSavedSecrets(app *App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "all",
		Short: "List all data currently saved in gopher-keeper",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			refreshTokenSecretKey, _ := viper.Get("REFRESH_SECRET").(string)

			if !session.IsSessionValid(refreshTokenSecretKey) {
				logger.Error("Session expired or not found. Please login again")
			}

			logger.Info("Session is valid")
		},
		Run: func(cmd *cobra.Command, args []string) {
			cards, err := app.Receiver.GetAllBankDetails()
			if err != nil {
				logger.Error("failed to get card secrets list", zap.Error(err))

				return
			}
			if len(cards) != 0 {
				logger.Info("At the moment, you have the following card secrets:")
			} else {
				logger.Info("You do nat have any card secrets saved. Please add your secrets to goph-keeper")
			}

			for _, card := range cards {

				logger.Info("Card details:", zap.Any("Secret ID", card.ID),
					zap.Any("metadata:", card.Info),
				)
			}

			creds, err := app.Receiver.GetAllCredentials()
			if err != nil {
				logger.Error("failed to get credential secrets list", zap.Error(err))

				return
			}
			if len(creds) != 0 {
				logger.Info("At the moment, you have the following credential secrets:")
			} else {
				logger.Info("You do nat have any credential secrets saved. Please add your secrets to goph-keeper")
			}

			for _, cred := range creds {

				logger.Info("Lodin && Password pair:", zap.Any("Secret ID", cred.ID),
					zap.Any("metadata:", cred.Service),
				)
			}

			files, err := app.Receiver.ListAllFiles()
			if err != nil {
				logger.Error("failed to get files list", zap.Error(err))

				return
			}

			if len(files) != 0 {
				logger.Info("At the moment, you have the following files saved:")
			} else {
				logger.Info("You do nat have any files saved. Please add your files to goph-keeper")
			}

			for _, file := range files {

				logger.Info("File details:", zap.Any("File ID", file.ID),
					zap.Any("metadata:", file.Info),
				)
			}

			texts, err := app.Receiver.GetAllTexts()
			if err != nil {
				logger.Error("failed to get files list", zap.Error(err))

				return
			}
			if len(texts) != 0 {
				logger.Info("At the moment, you have the following text data saved:")
			} else {
				logger.Info("You do nat have any any text data saved. Please add your text data to goph-keeper")
			}
			for _, text := range texts {

				logger.Info("Text data details:", zap.Any("Text data ID", text.ID),
					zap.Any("metadata:", text.Info),
				)
			}

		},
	}

	return cmd
}
