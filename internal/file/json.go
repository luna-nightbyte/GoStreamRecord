package file

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"
)

// ReadJSON reads and unmarshals JSON data from a file into the provided object.
func ReadJson(folder, filename string, v interface{}) error {
	file, err := os.Open(filepath.Join(dbDir, folder, filename))
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(v); err != nil {
		fmt.Println(err)
		return fmt.Errorf("failed to decode JSON: %w", err)
	}
	return nil
}

// WriteJSON marshals the given object and writes it to a file.
func WriteJson(folder, filename string, v interface{}) error {
	file, err := os.Create(filepath.Join(dbDir, folder, filename))
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ") // Pretty print
	if err := encoder.Encode(v); err != nil {
		return fmt.Errorf("failed to encode JSON: %w", err)
	}
	return nil
}

// CheckJson compares the JSON content of a file with the expected data.
// If they don't match, it returns an error explaining the difference.
func CheckJson(folder, filename string, expected any) bool {
	var actual map[string]interface{}
	if err := ReadJson(folder, filename, &actual); err != nil {
		log.Printf("failed to read JSON file: %v\n", err)
		return false
	}
	expectedBytes, err := json.Marshal(expected)
	if err != nil {
		log.Printf("failed to marshal expected data: %v\n", err)
		return false
	}
	var expectedMap map[string]interface{}
	if err := json.Unmarshal(expectedBytes, &expectedMap); err != nil {
		log.Printf("failed to unmarshal expected data: %v", err)
		return false
	}

	return compareJSON(actual, expectedMap)
}

// compareJSON compares two JSON objects and returns an error if they don't match.
func compareJSON(actual, expected map[string]interface{}) bool {
	ok := true
	for key, expectedValue := range expected {
		actualValue, exists := actual[key]
		if !exists {
			ok = false
			fmt.Printf("missing key: {\"%s\":}\n", key)
			log.Printf("missing key: {\"%s\":}\n", key)
		}
		if !compareValues(actualValue, expectedValue) {
			ok = false
			fmt.Printf("Wrong value! \nGot: {\"%s\":%s}, expected: {\"%s\":%s}\n", key, expectedValue, key, actualValue)
			log.Printf("Wrong value! \nGot: {\"%s\":%s}, expected: {\"%s\":%s}\n", key, expectedValue, key, actualValue)
		}
	}

	for key := range actual {
		if _, exists := expected[key]; !exists {
			ok = false
			fmt.Printf("unexpected key: {\"%s\":} in json\n", key)
			log.Printf("unexpected key: {\"%s\":} in json\n", key)
		}
	}

	return ok
}

// compareValues performs a deep comparison of two values.
func compareValues(actual, expected interface{}) bool {
	return reflect.TypeOf(actual) == reflect.TypeOf(expected)
}
