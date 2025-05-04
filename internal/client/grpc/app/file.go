package app

import (
	"fmt"
	"log"

	"github.com/google/uuid"
	serviceDown "github.com/igortoigildin/goph-keeper/internal/client/grpc/service/download"
	serviceUp "github.com/igortoigildin/goph-keeper/internal/client/grpc/service/upload"
	fl "github.com/igortoigildin/goph-keeper/pkg/file"
	"github.com/igortoigildin/goph-keeper/pkg/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func saveBinCmd(app *App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bin",
		Short: "Save binary data in storage",
		Run: func(cmd *cobra.Command, args []string) {
			pathStr, err := cmd.Flags().GetString("file_path")
			if err != nil {
				log.Fatalf("failed to get path: %s\n", err.Error())
			}

			info, err := cmd.Flags().GetString("info")
			if err != nil {
				logger.Fatal("failed to get metadata", zap.Error(err))
			}

			// Creating new uuid for the file to be saved
			id := uuid.New()

			// Creating Upload service
			clientService := serviceUp.New()

			serverAddr, _ := viper.Get("GRPC_PORT").(string)

			etag, err := clientService.SendFile(fmt.Sprintf(":%s", serverAddr), pathStr, batchSize, id.String(), info)
			if err != nil {
				logger.Fatal("failed to save binary file: ", zap.Error(err))
			}

			// save file to local client's storage
			err = app.Saver.SaveFile(id.String(), pathStr, info, etag)
			if err != nil {
				logger.Error("error saving file locally", zap.Error(err))
			}

			logger.Info("Your file saved successfully. Please keep your uuid and use it to retrive your data back from Goph-keeper.",
				zap.String("uuid:", id.String()))
		},
	}

	cmd.Flags().StringP("file_name", "n", "", "Name of the file to be saved")
	cmd.Flags().StringP("file_path", "p", "", "Path to the binary file, which need to be saved")
	cmd.Flags().StringP("info", "i", "", "Additional metadata, if necessary")

	return cmd
}

func downloadBinCmd(app *App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bin",
		Short: "Download binary data from storage",
		Run: func(cmd *cobra.Command, args []string) {
			idStr, err := cmd.Flags().GetString("id")
			if err != nil {
				logger.Fatal("failed to get file uuid:", zap.Error(err))
			}

			fileNameStr, err := cmd.Flags().GetString("file_name")
			if err != nil {
				logger.Fatal("failed to get file_name:", zap.Error(err))
			}

			// Initializing Download service
			clientService := serviceDown.New()

			serverAddr, _ := viper.Get("GRPC_PORT").(string)

			if err := clientService.DownloadFile(fmt.Sprintf(":%s", serverAddr), idStr, fileNameStr); err != nil {
				logger.Error("failed to obtain requested binary data from goph-keeper: ", zap.Error(err))

				// if remote server not responding, try to reach local storage
				res, err := app.GetFile(idStr)
				if err != nil {
					logger.Error("failed to obtain requested binary data from goph-keeper: ", zap.Error(err))
				}

				err = fl.SaveFileToDisk(res, "client_files")
				if err != nil {
					logger.Error("file saved locally")
				}

			}

		},
	}

	cmd.Flags().StringP("id", "i", "", "A Universally Unique Identifier of needed binary")
	cmd.Flags().StringP("file_name", "n", "", "Name of the file")

	return cmd
}
