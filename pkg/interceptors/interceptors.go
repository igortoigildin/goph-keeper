package interceptors

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/golang-jwt/jwt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func JwtUnaryInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		if strings.HasSuffix(info.FullMethod, "/Login") {
			return handler(ctx, req)
		}

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, errors.New("missing metadata")
		}

		tokens := md.Get("authorization")
		if len(tokens) == 0 {
			return nil, errors.New("authorization token not provided")
		}

		tokenStr := strings.TrimPrefix(tokens[0], "Bearer ")

		jwtSecret := os.Getenv("JWT_SECRET")
		// Decode the base64 secret key
		secretKey, err := base64.StdEncoding.DecodeString(jwtSecret)
		if err != nil {
			return nil, fmt.Errorf("invalid secret key format: %v", err)
		}

		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			return secretKey, nil
		})

		if err != nil || !token.Valid {
			return nil, fmt.Errorf("invalid token: %v", err)
		}

		return handler(ctx, req)
	}
}

// Интерсептор для проверки JWT в стриминговых запросах
func JwtStreamInterceptor() grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		ss grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		// Пропускаем проверку токена для метода Login
		if strings.HasSuffix(info.FullMethod, "/Login") {
			return handler(srv, ss)
		}

		md, ok := metadata.FromIncomingContext(ss.Context())
		if !ok {
			return fmt.Errorf("missing metadata")
		}

		tokens := md.Get("authorization")
		if len(tokens) == 0 {
			return fmt.Errorf("authorization token not provided")
		}

		tokenStr := strings.TrimPrefix(tokens[0], "Bearer ")

		jwtSecret := os.Getenv("JWT_SECRET")
		// Decode the base64 secret key
		secretKey, err := base64.StdEncoding.DecodeString(jwtSecret)
		if err != nil {
			return fmt.Errorf("invalid secret key format: %v", err)
		}

		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return secretKey, nil
		})

		if err != nil || !token.Valid {
			return fmt.Errorf("invalid token: %v", err)
		}

		return handler(srv, ss)
	}
}
