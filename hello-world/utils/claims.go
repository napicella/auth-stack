package utils

import "github.com/dgrijalva/jwt-go"

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}
