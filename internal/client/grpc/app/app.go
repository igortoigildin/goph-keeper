package app

import (
	"context"
	"log"
	"os"
	"time"

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
	sessionDuration = time.Minute * 7
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
		emailStr, err := cmd.Flags().GetString("email")
		if err != nil {
			log.Fatalf("failed to get email: %s\n", err.Error())
		}

		passStr, err := cmd.Flags().GetString("password")
		if err != nil {
			log.Fatalf("failed to get password: %s\n", err.Error())
		}

		authService := authService.New(serverAddr)
		token, err := authService.Login(context.Background(), emailStr, passStr)
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
			Email: emailStr,
			Token: token,
			ExpiresAt: time.Now().Add(sessionDuration),
		}

		err = session.SaveSession(sessionData)
		if err != nil {
			logger.Error("failed to save sesson", zap.Error(err))
		}

		log.Printf("user with %s email logged in successfully. Session saved\n", emailStr)
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
		// Check if the user is set (authenticated)
		user, _ := cmd.Flags().GetString("user")
		if user == "" {
			log.Println("You must be logged in to run this command")
			os.Exit(1)
		}
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

		log.Printf("login %s && password %s saved successfully\n", loginStr, passStr)
	},
}

// save text subcommand
var saveTextCmd = &cobra.Command{
	Use:   "text",
	Short: "Save arbitrary text data in storage",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Check if the user is set (authenticated)
		user, _ := cmd.Flags().GetString("user")
		if user == "" {
			log.Println("You must be logged in to run this command")
			os.Exit(1)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		_, err := cmd.Flags().GetString("file_name")
		if err != nil {
			log.Fatalf("failed to get file_name: %s\n", err.Error())
		}

		_, err = cmd.Flags().GetString("text")
		if err != nil {
			log.Fatalf("failed to get text: %s\n", err.Error())
		}

		// TODO: save text in minio

		log.Println("text saved successfully\n")
	},
}

// save binary data subcommand
var saveBinCmd = &cobra.Command{
	Use:   "bin",
	Short: "Save binary data in storage",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if !session.IsSessionValid(refreshTokenSecretKey)  {
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

		clientService := service.New(serverAddr, pathStr, batchSize)

		if err := clientService.SendFile(); err != nil {
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
		// Check if the user is set (authenticated)
		user, _ := cmd.Flags().GetString("user")
		if user == "" {
			log.Println("You must be logged in to run this command")
			os.Exit(1)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		_, err := cmd.Flags().GetString("card_number")
		if err != nil {
			log.Fatalf("failed to get card_number: %s\n", err.Error())
		}

		_, err = cmd.Flags().GetString("CVC")
		if err != nil {
			log.Fatalf("failed to get path: %s\n", err.Error())
		}

		_, err = cmd.Flags().GetString("expiration_date")
		if err != nil {
			log.Fatalf("failed to get path: %s\n", err.Error())
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
	createUserCmd.Flags().StringP("email", "e", "", "User email")
	createUserCmd.Flags().StringP("password", "p", "", "User password")
	createUserCmd.Flags().StringVarP(&serverAddr, "addr", "a", "", "server address")

	rootCmd.AddCommand(loginCmd)
	loginCmd.AddCommand(loginUserCmd)
	loginUserCmd.Flags().StringP("email", "e", "", "User email")
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
	saveTextCmd.Flags().StringP("file_name", "n", "", "Provided text will be saved in stated file")
	saveTextCmd.Flags().StringP("text", "t", "", "Text which need to be saved")

	// save binary data
	saveCmd.AddCommand(saveBinCmd)
	saveBinCmd.Flags().StringP("file_name", "n", "", "Name of the file to be saved")
	saveBinCmd.Flags().StringP("file_path", "p", "", "Path to the binary file, which need to be saved")

	// save card data
	saveCmd.AddCommand(saveCardInfoCmd)
	saveCardInfoCmd.Flags().StringP("card_number", "n", "", "Card number to be saved")
	saveCardInfoCmd.Flags().StringP("CVC", "c", "", "CVC to be saved")

	logger.Initialize(loggerLevel)

	if err := createUserCmd.MarkFlagRequired("addr"); err != nil {
		log.Fatal(err)
	}
}

// var (
// 	serverAddr  string
// 	filePath    string
// 	batchSize   int
// 	loggerLevel string
// 	rootCmd     = &cobra.Command{
// 		Use:   "transfer_client",
// 		Short: "Sending files via gRPC",
// 		Run: func(cmd *cobra.Command, args []string) {
// 			clientService := service.New(serverAddr, filePath, batchSize)

// 			if err := clientService.SendFile(); err != nil {
// 				log.Fatal(err)
// 			}
// 		},
// 	}
// )
