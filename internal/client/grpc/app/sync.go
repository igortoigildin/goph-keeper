package app

import (
	"fmt"
	"strings"

	serviceDown "github.com/igortoigildin/goph-keeper/internal/client/grpc/service/download"
	"github.com/igortoigildin/goph-keeper/pkg/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

const (
	loginPassword = "login_password_"
	bankData      = "bank_data_"
	textData      = "text_data_"
	binData       = "bin_data_"
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

			err := app.RunSync(fmt.Sprintf(":%s", serverAddr))
			if err != nil {
				logger.Error("failed to sync all data", zap.Error(err))
				return
			}

		},
	}
	return cmd
}

// RunSync obtains list of all objects from server, uses etag to check if object is up to date.
// if not, it downloads data from remote server via gRPC and updates local storage accordingly.
func (app *App) RunSync(addr string) error {
	objects, err := app.Syncer.ListAllData(addr)
	if err != nil {
		return fmt.Errorf("error getting object list: %w", err)
	}
	serverAddr, _ := viper.Get("GRPC_PORT").(string)
	clientService := serviceDown.New()

	for _, object := range objects {
		// text data
		if strings.HasPrefix(object.Key, textData) {
			objName := strings.TrimPrefix(object.Key, textData)
			res, err := app.ClientReceiver.GetText(objName)
			if err != nil {
				return fmt.Errorf("error getting text data: %w", err)
			}
			if res.Etag != object.Etag {
				// Get updated text from server
				// Requesting text with provided uuid.
				serverObj, err := clientService.DownloadText(fmt.Sprintf(":%s", serverAddr), objName)
				if err != nil {
					return fmt.Errorf("error downloading text data: %w", err)
				}

				// Update text in local storage
				err = app.ClientSaver.UpdateText(objName, serverObj.Text, object.Etag)
				if err != nil {
					return fmt.Errorf("error updating text data: %w", err)
				}
			}
		} else if strings.HasPrefix(object.Key, bankData) {
			// bank data
			objName := strings.TrimPrefix(object.Key, bankData)
			res, err := app.ClientReceiver.GetBankDetails(objName)
			if err != nil {
				return fmt.Errorf("error getting bank details: %w", err)
			}

			if res.Etag != object.Etag {
				// Get updated bank details from server
				serverObj, err := clientService.DownloadBankDetails(fmt.Sprintf(":%s", serverAddr), objName)
				if err != nil {
					return fmt.Errorf("error downloading bank details: %w", err)
				}

				// Update bank details in local storage
				err = app.ClientSaver.UpdateBankDetails(objName, serverObj.CardNumber, serverObj.Cvc, serverObj.ExpDate, serverObj.Info, object.Etag)
				if err != nil {
					return fmt.Errorf("error updating bank details: %w", err)
				}
			}
		} else if strings.HasPrefix(object.Key, binData) {
			// bin data
			objName := strings.TrimPrefix(object.Key, binData)
			res, err := app.ClientReceiver.GetFile(objName)
			if err != nil {
				return fmt.Errorf("error getting file: %w", err)
			}

			if res.Etag != object.Etag {
				// Get updated file from server
				serverObj, err := clientService.DownloadFile(fmt.Sprintf(":%s", serverAddr), objName, res.Filename)
				if err != nil {
					return fmt.Errorf("error downloading file: %w", err)
				}

				// Update file in local storage
				err = app.ClientSaver.UpdateFile(objName, object.Etag, serverObj.Data)
				if err != nil {
					return fmt.Errorf("error updating file: %w", err)
				}
			}
		} else if strings.HasPrefix(object.Key, loginPassword) {
			// credentials data
			objName := strings.TrimPrefix(object.Key, loginPassword)
			res, err := app.ClientReceiver.GetCredential(objName)
			if err != nil {
				return fmt.Errorf("error getting login password: %w", err)
			}

			if res.Etag != object.Etag {
				// Get updated login password from server
				serverObj, err := clientService.DownloadPassword(fmt.Sprintf(":%s", serverAddr), objName)
				if err != nil {
					return fmt.Errorf("error downloading login password: %w", err)
				}

				// Update login password in local storage
				err = app.ClientSaver.UpdateCredentials(objName, serverObj.Service, serverObj.Username, serverObj.Password, object.Etag)
				if err != nil {
					return fmt.Errorf("error updating login password: %w", err)
				}
			}
		}
	}

	logger.Info("Synchronization with server has been completed successfully. All files updated.")

	return nil
}
