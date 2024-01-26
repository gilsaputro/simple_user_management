package hash

import (
	"golang.org/x/crypto/bcrypt"
)

// HashConfig is list dependencies of Hash Package
type HashConfig struct {
	cost int
}

// HashMethod is list method for Hash Package
type HashMethod interface {
	HashValue(string) ([]byte, error)
	CompareValue(string, string) bool
}

// NewHashMethod func to create HashMethod interface
func NewHashMethod(cost int) HashMethod {
	return &HashConfig{
		cost: cost,
	}
}

// HashValue func to hash value
func (h *HashConfig) HashValue(value string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(value), h.cost)
}

// CompareValue func to hashed value with password
func (h *HashConfig) CompareValue(hash string, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}
