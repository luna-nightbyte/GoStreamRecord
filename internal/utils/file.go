package utils

import (
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
