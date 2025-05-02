package jwt

import "golang.org/x/crypto/bcrypt"

func VerifyPassword(hashedPassword string, candidatePassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(candidatePassword))
	return err == nil
}

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}
