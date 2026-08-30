package main

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// ONLY FOR DEVELOPMENT JWT SECRET KEY
var jwtSecret = []byte("RESPECT")

func generateJWT(email string, ttl time.Duration, jwtSecret []byte) (string, error) {
	expirationTime := time.Now().Add(ttl)

	claims := &jwt.RegisteredClaims{
		Subject:   email,
		ExpiresAt: jwt.NewNumericDate(expirationTime),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
