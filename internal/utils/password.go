package utils

import (
	"golang.org/x/crypto/bcrypt"
)

// hashPassword generates a bcrypt hash of the password.
func HashPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

// checkPasswordHash compares a plaintext password with its hash.
func CheckPasswordHash(password string, hash []byte) bool {
	err := bcrypt.CompareHashAndPassword(hash, []byte(password))
	return err == nil
}
