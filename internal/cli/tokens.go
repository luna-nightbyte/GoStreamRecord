package cli

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strconv"

	"github.com/gorilla/securecookie"
)

func generateCookieToken(args []string) (string, error) {
	lenght := 32
	if len(args) > 2 {
		tmpLength, err := strconv.Atoi(args[2])
		if err == nil {
			lenght = tmpLength
		} else {
			return string(securecookie.GenerateRandomKey(lenght)), fmt.Errorf("not a valid int. Defaulting to 32.")
		}
	}
	return string(securecookie.GenerateRandomKey(lenght)), nil
}
func generateSessionKey(args []string) (string, error) {
	lenght := 32
	if len(args) > 2 {
		tmpLength, err := strconv.Atoi(args[2])
		if err == nil {
			lenght = tmpLength
		}
	}
	bytes := make([]byte, lenght)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
