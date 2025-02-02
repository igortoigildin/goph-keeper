package app

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
	authService "github.com/igortoigildin/goph-keeper/internal/client/grpc/service/auth"
	service "github.com/igortoigildin/goph-keeper/internal/client/grpc/service/upload"
	"github.com/igortoigildin/goph-keeper/pkg/logger"
	"github.com/igortoigildin/goph-keeper/pkg/session"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var (
	// File to store JWT token
	tokenFile             = ".jwt_token"
	refreshTokenSecretKey = "W4/X+LLjehdxptt4YgGFCvMpq5ewptpZZYRHY6A72g0="
	accessTokenSecretKey  = "VqvguGiffXILza1f44TWXowDT4zwf03dtXmqWW4SYyE="
	sessionDuration       = time.Minute * 7
)

var (
	batchSize   int
	loggerLevel string
	serverAddr  string
	rootCmd     = &cobra.Command{
		Use:   "goph-keeper-app",
		Short: "My cli app",
	}
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
		emailStr, err := cmd.Flags().GetString("email")
		if err != nil {
			log.Fatalf("failed to get email: %s\n", err.Error())
		}

		passStr, err := cmd.Flags().GetString("password")
		if err != nil {
			log.Fatalf("failed to get password: %s\n", err.Error())
		}

		authService := authService.New(serverAddr)

		if err = authService.RegisterNewUser(context.Background(), emailStr, passStr); err != nil {
			log.Fatalf("registration failed: %s\n", err.Error())
		}

		log.Printf("user with %s email created successfully\n", emailStr)
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
	Run: func(cmd *cobra.Command, args []string) {
		loginStr, err := cmd.Flags().GetString("login")
		if err != nil {
			log.Fatalf("failed to get email: %s\n", err.Error())
		}

		passStr, err := cmd.Flags().GetString("password")
		if err != nil {
			log.Fatalf("failed to get password: %s\n", err.Error())
		}

		authService := authService.New(serverAddr)
		token, err := authService.Login(context.Background(), loginStr, passStr)
		if err != nil {
			log.Fatalf("failed to login: %s\n", err.Error())
		} else if token == "" {
			log.Fatalf("got invalid jwt token: %s\n", err.Error())
		}

		// err = os.WriteFile(tokenFile, []byte(token), 0644)
		// if err != nil {
		// 	logger.Error("error saving JWT token", zap.Error(err))
		// 	return
		// }

		sessionData := &session.Session{
			Login:     loginStr,
			Token:     token,
			ExpiresAt: time.Now().Add(sessionDuration),
		}

		err = session.SaveSession(sessionData)
		if err != nil {
			logger.Error("failed to save sesson", zap.Error(err))
		}

		log.Printf("user with %s email logged in successfully. Session saved\n", loginStr)
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
		if !session.IsSessionValid(refreshTokenSecretKey) {
			logger.Error("Session expired or not found. Please login again")

			os.Exit(1)
			return
		}

		logger.Info("Session is valid")
	},
	Run: func(cmd *cobra.Command, args []string) {
		loginStr, err := cmd.Flags().GetString("login")
		if err != nil {
			log.Fatalf("failed to get login: %s\n", err.Error())
		}

		passStr, err := cmd.Flags().GetString("password")
		if err != nil {
			log.Fatalf("failed to get password: %s\n", err.Error())
		}

		serverAddr = ":9000" // TO BE UPDATED

		clientService := service.New()

		id := uuid.New()

		if err := clientService.SendPassword(serverAddr, loginStr, passStr, id.String()); err != nil {
			log.Fatal("failed to send password: ", zap.Error(err))
		}

		log.Printf("login %s && password %s saved successfully\n", loginStr, passStr)
	},
}

// save text subcommand
var saveTextCmd = &cobra.Command{
	Use:   "text",
	Short: "Save arbitrary text data in storage",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if !session.IsSessionValid(refreshTokenSecretKey) {
			logger.Error("Session expired or not found. Please login again")

			os.Exit(1)
			return
		}

		logger.Info("Session is valid")
	},
	Run: func(cmd *cobra.Command, args []string) {
		// _, err := cmd.Flags().GetString("file_name")
		// if err != nil {
		// 	log.Fatalf("failed to get file_name: %s\n", err.Error())
		// }

		textData, err := cmd.Flags().GetString("text")
		if err != nil {
			log.Fatalf("failed to get text: %s\n", err.Error())
		}

		serverAddr = ":9000" // TO BE UPDATED

		id := uuid.New()

		clientService := service.New()

		if err := clientService.SendText(serverAddr, textData, id.String()); err != nil {
			log.Fatal("failed to send text: ", zap.Error(err))
		}

		log.Println("text saved successfully\n")
	},
}

