package utils

import (
	"os"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword hashes a plain text password using bcrypt
func HashPassword(password string) (string, error) {
	cost := bcrypt.DefaultCost
	if os.Getenv("APP_ENV") == "test" || os.Getenv("GO_ENV") == "test" || os.Getenv("TESTING") == "true" {
		cost = bcrypt.MinCost
	}
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	return string(bytes), err
}

// CheckPasswordHash compares a hashed password with a plain text candidate
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
