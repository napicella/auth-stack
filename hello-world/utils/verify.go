package utils

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
)

func VerifyToken(tknStr, jwtKey string) (*Claims, error){
	// Initialize a new instance of `Claims`
	claims := &Claims{}
	// Parse the JWT string and store the result in `claims`.
	// Note that we are passing the key in this method as well. This method will return an error
	// if the token is invalid (if it has expired according to the expiry time we set on sign in),
	// or if the signature does not match
	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		// the api lies, the key needs to be an array of bytes.
		// The interface{} return type is misleading.
		return []byte(jwtKey), nil
	})
	if err != nil {
		return nil, err
	}
	if !tkn.Valid {
		return nil, errors.New("token is invalid")
	}
	return claims, nil
}