package app

import (
	"fmt"
	"os"
	"time"

	"github.com/igortoigildin/goph-keeper/internal/client/config"
	"github.com/igortoigildin/goph-keeper/internal/client/grpc/models"
	syncService "github.com/igortoigildin/goph-keeper/internal/client/grpc/service/sync"
	storage "github.com/igortoigildin/goph-keeper/internal/client/grpc/storage/sqlite"
	"github.com/igortoigildin/goph-keeper/pkg/logger"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var (
	loggerLevel string
	rootCmd     = &cobra.Command{
		Use:   "goph-keeper-app",
		Short: "My cli app",
	}
	sessionDuration = time.Minute * 7
	batchSize       = 1024 * 1024
)

type App struct {
	DBPath string
	Saver
	Downloader
	Receiver
	Syncer
}

type Syncer interface {
	SyncAllData(addr string) error
}

type Saver interface {
	SaveText(id, info, text, etag string) error
	SaveCredentials(id, service, username, password, etag string) error
	SaveBankDetails(cardNumber, cvc, expDate, id, bankName, etag string) error
	SaveFile(id, filePath, info, etag string) error
}

type Downloader interface {
	GetAllTexts() ([]models.Text, error)
	GetText(id string) (models.Text, error)
	GetAllCredentials() ([]models.Credential, error)
	GetCredential(id string) (models.Credential, error)
	GetAllBankDetails() ([]models.BankDetails, error)
	GetBankDetails(id string) (models.BankDetails, error)
	GetFile(id string) (models.File, error)
}

type Receiver interface {
	GetAllTexts() ([]models.Text, error)
	GetAllCredentials() ([]models.Credential, error)
	GetAllBankDetails() ([]models.BankDetails, error)
	ListAllFiles() ([]models.File, error)
}

func NewApp(dbPath string) (*App, error) {
	storage, err := storage.NewClientRepository(dbPath)
	if err != nil {
		return nil, fmt.Errorf("Failed to connect to DB: %w", err)
	}

	return &App{
		Saver:      storage,
		Downloader: storage,
		DBPath:     dbPath,
		Receiver:   storage,
		Syncer:     syncService.New(),
	}, nil
}

// save command
var saveCmd = &cobra.Command{
	Use:   "save",
	Short: "Save data in storage",
}

// download command
var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Download data from storage",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		logger.Error("error executing root cmd", zap.Error(err))

		fmt.Println("Error:", err)
		os.Exit(1)
	}
}

func init() {
	logger.Initialize(loggerLevel)

	if err := config.LoadConfig(); err != nil {
		logger.Fatal("error loading config", zap.Error(err))
	}

	app, err := NewApp("sqlite3")
	if err != nil {
		logger.Error("Failed to init app", zap.Error(err))
		os.Exit(1)
	}

	rootCmd.Flags().StringVarP(&loggerLevel, "log", "l", "info", "logger level")
	rootCmd.AddCommand(createCmd)
	createCmd.AddCommand(createUserCmd)
	createUserCmd.Flags().StringP("login", "l", "", "User login")
	createUserCmd.Flags().StringP("password", "p", "", "User password")

	rootCmd.AddCommand(loginCmd)
	loginCmd.AddCommand(loginUserCmd)
	loginUserCmd.Flags().StringP("login", "l", "", "User login")
	loginUserCmd.Flags().StringP("password", "p", "", "User password")

	rootCmd.AddCommand(saveCmd)

	// save login && password
	saveCmd.AddCommand(savePasswordCmd(app))

	rootCmd.AddCommand(downloadCmd)

	// download login && password
	downloadCmd.AddCommand(downloadPassCmd(app))

	// save text data
	saveCmd.AddCommand(saveTextCmd(app))

	// download text
	downloadCmd.AddCommand(downloadTextCmd(app))

	// save binary data
	saveCmd.AddCommand(saveBinCmd(app))

	// download binary data
	downloadCmd.AddCommand(downloadBinCmd(app))

	// save card details
	saveCmd.AddCommand(saveCardInfoCmd(app))

	// download card details
	downloadCmd.AddCommand(downloadCardInfoCmd(app))

	// list all saved secrets
	listCmd.AddCommand(listAllSavedSecrets(app))

	rootCmd.AddCommand(listCmd)

	syncCmd.AddCommand(syncAllData(app))

	// sync data with server
	rootCmd.AddCommand(syncCmd)
}
