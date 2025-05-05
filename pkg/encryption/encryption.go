package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"strings"

	"github.com/igortoigildin/goph-keeper/pkg/logger"
	"go.uber.org/zap"
)

func Encrypt(plaintext string, key []byte) (string, error) {
	// Decode the base64 key
	decodedKey, err := base64.StdEncoding.DecodeString(string(key))
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(decodedKey)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func Decrypt(encoded string, key []byte) (string, error) {
	logger.Debug("Starting decryption",
		zap.String("encoded_length", fmt.Sprintf("%d", len(encoded))),
		zap.String("encoded_data", encoded),
	)

	// Trim any surrounding quotes from the encoded string
	encoded = strings.Trim(encoded, "\"")

	// Decode the base64 key
	decodedKey, err := base64.StdEncoding.DecodeString(string(key))
	if err != nil {
		logger.Error("Failed to decode key", zap.Error(err))
		return "", err
	}

	logger.Debug("Attempting to decode ciphertext")
	ciphertext, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		logger.Error("Failed to decode ciphertext",
			zap.Error(err),
			zap.String("encoded", encoded),
		)
		return "", err
	}

	block, err := aes.NewCipher(decodedKey)
	if err != nil {
		logger.Error("Failed to create cipher", zap.Error(err))
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		logger.Error("Failed to create GCM", zap.Error(err))
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		logger.Error("Ciphertext too short",
			zap.Int("ciphertext_length", len(ciphertext)),
			zap.Int("nonce_size", nonceSize),
		)
		return "", err
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		logger.Error("Failed to decrypt", zap.Error(err))
		return "", err
	}

	return string(plaintext), nil
}
