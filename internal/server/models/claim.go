package model

import "github.com/golang-jwt/jwt"

type UserClaims struct {
	jwt.StandardClaims
	Email string `json:"email"`
}
