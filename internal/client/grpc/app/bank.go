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

func saveCardInfoCmd(app *App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "card",
		Short: "Save bank card details in storage",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			refreshTokenSecretKey, _ := viper.Get("REFRESH_SECRET").(string)

			if !session.IsSessionValid(refreshTokenSecretKey) {
				logger.Error("Session expired or not found. Please login again")
			}

			logger.Info("Session is valid")
		},
		Run: func(cmd *cobra.Command, args []string) {
			cardNumber, err := cmd.Flags().GetString("card_number")
			if err != nil {
				logger.Fatal("failed to get card_number", zap.Error(err))
			}

			cvc, err := cmd.Flags().GetString("CVC")
			if err != nil {
				logger.Fatal("failed to get CVC", zap.Error(err))
			}

			expDate, err := cmd.Flags().GetString("expiration_date")
			if err != nil {
				logger.Fatal("failed to get expiration_date", zap.Error(err))
			}

			meta, err := cmd.Flags().GetString("info")
			if err != nil {
				logger.Fatal("failed to get metadata", zap.Error(err))
			}

			// Creating new uuid for the bank details to be saved
			id := uuid.New()

			// Creating Upload service
			clientService := serviceUp.New()

			serverAddr, _ := viper.Get("GRPC_PORT").(string)

			err = app.Saver.SaveBankDetails(cardNumber, cvc, expDate, id.String(), meta)
			if err != nil {
				logger.Error("failed to save bank details locally", zap.Error(err))
			}

			if err := clientService.SendBankDetails(fmt.Sprintf(":%s", serverAddr), cardNumber, cvc, expDate, id.String(), meta); err != nil {
				logger.Fatal("failed to save bank details: ", zap.Error(err))
			}

			logger.Info("Your bank details saved successfully. Please keep your uuid and use it to retrive your data back from Goph-keeper.",
				zap.String("uuid:", id.String()))
		},
	}
	cmd.Flags().StringP("card_number", "n", "", "Card number to be saved")
	cmd.Flags().StringP("CVC", "c", "", "CVC to be saved")
	cmd.Flags().StringP("expiration_date", "e", "", "expiration_date to be saved")
	cmd.Flags().StringP("info", "i", "", "Additional metadata, if necessary")

	return cmd
}
