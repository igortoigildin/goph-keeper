package app

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/igortoigildin/goph-keeper/internal/client/config"
	authService "github.com/igortoigildin/goph-keeper/internal/client/grpc/service/auth"
	serviceDown "github.com/igortoigildin/goph-keeper/internal/client/grpc/service/download"
	serviceUp "github.com/igortoigildin/goph-keeper/internal/client/grpc/service/upload"
	"github.com/igortoigildin/goph-keeper/pkg/logger"
	"github.com/igortoigildin/goph-keeper/pkg/session"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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

// user registration
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "New user registration",
}

var createUserCmd = &cobra.Command{
	Use:   "user",
	Short: "New user registration",
	Run: func(cmd *cobra.Command, args []string) {
		loginStr, err := cmd.Flags().GetString("login")
		if err != nil {
			logger.Error("failed to get login:", zap.Error(err))

			return
		}

		passStr, err := cmd.Flags().GetString("password")
		if err != nil {
			logger.Error("failed to get password:", zap.Error(err))

			return
		}

		serverAddr, _ := viper.Get("GRPC_PORT").(string)
		authService := authService.New(fmt.Sprintf(":%s", serverAddr))

		if err = authService.RegisterNewUser(context.Background(), loginStr, passStr); err != nil {
			logger.Error("registration failed:", zap.Error(err))

			return
		}

		logger.Info("User created successfully:", zap.String("login", loginStr))
	},
}

// user login
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "User authentication",
}

var loginUserCmd = &cobra.Command{
	Use:   "user",
	Short: "User authentication",
	RunE: func(cmd *cobra.Command, args []string) error {
		loginStr, err := cmd.Flags().GetString("login")
		if err != nil {
			logger.Error("failed to get login:", zap.Error(err))

			return fmt.Errorf("login not provided: %w", err)
		}

		passStr, err := cmd.Flags().GetString("password")
		if err != nil {
			logger.Error("failed to get password:", zap.Error(err))

			return fmt.Errorf("password not provided: %w", err)
		}

		serverAddr, _ := viper.Get("GRPC_PORT").(string)

		authService := authService.New(fmt.Sprintf(":%s", serverAddr))

		token, err := authService.Login(context.Background(), loginStr, passStr)
		if err != nil {
			logger.Error("failed to login:", zap.Error(err))

			return fmt.Errorf("authentication error: %w", err)
		}

		if token == "" {
			logger.Error("failed to login, jwt token has not been received")

			return fmt.Errorf("authentication error: %w", err)
		}

		sessionData := &session.Session{
			Login:     loginStr,
			Token:     token,
			ExpiresAt: time.Now().Add(sessionDuration),
		}

		err = session.SaveSession(sessionData)
		if err != nil {
			logger.Error("failed to save sesson", zap.Error(err))

			return fmt.Errorf("failed to save sesson: %w", err)
		}

		logger.Info("Session saved. User logged in successfully:", zap.String("login", loginStr))

		return nil
	},
}

// download command
var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Download data from storage",
}

// download password subcommmand
var downloadPassCmd = &cobra.Command{
	Use:   "password",
	Short: "Download login && password from storage",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		refreshTokenSecretKey, _ := viper.Get("REFRESH_SECRET").(string)

		if !session.IsSessionValid(refreshTokenSecretKey) {
			logger.Fatal("Session expired or not found. Please login again")
		}

		logger.Info("Session is valid")
	},
	Run: func(cmd *cobra.Command, args []string) {
		idStr, err := cmd.Flags().GetString("id")
		if err != nil {
			logger.Fatal("failed to get credentials id", zap.Error(err))
		}

		clientService := serviceDown.New()
		serverAddr, _ := viper.Get("GRPC_PORT").(string)

		if err := clientService.DownloadPassword(fmt.Sprintf(":%s", serverAddr), idStr); err != nil {
			logger.Error("failed to obtain requested credentials from goph-keeper", zap.Error(err))
		}
	},
}

// save command
var saveCmd = &cobra.Command{
	Use:   "save",
	Short: "Save data in storage",
}