// save binary data subcommand
var saveBinCmd = &cobra.Command{
	Use:   "bin",
	Short: "Save binary data in storage",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if !session.IsSessionValid(refreshTokenSecretKey) {
			logger.Error("Session expired or not found. Please login again")

			os.Exit(1)
			return
		}

		logger.Info("Session is valid")
	},
	Run: func(cmd *cobra.Command, args []string) {
		pathStr, err := cmd.Flags().GetString("file_path")
		if err != nil {
			log.Fatalf("failed to get path: %s\n", err.Error())
		}

		serverAddr = ":9000" // TO BE UPDATED

		id := uuid.New()

		clientService := service.New()

		if err := clientService.SendFile(serverAddr, pathStr, batchSize, id.String()); err != nil {
			log.Fatal("failed to send binary file: ", zap.Error(err))
		}

		log.Printf("biniry data %s saved successfully\n", pathStr)
	},
}

// save card bank details subcommand
var saveCardInfoCmd = &cobra.Command{
	Use:   "card",
	Short: "Save bank card details in storage",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if !session.IsSessionValid(refreshTokenSecretKey) {
			logger.Error("Session expired or not found. Please login again")

			os.Exit(1)
			return
		}

		logger.Info("Session is valid")
	},
	Run: func(cmd *cobra.Command, args []string) {
		cardNumber, err := cmd.Flags().GetString("card_number")
		if err != nil {
			log.Fatalf("failed to get card_number: %s\n", err.Error())
		}

		cvc, err := cmd.Flags().GetString("CVC")
		if err != nil {
			log.Fatalf("failed to get path: %s\n", err.Error())
		}

		expDate, err := cmd.Flags().GetString("expiration_date")
		if err != nil {
			log.Fatalf("failed to get path: %s\n", err.Error())
		}

		serverAddr = ":9000" // TO BE UPDATED

		id := uuid.New()

		clientService := service.New()

		if err := clientService.SendBankDetails(serverAddr, cardNumber, cvc, expDate, id.String()); err != nil {
			log.Fatal("failed to send binary file: ", zap.Error(err))
		}

		// TODO: save bank card details in minio

		log.Println("bank card details saved successfully\n")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		logger.Error("error while executing root cmd", zap.Error(err))
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVarP(&loggerLevel, "log", "l", "info", "logger level")
	// rootCmd.Flags().StringVarP(&filePath, "file", "f", "", "file path")
	rootCmd.Flags().IntVarP(&batchSize, "batch", "b", 1024*1024, "batch size for sending")
	rootCmd.AddCommand(createCmd)
	createCmd.AddCommand(createUserCmd)
	createUserCmd.Flags().StringP("login", "l", "", "User login")
	createUserCmd.Flags().StringP("password", "p", "", "User password")
	createUserCmd.Flags().StringVarP(&serverAddr, "addr", "a", "", "server address")

	rootCmd.AddCommand(loginCmd)
	loginCmd.AddCommand(loginUserCmd)
	loginUserCmd.Flags().StringP("login", "l", "", "User login")
	loginUserCmd.Flags().StringP("password", "p", "", "User password")
	loginUserCmd.Flags().StringVarP(&serverAddr, "addr", "a", "", "server address")

	rootCmd.AddCommand(saveCmd)

	// save login && password
	saveCmd.AddCommand(savePasswordCmd)
	savePasswordCmd.Flags().StringP("login", "l", "", "Login to be saved")
	savePasswordCmd.Flags().StringP("password", "p", "", "Password to be saved")
	savePasswordCmd.Flags().StringVarP(&serverAddr, "addr", "a", "", "server address")

	// save text data
	saveCmd.AddCommand(saveTextCmd)
	saveTextCmd.Flags().StringP("text", "t", "", "Text which need to be saved")

	// save binary data
	saveCmd.AddCommand(saveBinCmd)
	saveBinCmd.Flags().StringP("file_name", "n", "", "Name of the file to be saved")
	saveBinCmd.Flags().StringP("file_path", "p", "", "Path to the binary file, which need to be saved")

	// save card data
	saveCmd.AddCommand(saveCardInfoCmd)
	saveCardInfoCmd.Flags().StringP("card_number", "n", "", "Card number to be saved")
	saveCardInfoCmd.Flags().StringP("CVC", "c", "", "CVC to be saved")
	saveCardInfoCmd.Flags().StringP("expiration_date", "e", "", "expiration_date to be saved")

	logger.Initialize(loggerLevel)

	if err := createUserCmd.MarkFlagRequired("addr"); err != nil {
		log.Fatal(err)
	}
}
