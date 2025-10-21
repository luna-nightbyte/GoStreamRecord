package utils

import (
	"math/rand"
	"strings"
	"time"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

// RandString generates a random string of length n.
func RandString(n int) string {
	rand.Seed(time.Now().UnixNano())
	var sb strings.Builder
	sb.Grow(n) // Pre-allocate memory for efficiency

	k := len(letterRunes)
	for i := 0; i < n; i++ {
		// Pick a random rune from the letterRunes slice
		sb.WriteRune(letterRunes[rand.Intn(k)])
	}

	return sb.String()
}
