package version

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
)

var Version string = "dev"
var Shasum string = "hello"

func CheckFileSHA256(filePath, shaKey string) (bool, error) {
	// Open the file for reading.
	file, err := os.Open(filePath)
	if err != nil {
		return false, fmt.Errorf("could not open file: %w", err)
	}
	defer file.Close() // Ensure the file is closed when the function exits.

	// Create a new SHA-256 hash object.
	hash := sha256.New()

	// Copy the file's content to the hash object.
	// This calculates the hash of the entire file.
	if _, err := io.Copy(hash, file); err != nil {
		return false, fmt.Errorf("could not copy file content to hash: %w", err)
	}

	// Get the final hash sum and encode it as a hex string.
	calculatedKey := hex.EncodeToString(hash.Sum(nil)) 
	// Compare the calculated hash with the provided SHA key.
	return calculatedKey == shaKey, nil
}
