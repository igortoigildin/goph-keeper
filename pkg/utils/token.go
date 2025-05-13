package jwt

import (
	"encoding/base64"
	"time"

	"github.com/golang-jwt/jwt"
	model "github.com/igortoigildin/goph-keeper/internal/server/models"
	"github.com/pkg/errors"
)

func GenerateToken(info model.UserInfo, secretKey []byte, duration time.Duration) (string, error) {
	// Decode the base64 secret key
	decodedKey, err := base64.StdEncoding.DecodeString(string(secretKey))
	if err != nil {
		return "", errors.Errorf("invalid secret key format: %v", err)
	}

	claims := model.UserClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(duration).Unix(),
		},
		Login: info.Login,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(decodedKey)
}

func VeryfyToken(tokenStr string, secretKey []byte) (*model.UserClaims, error) {
	// Decode the base64 secret key
	decodedKey, err := base64.StdEncoding.DecodeString(string(secretKey))
	if err != nil {
		return nil, errors.Errorf("invalid secret key format: %v", err)
	}

	token, err := jwt.ParseWithClaims(
		tokenStr,
		&model.UserClaims{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.Errorf("unexpected token signing method")
			}

			return decodedKey, nil
		},
	)

	if err != nil {
		return nil, errors.Errorf("invalid token: %s", err.Error())
	}

	claims, ok := token.Claims.(*model.UserClaims)
	if !ok {
		return nil, errors.Errorf("invalid token claims")
	}

	return claims, nil
}
