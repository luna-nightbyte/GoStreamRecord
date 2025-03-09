package dbapi

import (
	"GoStreamRecord/internal/file"
	"crypto/rand"
	"encoding/hex"

	"golang.org/x/crypto/bcrypt"
)

type API_secrets struct {
	Keys []ApiKeys `json:keys`
}

type ApiKeys struct {
	User string `json:user`
	Name string `json:name`
	Key  string `json:secret`
}

func (a *API_secrets) Load() error {
	return file.ReadJson("api", file.API_json, &a)
}

func (a API_secrets) NewKey() ApiKeys {
	return ApiKeys{}
}

// GenerateAPIKey creates a secure random API key
func (a *ApiKeys) GenerateAPIKey(length int) (string, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func (a *ApiKeys) HashAPIKey(apiKey string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(apiKey), bcrypt.DefaultCost)
	return string(hash), err
}

func VerifyAPIKey(hashedKey, apiKey string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashedKey), []byte(apiKey)) == nil
}
