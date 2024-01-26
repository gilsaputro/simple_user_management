package token

import (
	"crypto/rsa"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// TokenConfig is list dependencies of token package
type TokenConfig struct {
	signKey       *rsa.PrivateKey
	verifyKey     *rsa.PublicKey
	expTimeInHour int64
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

type NewTokenConfig struct {
	PrivateKeyLocation string
	PublicKeyLocation  string
	ExpinHour          int64
}

// NewTokenMethod is func to generate TokenMethod interface
func NewTokenMethod(cfg NewTokenConfig) (TokenMethod, error) {
	privateKey, err := os.ReadFile(cfg.PrivateKeyLocation)
	if err != nil {
		return nil, fmt.Errorf("failed read private key file %s authenticator, err: %s", cfg.PrivateKeyLocation, err)
	}

	publicKey, err := os.ReadFile(cfg.PublicKeyLocation)
	if err != nil {
		return nil, fmt.Errorf("failed read public key file %s authenticator, err: %s", cfg.PublicKeyLocation, err)
	}

	return buildAuthenticator(string(privateKey), string(publicKey), cfg.ExpinHour)
}

// the two keys PrivateKey & PublicKey from generate rsa from openssl
// $ openssl genrsa -out demo.rsa 1024 # the 1024 is the size of the key we are generating
// $ openssl rsa -in demo.rsa -pubout > demo.rsa.pub
// PrivateKey private key generate from "openssl genrsa -out app.rsa keysize"
// PublicKey  public key generate from "openssl rsa -in app.rsa -pubout > app.rsa.pub"
func buildAuthenticator(privateKey, publicKey string, expInHour int64) (TokenMethod, error) {
	signKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(strings.Trim(strings.TrimSpace(privateKey), "\n")))
	if err != nil {
		return nil, fmt.Errorf("failed parse private key, err: %s", err)
	}

	verifyKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(publicKey))
	if err != nil {
		return nil, fmt.Errorf("failed parse public key, err: %s", err)
	}

	return TokenConfig{
		signKey:       signKey,
		verifyKey:     verifyKey,
		expTimeInHour: expInHour,
	}, nil
}

// GenerateToken is func to generate token from body
func (t TokenConfig) GenerateToken(body TokenBody) (string, error) {
	claims := jwt.MapClaims{
		"id":  body.UserID,
		"exp": time.Now().Add(time.Hour * time.Duration(t.expTimeInHour)).Unix(),
	}

	jwtClaim := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	tokenString, err := jwtClaim.SignedString(t.signKey)
	if err != nil {
		err = fmt.Errorf("error while signing token!, err: %s", err)
	}

	return "Bearer " + tokenString, nil
}

// ValidateToken is func to validate and generate body from token
func (t TokenConfig) ValidateToken(tokenString string) (TokenBody, error) {
	// check if it is empty
	if tokenString == "" {
		return TokenBody{}, fmt.Errorf("Invalid Token")
	}

	// validate the tokenCookie
	token, err := jwt.ParseWithClaims(tokenString, &jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		if token.Method == jwt.SigningMethodES256 && token.Valid {
			return nil, fmt.Errorf("unexpected signing method: %v, token not valid", token.Header["alg"])
		}
		return t.verifyKey, nil
	})

	if err != nil {
		return TokenBody{}, fmt.Errorf("Invalid Token")
	}

	claims, ok := token.Claims.(*jwt.MapClaims)
	if !ok || !token.Valid {
		return TokenBody{}, fmt.Errorf("Invalid Token")
	}

	userID, ok := (*claims)["id"].(float64)
	if !ok {
		return TokenBody{}, fmt.Errorf("Invalid Token")
	}

	return TokenBody{
		UserID: int(userID),
	}, nil
}
