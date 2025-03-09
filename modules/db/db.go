package db

import (
	dbapi "GoStreamRecord/modules/db/api"
	dblogin "GoStreamRecord/modules/db/login"
	"GoStreamRecord/modules/db/settings"
	"GoStreamRecord/modules/db/streamers"
	"GoStreamRecord/modules/file"
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type configs struct {
	APIKeys   dbapi.API_secrets
	Settings  settings.Settings
	Streamers streamers.List
	Users     dblogin.Logins
}

var (
	Config = configs{
		APIKeys:   dbapi.API_secrets{Keys: []dbapi.ApiKeys{}},
		Settings:  settings.Settings{},
		Streamers: streamers.List{Streamers: []streamers.Streamer{}},
		Users:     dblogin.Logins{Users: []dblogin.Login{}},
	}
	Version string = "dev"
)

func init() {
	loadConfigurations()
}
func loadConfigurations() {
	loadConfig(file.API_keys_file, &Config.APIKeys)
	err := loadConfig(file.Config_json_path, &Config.Settings)
	if err != nil {
		log.Fatal(err)
	}
	loadConfig(file.Streamers_json_path, &Config.Streamers)

	err = loadConfig(file.Users_json_path, &Config.Users)
	if err != nil {
		log.Fatal(err)
	}
}
func loadConfig(path string, target any) error {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("Config file %s not found.", path)
		}
		return fmt.Errorf("Error reading file %s: %v", path, err)
	}

	if ok, err := file.CheckJson(path, target); !ok {
		return fmt.Errorf("Invalid JSON format in %s: %v", path, err)
	}

	if err = json.Unmarshal(data, target); err != nil {
		return fmt.Errorf("Failed to parse JSON in %s: %v", path, err)
	}
	return nil
}

// ----------------- Streamers -----------------

func (c *configs) AddStreamer(streamerName, provider string) string {
	c.Streamers.Add(streamerName, provider)
	ok := write(file.Streamers_json_path, c.Streamers)
	if !ok {
		log.Printf("Error adding %s..\n", streamerName)
		return fmt.Sprintf("Error adding %s..\n", streamerName)
	}
	log.Printf("%s has been added", streamerName)
	return ""
}

func (c *configs) RemoveStreamer(streamerName string) string {
	output := c.Streamers.Remove(streamerName)
	if output == "" {
		return ""
	}
	ok := write(file.Streamers_json_path, c.Streamers)
	if !ok {
		log.Printf("Error removing %s..\n", streamerName)
		return fmt.Sprintf("Error removing %s..\n", streamerName)
	}
	log.Printf("%s has been deleted", streamerName)
	return ""
}

// ----------------- Global General -----------------

func (c *configs) Reload(path string, target any) {
	if err := file.ReadJson(path, target); err != nil {
		log.Printf("Error reloading config: %v", err)
	}
}

func (c *configs) Update(path string, newConfig any) {
	var backup any
	if !read(path, &backup) {
		return
	}
	if !write(path, newConfig) || !verify(path, newConfig) {
		write(path, backup)
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

// ----------------- Local General -----------------

func verify(path string, config any) bool {
	ok, err := file.CheckJson(path, config)
	if !ok {
		log.Printf("Config verification failed: %v", err)
		return false
	}
	return true
}

func read(path string, target any) bool {
	if err := file.ReadJson(path, target); err != nil {
		log.Printf("Error reading config: %v", err)
		return false
	}
	return true
}

func write(path string, config any) bool {
	if err := file.WriteJson(path, config); err != nil {
		log.Printf("Error writing config: %v", err)
		return false
	}
	return true
}
