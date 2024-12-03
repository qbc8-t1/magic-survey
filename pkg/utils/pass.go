package utils

import (
	"golang.org/x/crypto/bcrypt"
)

// HashPassword hashes the given password
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), err
}

// CheckPasswordHash compare a hashed password with the plain-text
func CheckPasswordHash(password, hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))

	if err != nil {
		return err
	}
	return nil
}
