package login

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

func HashedPassword(password string) []byte {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal("Error generating password hash:", err)
	}
	return hashedPassword
}

func CompareHashAndPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
