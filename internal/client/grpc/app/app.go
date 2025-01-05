package app

import (
	"context"
	"log"
	"os"

	authService "github.com/igortoigildin/goph-keeper/internal/client/grpc/service/register"
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

// registration
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
			log.Fatalf("failed to login: %s\n", err.Error())
		}

		log.Printf("user with %s email created successfully\n", emailStr)
	},
}

// login
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

		_, err = cmd.Flags().GetString("password")
		if err != nil {
			log.Fatalf("failed to get password: %s\n", err.Error())
		}

		log.Printf("user with %s email logged in successfully\n", emailStr)
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
