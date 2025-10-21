package utils

import (
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

func RemoveAll(path string) error {
	return os.RemoveAll(path)

}

// Returns "data" []string
func ReadAll(pwd string, data []string) ([]string, error) {
	folder, err := os.ReadDir(pwd)
	if err != nil {
		log.Println(err)
		return data, err
	}
	for _, file := range folder {
		if file.IsDir() {
			data, _ = ReadAll(filepath.Join(pwd, file.Name()), data)
			continue
		}
		tmp := filepath.Join(pwd, file.Name())
		data = append(data, tmp)
	}
	return data, err
}

func CopyFile(src, dst string) error {
	// Open source file
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	// Create destination file
	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	// Copy contents from source to destination
	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return err
	}

	// Flush to disk
	err = destFile.Sync()
	if err != nil {
		return err
	}

	return nil
}

func FileSHA256(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	hasher := sha256.New()
	if _, err := io.Copy(hasher, file); err != nil {
		return "", fmt.Errorf("failed to copy file contents to hasher: %w", err)
	}
	hashInBytes := hasher.Sum(nil)

	// 6. Convert the byte slice to a hexadecimal string.
	hashString := fmt.Sprintf("%x", hashInBytes)

	return hashString, nil
}
