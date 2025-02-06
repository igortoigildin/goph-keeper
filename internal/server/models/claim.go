package model

import "github.com/golang-jwt/jwt"

type UserClaims struct {
	jwt.StandardClaims
	Login string `db:"email"`
}
