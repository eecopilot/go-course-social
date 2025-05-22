package auth

import (
	"github.com/golang-jwt/jwt/v5"
)

type JWTAuthenticator struct {
	secretKey string
	aud       string
	iss       string
}

func NewJWTAuthenticator(secretKey string, aud string, iss string) *JWTAuthenticator {
	return &JWTAuthenticator{
		secretKey: secretKey,
		aud:       aud,
		iss:       iss,
	}
}

func (j *JWTAuthenticator) GenerateToken(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(j.secretKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (j *JWTAuthenticator) ValidateToken(token string) (*jwt.Claims, error) {
	return nil, nil
}
