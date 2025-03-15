package db

import (
	dbapi "GoStreamRecord/internal/db/api"
	"GoStreamRecord/internal/db/settings"
	"GoStreamRecord/internal/db/streamers"
	dbuser "GoStreamRecord/internal/db/users"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"
)

type configs struct {
	APIKeys   dbapi.API_secrets
	Streamers streamers.List
	Users     dbuser.Logins

	Settings settings.Settings // more or less old and deprecated settings. most should be removed
	Recorder settings.YP_DLP
}

var (
	Config = configs{
		APIKeys:   dbapi.API_secrets{Keys: []dbapi.ApiKeys{}},
		Streamers: streamers.List{Streamers: []streamers.Streamer{}},
		Users:     dbuser.Logins{Users: []dbuser.Login{}},

		Settings: settings.Settings{},
		Recorder: settings.YP_DLP{},
	}
	DbDir   string = "./db"
	Version string = "dev"
)

func init() {

	err := loadConfig("settings", "settings.json", &Config.Settings)
	if err != nil {
		log.Fatal(err)
	}

	err = loadConfig("users", "users.json", &Config.Users)
	if err != nil {
		log.Fatal(err)
	}
	Config.Recorder.Read(filepath.Join(DbDir, "settings", "yt-dlp.json"), &Config.Recorder.YT_DLP)
	// Veriify the commands
	err = loadConfig("settings", "yt-dlp.json", &Config.Recorder.YT_DLP)
	if err != nil {
		log.Fatal(err)
	}

	// Dont need to "break" on errors with streamers. They can be added later.
	loadConfig("streamers", "streamers.json", &Config.Streamers)
	loadConfig("api", "api.json", &Config.APIKeys)

}
func loadConfig(folder, filename string, target any) error {

	if ok := CheckJson(folder, filename, target); !ok {
		return fmt.Errorf("Invalid JSON format in %s", filename)
	}
	err := Read(folder, filename, &target)
	if err != nil {
		return fmt.Errorf("Error reading file %s: %v", filename, err)
	}

	return nil
}

func (c *configs) Update(folder, filename string, newConfig any) {
	var backup any
	if Read(folder, filename, &backup) != nil {
		return
	}
	if Write(folder, filename, &newConfig) != nil || !CheckJson(folder, filename, &newConfig) {
		Write(folder, filename, &backup)
	}
}

func (c *configs) GenerateDefault(path string, jsonFile any) {
	f, err := os.Create(path)
	if err != nil {
		log.Fatalf("Failed to create config file: %v", err)
	}
	defer f.Close()
	data, _ := json.MarshalIndent(&jsonFile, "", "  ")
	f.Write(data)
}

// ReadJSON reads and unmarshals JSON data from a file into the provided object.
func Read(folder, filename string, v interface{}) error {
	f, err := os.Open(filepath.Join(DbDir, folder, filename))
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer f.Close()

	decoder := json.NewDecoder(f)
	if err := decoder.Decode(v); err != nil {
		fmt.Println(err)
		return fmt.Errorf("failed to decode JSON: %w", err)
	}
	return nil
}

// WriteJSON marshals the given object and writes it to a
func Write(folder, filename string, v interface{}) error {
	f, err := os.Create(filepath.Join(DbDir, folder, filename))
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

// CheckJson compares the JSON content of a file with the expected data.
// If they don't match, it returns an error explaining the difference.
func CheckJson(folder, filename string, expected any) bool {
	var actual map[string]interface{}
	if err := Read(folder, filename, &actual); err != nil {
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
