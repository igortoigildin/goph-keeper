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

func savePasswordCmd(app *App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "password",
		Short: "Save login && password in storage",
		Run: func(cmd *cobra.Command, args []string) {
			loginStr, err := cmd.Flags().GetString("login")
			if err != nil {
				logger.Fatal("failed to get login:", zap.Error(err))
			}

			encryptedLogin, err := encryption.Encrypt(loginStr, []byte(viper.Get("ENCRYPTION_KEY").(string)))
			if err != nil {
				logger.Error("failed to encrypt login", zap.Error(err))
			}

			passStr, err := cmd.Flags().GetString("password")
			if err != nil {
				logger.Fatal("failed to get password:", zap.Error(err))
			}

			encryptedPassword, err := encryption.Encrypt(passStr, []byte(viper.Get("ENCRYPTION_KEY").(string)))
			if err != nil {
				logger.Error("failed to encrypt password", zap.Error(err))
			}

			meta, err := cmd.Flags().GetString("service")
			if err != nil {
				logger.Fatal("failed to get metadata", zap.Error(err))
			}

			// Initializing Upload service.
			clientService := serviceUp.New()

			// Creating new uuid for credentials to be saved.
			id := uuid.New()

			serverAddr, _ := viper.Get("GRPC_PORT").(string)

			// Sending credentials with created uuid to server.
			etag, err := clientService.SendPassword(fmt.Sprintf(":%s", serverAddr), encryptedLogin, encryptedPassword, id.String(), meta)
			if err != nil {
				logger.Error("failed to send credentials to server:", zap.Error(err))
			}

			// Saving credentials to local client storage
			err = app.ClientSaver.SaveCredentials(id.String(), meta, encryptedLogin, encryptedPassword, etag)
			if err != nil {
				logger.Error("failed to save credentials locally", zap.Error(err))
			}

			logger.Info("Credentials saved successfully. Please save your uuid and use it to retrive your data back from Goph-keeper.",
				zap.String("uuid:", id.String()))
		},
	}

	cmd.Flags().StringP("login", "l", "", "Login to be saved")
	cmd.Flags().StringP("password", "p", "", "Password to be saved")
	cmd.Flags().StringP("service", "d", "", `Name of the site, app, or other platform
	for which the login and password were created.`)

	return cmd
}

func downloadPassCmd(app *App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "password",
		Short: "Download login && password from storage",
		Run: func(cmd *cobra.Command, args []string) {
			idStr, err := cmd.Flags().GetString("id")
			if err != nil {
				logger.Fatal("failed to get credentials id", zap.Error(err))
			}

			clientService := serviceDown.New()
			serverAddr, _ := viper.Get("GRPC_PORT").(string)

			_, err = clientService.DownloadPassword(fmt.Sprintf(":%s", serverAddr), idStr)
			if err != nil {
				// if remote server is not available, try to reach local storage
				res, err := app.ClientReceiver.GetCredential(idStr)
				if err != nil {
					logger.Error("failed to download date from local storage", zap.Error(err))
				}

				decryptedLogin, err := encryption.Decrypt(res.Username, []byte(viper.Get("ENCRYPTION_KEY").(string)))
				if err != nil {
					logger.Error("failed to decrypt login", zap.Error(err))
				}

				decryptedPassword, err := encryption.Decrypt(res.Password, []byte(viper.Get("ENCRYPTION_KEY").(string)))
				if err != nil {
					logger.Error("failed to decrypt password", zap.Error(err))
				}

				logger.Info("Your data: ", zap.Any("login", decryptedLogin), zap.Any("password", decryptedPassword))

			}
		},
	}

	cmd.Flags().StringP("id", "i", "", "A Universally Unique Identifier of the saved password")

	return cmd
}
