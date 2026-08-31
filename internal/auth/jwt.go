package auth

import (
	"time"

	"github.com/fistos3rr/ideagen/internal/data"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JwtClaims struct {
	UserID int64  `json:"user_id"`
	Role   string `json:"role"`
	Type   string `json:"type"`
	JTI    string `json:"jti"`
	jwt.RegisteredClaims
}

func generateJWT(user *data.User, tokenType string, ttl time.Duration, jwtSecret []byte, jti string) (string, error) {
	claims := JwtClaims{
		UserID: user.ID,
		Role:   user.Role,
		Type:   tokenType,
		JTI:    jti,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func GenerateJWTAccessToken(user *data.User, ttl time.Duration, jwtSecret []byte) (string, error) {
	return generateJWT(user, "access", ttl, jwtSecret, "")
}

func GenerateJWTRefreshToken(user *data.User, ttl time.Duration, jwtSecret []byte, jti string) (string, error) {
	return generateJWT(user, "refresh", ttl, jwtSecret, jti)
}

func parseAndValidate(tokenString string, jwtSecret []byte) (*JwtClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JwtClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*JwtClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, jwt.ErrSignatureInvalid
}

func ParseAccessToken(tokenString string, jwtSecret []byte) (*JwtClaims, error) {
	claims, err := parseAndValidate(tokenString, jwtSecret)
	if err != nil {
		return nil, err
	}
	if claims.Type != "access" {
		return nil, jwt.ErrInvalidType
	}
	return claims, nil
}

func ParseRefreshToken(tokenString string, jwtSecret []byte) (*JwtClaims, error) {
	claims, err := parseAndValidate(tokenString, jwtSecret)
	if err != nil {
		return nil, err
	}
	if claims.Type != "refresh" {
		return nil, jwt.ErrInvalidType
	}
	return claims, nil
}

func NewTokenID() string {
	return uuid.New().String()
}
