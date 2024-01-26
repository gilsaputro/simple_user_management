package token

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// TokenConfig is list dependencies of token package
type TokenConfig struct {
	Secret        string
	ExpTimeInHour int64
}

// TokenMethod is method for Token Package
type TokenMethod interface {
	GenerateToken(TokenBody) (string, error)
	ValidateToken(string) (TokenBody, error)
}

// TokenBody is list parameter that will be stored as token
type TokenBody struct {
	UserID int
}

// NewTokenMethod is func to generate TokenMethod interface
func NewTokenMethod(secret string, expinHour int64) TokenMethod {
	return TokenConfig{
		Secret:        secret,
		ExpTimeInHour: expinHour,
	}
}

// GenerateToken is func to generate token from body
func (t TokenConfig) GenerateToken(body TokenBody) (string, error) {
	claims := jwt.MapClaims{
		"id":  body.UserID,
		"exp": time.Now().Add(time.Hour * time.Duration(t.ExpTimeInHour)).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(t.Secret))
}

// ValidateToken is func to validate and generate body from token
func (t TokenConfig) ValidateToken(tokenString string) (TokenBody, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(t.Secret), nil
	})

	if err != nil {
		return TokenBody{}, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userIDFloat64, ok := claims["id"].(float64)
		if !ok {
			return TokenBody{}, fmt.Errorf("Invalid Token")
		}

		if userIDFloat64 > 0 {
			return TokenBody{UserID: int(userIDFloat64)}, nil
		}
	}
	return TokenBody{}, fmt.Errorf("Invalid Token")
}
