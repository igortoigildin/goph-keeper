package app

import (
	"fmt"

	"github.com/igortoigildin/goph-keeper/pkg/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// sync all data with server
var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync all data with server",
}

func syncAllData(app *App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "all",
		Short: "Sync all data currently saved in gopher-keeper",
		Run: func(cmd *cobra.Command, args []string) {

			serverAddr, _ := viper.Get("GRPC_PORT").(string)

			// cards, err := app.Receiver.GetAllBankDetails()
			// if err != nil {
			// 	logger.Error("failed to get card secrets list", zap.Error(err))

			// 	return
			// }

			// for _, card := range cards {

			// 	cardData, err := app.Downloader.GetFile(card.ID)
			// 	if err != nil {
			// 		logger.Error("failed to get file", zap.Error(err))

			// 		return
			// 	}

			// 	if cardData.Etag != card.Etag {
			// 		err = app.Saver.SaveBankDetails()
			// 	}
			// }

			err := app.Syncer.SyncAllData(fmt.Sprintf(":%s", serverAddr))
			if err != nil {
				logger.Error("failed to sync all data", zap.Error(err))

				return
			}

		},
	}
	return cmd
}
