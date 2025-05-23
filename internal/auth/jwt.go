package auth

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

type JWTAuthenticator struct {
	secretKey string
	aud       string
	iss       string
}

// NewJWTAuthenticator creates a new JWTAuthenticator
func NewJWTAuthenticator(secretKey string, aud string, iss string) *JWTAuthenticator {
	return &JWTAuthenticator{
		secretKey: secretKey,
		aud:       aud,
		iss:       iss,
	}
}

// GenerateToken generates a new token for the given claims
func (j *JWTAuthenticator) GenerateToken(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(j.secretKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// ValidateToken validates the token and returns the token if it is valid
func (j *JWTAuthenticator) ValidateToken(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(j.secretKey), nil
	},
		jwt.WithValidMethods([]string{"HS256"}),
		jwt.WithAudience(j.aud),
		jwt.WithIssuer(j.iss),
		jwt.WithExpirationRequired())
}
