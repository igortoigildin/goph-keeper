package app

import (
	"context"
	"fmt"
	"time"

	authService "github.com/igortoigildin/goph-keeper/internal/client/grpc/service/auth"
	"github.com/igortoigildin/goph-keeper/pkg/logger"
	"github.com/igortoigildin/goph-keeper/pkg/session"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
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