// save password subcommand
var savePasswordCmd = &cobra.Command{
	Use:   "password",
	Short: "Save login && password in storage",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		refreshTokenSecretKey, _ := viper.Get("REFRESH_SECRET").(string)

		if !session.IsSessionValid(refreshTokenSecretKey) {
			logger.Fatal("Session expired or not found. Please login again")
		}

		logger.Info("Session is valid")
	},
	Run: func(cmd *cobra.Command, args []string) {
		loginStr, err := cmd.Flags().GetString("login")
		if err != nil {
			logger.Fatal("failed to get login:", zap.Error(err))
		}

		passStr, err := cmd.Flags().GetString("password")
		if err != nil {
			logger.Fatal("failed to get password:", zap.Error(err))
		}

		meta, err := cmd.Flags().GetString("destination")
		if err != nil {
			logger.Fatal("failed to get metadata", zap.Error(err))
		}

		// Initializing Upload service.
		clientService := serviceUp.New()

		// Creating new uuid for credentials to be saved.
		id := uuid.New()

		serverAddr, _ := viper.Get("GRPC_PORT").(string)

		// Sending credentials with created uuid to server.
		if err := clientService.SendPassword(fmt.Sprintf(":%s", serverAddr), loginStr, passStr, id.String(), meta); err != nil {
			logger.Error("failed to send credentials to server:", zap.Error(err))
		}

		logger.Info("Credentials saved successfully. Please save your uuid and use it to retrive your data back from Goph-keeper.",
			zap.String("uuid:", id.String()))
	},
}

// download text subcommand
var downloadTextCmd = &cobra.Command{
	Use:   "text",
	Short: "Download arbitrary text data from storage",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		refreshTokenSecretKey, _ := viper.Get("REFRESH_SECRET").(string)

		if !session.IsSessionValid(refreshTokenSecretKey) {
			logger.Fatal("Session expired or not found. Please login again")
		}

		logger.Info("Session is valid")
	},
	Run: func(cmd *cobra.Command, args []string) {
		idStr, err := cmd.Flags().GetString("id")
		if err != nil {
			logger.Fatal("failed to get text uuid:", zap.Error(err))
		}

		// Initializing download service.
		clientService := serviceDown.New()

		serverAddr, _ := viper.Get("GRPC_PORT").(string)

		// Requesting text with provided uuid.
		if err := clientService.DownloadText(fmt.Sprintf(":%s", serverAddr), idStr); err != nil {
			logger.Fatal("failed to obtain text data from goph-keeper: ", zap.Error(err))
		}
	},
}

// save text subcommand
var saveTextCmd = &cobra.Command{
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

		meta, err := cmd.Flags().GetString("info")
		if err != nil {
			logger.Fatal("failed to get metadata", zap.Error(err))
		}

		// Creating new uuid for text to be saved
		id := uuid.New()

		// Initializing Upload service
		clientService := serviceUp.New()

		serverAddr, _ := viper.Get("GRPC_PORT").(string)

		if err := clientService.SendText(fmt.Sprintf(":%s", serverAddr), textData, id.String(), meta); err != nil {
			logger.Fatal("failed to save text", zap.Error(err))
		}

		logger.Info("Your text saved successfully. Please save your uuid and use it to retrive your data back from Goph-keeper.",
			zap.String("uuid:", id.String()))
	},
}

// download binary data subcommand
var downloadBinCmd = &cobra.Command{
	Use:   "bin",
	Short: "Download binary data from storage",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		refreshTokenSecretKey, _ := viper.Get("REFRESH_SECRET").(string)

		if !session.IsSessionValid(refreshTokenSecretKey) {
			logger.Fatal("Session expired or not found. Please login again")
		}

		logger.Info("Session is valid")
	},
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
			logger.Fatal("failed to obtain requested binary data from goph-keeper: ", zap.Error(err))
		}
	},
}

// save binary data subcommand
var saveBinCmd = &cobra.Command{
	Use:   "bin",
	Short: "Save binary data in storage",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		refreshTokenSecretKey, _ := viper.Get("REFRESH_SECRET").(string)

		if !session.IsSessionValid(refreshTokenSecretKey) {
			logger.Fatal("Session expired or not found. Please login again")
		}

		logger.Info("Session is valid")
	},
	Run: func(cmd *cobra.Command, args []string) {
		pathStr, err := cmd.Flags().GetString("file_path")
		if err != nil {
			log.Fatalf("failed to get path: %s\n", err.Error())
		}

		meta, err := cmd.Flags().GetString("info")
		if err != nil {
			logger.Fatal("failed to get metadata", zap.Error(err))
		}

		// Creating new uuid for the file to be saved
		id := uuid.New()

		// Creating Upload service
		clientService := serviceUp.New()

		serverAddr, _ := viper.Get("GRPC_PORT").(string)

		if err := clientService.SendFile(fmt.Sprintf(":%s", serverAddr), pathStr, batchSize, id.String(), meta); err != nil {
			logger.Fatal("failed to save binary file: ", zap.Error(err))
		}

		logger.Info("Your file saved successfully. Please keep your uuid and use it to retrive your data back from Goph-keeper.",
			zap.String("uuid:", id.String()))
	},
}

