package app

import (
	"context"
	"log"
	"os"

	authService "github.com/igortoigildin/goph-keeper/internal/client/grpc/service/auth"
	"github.com/igortoigildin/goph-keeper/pkg/logger"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

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

var (
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
		if err = authService.Login(context.Background(), emailStr, passStr); err != nil {
			log.Fatalf("failed to login: %s\n", err.Error())
		}

		log.Printf("user with %s email logged in successfully\n", emailStr)
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

	Run: func(cmd *cobra.Command, args []string) {
		loginStr, err := cmd.Flags().GetString("login")
		if err != nil {
			log.Fatalf("failed to get login: %s\n", err.Error())
		}

		passStr, err := cmd.Flags().GetString("password")
		if err != nil {
			log.Fatalf("failed to get password: %s\n", err.Error())
		}

		// TODO: save login && pass in minio

		log.Printf("login %s && password %s saved successfully\n", loginStr, passStr)
	},
}

// save text subcommand
var saveTextCmd = &cobra.Command{
	Use:   "text",
	Short: "Save arbitrary text data in storage",

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

// save bin data subcommand
var saveBinCmd = &cobra.Command{
	Use:   "bin",
	Short: "Save binary data in storage",

	Run: func(cmd *cobra.Command, args []string) {
		_, err := cmd.Flags().GetString("file_name")
		if err != nil {
			log.Fatalf("failed to get file_name: %s\n", err.Error())
		}

		pathStr, err := cmd.Flags().GetString("path")
		if err != nil {
			log.Fatalf("failed to get path: %s\n", err.Error())
		}

		// TODO: save binary data in minio

		log.Printf("biniry data %s saved successfully\n", pathStr)
	},
}

// save card bank details subcommand
var saveCardInfoCmd = &cobra.Command{
	Use:   "card",
	Short: "Save bank card details in storage",

	Run: func(cmd *cobra.Command, args []string) {
		_, err := cmd.Flags().GetString("card_number")
		if err != nil {
			log.Fatalf("failed to get card_number: %s\n", err.Error())
		}

		_, err = cmd.Flags().GetString("CVC")
		if err != nil {
			log.Fatalf("failed to get path: %s\n", err.Error())
		}

		_, err = cmd.Flags().GetString("Expiration date")
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
	//rootCmd.Flags().StringVarP(&serverAddr, "addr", "a", "", "server address")
	rootCmd.Flags().StringVarP(&loggerLevel, "log", "l", "info", "logger level")
	// rootCmd.Flags().StringVarP(&filePath, "file", "f", "", "file path")
	// rootCmd.Flags().IntVarP(&batchSize, "batch", "b", 1024*1024, "batch size for sending")
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

	logger.Initialize(loggerLevel)

	// if err := rootCmd.MarkFlagRequired("file"); err != nil {
	// 	log.Fatal(err)
	// }

	if err := createUserCmd.MarkFlagRequired("addr"); err != nil {
		log.Fatal(err)
	}
}
