package cookie

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"log"

	"github.com/gorilla/securecookie"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/hkdf"
)

func HkdfBytes(key []byte, info string, n int) []byte {
	out := make([]byte, n)
	k := hkdf.New(sha256.New, key, nil, []byte(info))
	if _, err := k.Read(out); err != nil {
		log.Fatalf("hkdf error: %v", err)
	}
	return out
}

func NewAEADKey(master []byte) []byte { return HkdfBytes(master, "data-aead", 32) }

func NewSecret(length int) []byte {

	key := make([]byte, length)

	_, err := rand.Read(key)
	if err != nil {
		return nil
	}

	return []byte(base64.StdEncoding.EncodeToString(key))
}

func HashString(pw string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
}
func CheckHash(hash []byte, pw string) error {
	return bcrypt.CompareHashAndPassword(hash, []byte(pw))
}

func Encrypt(aeadKey []byte, plaintext string) (string, error) {
	if plaintext == "" {
		return "", nil
	}
	block, err := aes.NewCipher(aeadKey)
	if err != nil {
		return "", err
	}
	aead, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, aead.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return "", err
	}
	ct := aead.Seal(nil, nonce, []byte(plaintext), nil)
	return base64.RawStdEncoding.EncodeToString(append(nonce, ct...)), nil
}

func Decrypt(aeadKey []byte, b64 string) (string, error) {
	if b64 == "" {
		return "", nil
	}
	raw, err := base64.RawStdEncoding.DecodeString(b64)
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher(aeadKey)
	if err != nil {
		return "", err
	}
	aead, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	ns := aead.NonceSize()
	if len(raw) < ns {
		return "", errors.New("ciphertext too short")
	}
	pt, err := aead.Open(nil, raw[:ns], raw[ns:], nil)
	if err != nil {
		return "", err
	}
	return string(pt), nil
}

func NewSecureCookie(master []byte) *securecookie.SecureCookie {
	hashKey := HkdfBytes(master, "cookie-sign", 32)
	blockKey := HkdfBytes(master, "cookie-enc", 32)
	sc := securecookie.New(hashKey, blockKey)
	sc.SetSerializer(securecookie.JSONEncoder{})
	sc.MaxAge(60 * 60 * 8)
	return sc
}