// download card details subcommand
var downloadCardInfoCmd = &cobra.Command{
	Use:   "card",
	Short: "Download card details from storage",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		refreshTokenSecretKey, _ := viper.Get("REFRESH_SECRET").(string)

		if !session.IsSessionValid(refreshTokenSecretKey) {
			logger.Fatal("Session expired or not found. Please login again")
		}

		logger.Info("Session is valid")
	},
	Run: func(cmd *cobra.Command, args []string) {
		idStr, err := cmd.Flags().GetString("id")
		if err != nil {
			logger.Fatal("failed to get bank details uuid:", zap.Error(err))
		}

		// Initializing Download service
		clientService := serviceDown.New()

		serverAddr, _ := viper.Get("GRPC_PORT").(string)

		if err := clientService.DownloadBankDetails(fmt.Sprintf(":%s", serverAddr), idStr); err != nil {
			logger.Fatal("failed to obtain card details from goph-keeper: ", zap.Error(err))
		}

	},
}

// save card bank details subcommand
var saveCardInfoCmd = &cobra.Command{
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

		if err := clientService.SendBankDetails(fmt.Sprintf(":%s", serverAddr), cardNumber, cvc, expDate, id.String(), meta); err != nil {
			logger.Fatal("failed to save bank details: ", zap.Error(err))
		}

		logger.Info("Your bank details saved successfully. Please keep your uuid and use it to retrive your data back from Goph-keeper.",
			zap.String("uuid:", id.String()))
	},
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
	saveCmd.AddCommand(savePasswordCmd)
	savePasswordCmd.Flags().StringP("login", "l", "", "Login to be saved")
	savePasswordCmd.Flags().StringP("password", "p", "", "Password to be saved")
	savePasswordCmd.Flags().StringP("destination", "d", "", `Name of the site, app, or other platform
	 for which the login and password were created.`)

	rootCmd.AddCommand(downloadCmd)

	// download login && password
	downloadCmd.AddCommand(downloadPassCmd)
	downloadPassCmd.Flags().StringP("id", "i", "", "A Universally Unique Identifier of the saved password")

	// save text data
	saveCmd.AddCommand(saveTextCmd)
	saveTextCmd.Flags().StringP("text", "t", "", "Text which need to be saved")
	saveTextCmd.Flags().StringP("info", "i", "", "Additional metadata, if necessary")

	// download text
	downloadCmd.AddCommand(downloadTextCmd)
	downloadTextCmd.Flags().StringP("id", "i", "", "A Universally Unique Identifier of saved text")

	// save binary data
	saveCmd.AddCommand(saveBinCmd)
	saveBinCmd.Flags().StringP("file_name", "n", "", "Name of the file to be saved")
	saveBinCmd.Flags().StringP("file_path", "p", "", "Path to the binary file, which need to be saved")
	saveBinCmd.Flags().StringP("info", "i", "", "Additional metadata, if necessary")

	// download binary data
	downloadCmd.AddCommand(downloadBinCmd)
	downloadBinCmd.Flags().StringP("id", "i", "", "A Universally Unique Identifier of needed binary")
	downloadBinCmd.Flags().StringP("file_name", "n", "", "Name of the file")

	// save card details
	saveCmd.AddCommand(saveCardInfoCmd)
	saveCardInfoCmd.Flags().StringP("card_number", "n", "", "Card number to be saved")
	saveCardInfoCmd.Flags().StringP("CVC", "c", "", "CVC to be saved")
	saveCardInfoCmd.Flags().StringP("expiration_date", "e", "", "expiration_date to be saved")
	saveCardInfoCmd.Flags().StringP("info", "i", "", "Additional metadata, if necessary")

	// download card details
	downloadCmd.AddCommand(downloadCardInfoCmd)
	downloadCardInfoCmd.Flags().StringP("id", "i", "", "A Universally Unique Identifier of the saved card details")
}
