package db

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func LoadConfig(filePath string, target any) error {
	err := read(filePath, &target)
	if err != nil {
		return fmt.Errorf("Error reading file %s: %v", filePath, err)
	}

	return nil
}
func ForceLoadConfig(filePath string, target any) error {

	err := read(filePath, &target)
	if err != nil {
		return fmt.Errorf("Error reading file %s: %v", filePath, err)
	}

	return nil
}

func Update(filePath string, newConfig any) {
	var backup any
	if read(filePath, &backup) != nil {
		return
	}
	if Write(filePath, &newConfig) != nil {
		Write(filePath, &backup)
	}
}

func GenerateDefault(filePath string, jsonFile any) {

	err := os.MkdirAll(filepath.Dir(filePath), 0755)
	if err != nil {
		log.Fatalf("Failed to create folders: %v", err)
	}
	f, err := os.Create(filepath.Join(filePath))
	if err != nil {
		log.Fatalf("Failed to create config file: %v", err)
	}
	defer f.Close()
	data, _ := json.MarshalIndent(&jsonFile, "", "  ")
	f.Write(data)
}

// ReadJSON reads and unmarshals JSON data from a file into the provided object.
func read(filePath string, v interface{}) error {
	f, err := os.Open(filepath.Join(filePath))
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer f.Close()

	decoder := json.NewDecoder(f)
	if err := decoder.Decode(v); err != nil { 
		return fmt.Errorf("failed to decode JSON: %w", err)
	}
	return nil
}

// WriteJSON marshals the given object and writes it to a
func Write(filePath string, v interface{}) error {
	f, err := os.Create(filepath.Join(filePath))
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer f.Close()

	encoder := json.NewEncoder(f)
	encoder.SetIndent("", "  ") // Pretty print
	if err := encoder.Encode(v); err != nil {
		return fmt.Errorf("failed to encode JSON: %w", err)
	}
	return nil
}
