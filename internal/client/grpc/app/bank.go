package app

import (
	"fmt"

	"github.com/google/uuid"
	serviceDown "github.com/igortoigildin/goph-keeper/internal/client/grpc/service/download"
	serviceUp "github.com/igortoigildin/goph-keeper/internal/client/grpc/service/upload"
	"github.com/igortoigildin/goph-keeper/pkg/encryption"
	"github.com/igortoigildin/goph-keeper/pkg/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func saveCardInfoCmd(app *App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "card",
		Short: "Save bank card details in storage",
		// PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// 	// refreshTokenSecretKey, _ := viper.Get("REFRESH_SECRET").(string)

		// 	// if !session.IsSessionValid(refreshTokenSecretKey) {
		// 	// 	logger.Error("Session expired or not found. Please login again")
		// 	// }

		// 	// logger.Info("Session is valid")
		// },
		Run: func(cmd *cobra.Command, args []string) {
			cardNumber, err := cmd.Flags().GetString("card_number")
			if err != nil {
				logger.Fatal("failed to get card_number", zap.Error(err))
			}

			// Encrypting card number
			encryptedCardNumber, err := encryption.Encrypt(cardNumber, []byte(viper.Get("ENCRYPTION_KEY").(string)))
			if err != nil {
				logger.Error("failed to encrypt card number", zap.Error(err))
			}

			cvc, err := cmd.Flags().GetString("CVC")
			if err != nil {
				logger.Fatal("failed to get CVC", zap.Error(err))
			}

			// Encrypting cvc
			encryptedCVC, err := encryption.Encrypt(cvc, []byte(viper.Get("ENCRYPTION_KEY").(string)))
			if err != nil {
				logger.Error("failed to encrypt cvc", zap.Error(err))
			}

			expDate, err := cmd.Flags().GetString("expiration_date")
			if err != nil {
				logger.Fatal("failed to get expiration_date", zap.Error(err))
			}

			// Encrypting expiration date
			encryptedExpDate, err := encryption.Encrypt(expDate, []byte(viper.Get("ENCRYPTION_KEY").(string)))
			if err != nil {
				logger.Error("failed to encrypt expiration date", zap.Error(err))
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

			// Save data to local storate
			err = app.Saver.SaveBankDetails(encryptedCardNumber, encryptedCVC, encryptedExpDate, id.String(), meta)
			if err != nil {
				logger.Error("failed to save bank details locally", zap.Error(err))
			}

			// Upload data to remote server
			if err := clientService.SendBankDetails(fmt.Sprintf(":%s", serverAddr), encryptedCardNumber, encryptedCVC, encryptedExpDate, id.String(), meta); err != nil {
				logger.Error("failed to save bank details: ", zap.Error(err))
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

func downloadCardInfoCmd(app *App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "card",
		Short: "Download card details from storage",
		// PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// 	refreshTokenSecretKey, _ := viper.Get("REFRESH_SECRET").(string)

		// 	if !session.IsSessionValid(refreshTokenSecretKey) {
		// 		logger.Fatal("Session expired or not found. Please login again")
		// 	}

		// 	logger.Info("Session is valid")
		// },
		Run: func(cmd *cobra.Command, args []string) {
			idStr, err := cmd.Flags().GetString("id")
			if err != nil {
				logger.Fatal("failed to get bank details uuid:", zap.Error(err))
			}

			// Initializing Download service
			clientService := serviceDown.New()

			serverAddr, _ := viper.Get("GRPC_PORT").(string)

			// obtain data from remote server
			if err := clientService.DownloadBankDetails(fmt.Sprintf(":%s", serverAddr), idStr); err != nil {
				logger.Error("failed to obtain card details from goph-keeper: ", zap.Error(err))

				// if remote server not responding, try reach local storage
				logger.Info("trying to obtain data locally")

				res, err := app.Downloader.GetBankDetails(idStr)
				if err != nil {
					logger.Error("failed to download bank details locally: ", zap.Error(err))
				}

				// Decrypting card number
				decryptedCardNumber, err := encryption.Decrypt(res.CardNumber, []byte(viper.Get("ENCRYPTION_KEY").(string)))
				if err != nil {
					logger.Error("failed to decrypt card number", zap.Error(err))
				}

				// Decrypting cvc
				decryptedCVC, err := encryption.Decrypt(res.Cvc, []byte(viper.Get("ENCRYPTION_KEY").(string)))
				if err != nil {
					logger.Error("failed to decrypt cvc", zap.Error(err))
				}

				// Decrypting expiration date
				decryptedExpDate, err := encryption.Decrypt(res.ExpDate, []byte(viper.Get("ENCRYPTION_KEY").(string)))
				if err != nil {
					logger.Error("failed to decrypt expiration date", zap.Error(err))
				}

				logger.Info("data from local storage:", zap.Any("card_number", decryptedCardNumber),
					zap.Any("CVC", decryptedCVC),
					zap.Any("expiration_date", decryptedExpDate),
				)

			}

		},
	}
	cmd.Flags().StringP("id", "i", "", "A Universally Unique Identifier of the saved card details")

	return cmd
}
