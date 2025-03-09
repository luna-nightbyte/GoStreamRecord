package db

import (
	dbapi "GoStreamRecord/internal/db/api"
	"GoStreamRecord/internal/db/settings"
	"GoStreamRecord/internal/db/streamers"
	dbuser "GoStreamRecord/internal/db/users"
	"GoStreamRecord/internal/file"
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type configs struct {
	APIKeys   dbapi.API_secrets
	Settings  settings.Settings
	Streamers streamers.List
	Users     dbuser.Logins
}

var (
	Config = configs{
		APIKeys:   dbapi.API_secrets{Keys: []dbapi.ApiKeys{}},
		Settings:  settings.Settings{},
		Streamers: streamers.List{Streamers: []streamers.Streamer{}},
		Users:     dbuser.Logins{Users: []dbuser.Login{}},
	}
	Version string = "dev"
)

func init() {
	loadConfigurations()
}
func loadConfigurations() {

	err := loadConfig("settings", file.Settings_json, &Config.Settings)
	if err != nil {
		log.Fatal(err)
	}

	err = loadConfig("users", file.Users_json, &Config.Users)
	if err != nil {
		log.Fatal(err)
	}

	// Dont need to "break" on errors with streamers. They can be added later.
	loadConfig("streamers", file.Streamers_json, &Config.Streamers)
	loadConfig("api", file.API_json, &Config.APIKeys)

}
func loadConfig(folder, filename string, target any) error {

	if ok := file.CheckJson(folder, filename, target); !ok {
		return fmt.Errorf("Invalid JSON format in %s", filename)
	}
	err := file.ReadJson(folder, filename, &target)
	if err != nil {
		return fmt.Errorf("Error reading file %s: %v", filename, err)
	}

	return nil
}

// ----------------- Global General -----------------

func (c *configs) Update(folder, filename string, newConfig any) {
	var backup any
	if file.ReadJson(folder, filename, &backup) != nil {
		return
	}
	if file.WriteJson(folder, filename, &newConfig) != nil || !file.CheckJson(folder, filename, &newConfig) {
		file.WriteJson(folder, filename, &backup)
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
