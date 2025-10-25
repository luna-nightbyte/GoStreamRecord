package utils

import (
	"encoding/json"
	"fmt"
)

// marshalIntSlice converts a slice of integers to a JSON string for database storage.
// SQLite does not have a native array type, so storing as JSON is a common workaround.
func marshalIntSlice(slice []int) (string, error) {
	bytes, err := json.Marshal(slice)
	if err != nil {
		return "", fmt.Errorf("could not marshal int slice: %w", err)
	}
	return string(bytes), nil
}

// unmarshalIntSlice converts a JSON string from the database back into a slice of integers.
func unmarshalIntSlice(data string) ([]int, error) {
	var slice []int
	if err := json.Unmarshal([]byte(data), &slice); err != nil {
		return nil, fmt.Errorf("could not unmarshal int slice from json '%s': %w", data, err)
	}
	return slice, nil
}
