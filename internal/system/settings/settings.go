package settings

import (
	"fmt"
	"remoteCtrl/internal/db/jsondb"
	"remoteCtrl/internal/system/prettyprint"
)

type app struct {
	Port          int        `json:"port"`
	Loop_interval int        `json:"loop_interval_in_minutes"`
	Files_folder  string     `json:"output_folder"`
	RateLimit     rate_limit `json:"rate_limit"`
	//Cookie        string     `json:"cookie,omit_empty"`
}
type rate_limit struct {
	Enable bool `json:"enable"`
	Time   int  `json:"time"`
}

type Settings struct {
	App         app      `json:"app"`
	GoogleDrive gdrive   `json:"google_drive"`
	Telegram    telegram `json:"telegram"`
}

type StreamList struct {
	List []Streamer `json:"streamers"`
}

type Streamer struct {
	Name     string `json:"name"`
	Provider string `json:"provider"`
}

type DB struct {
	APIKeys API_secrets
	//Users     Logins
	Settings  Settings
	Streamers StreamList
}

const (
	CONFIGS_FOLDER        string = "settings"
	CONFIG_SETTINGS_PATH  string = CONFIGS_FOLDER + "/settings.json"
	CONFIG_USERS_PATH     string = CONFIGS_FOLDER + "/users.json"
	CONFIG_API_PATH       string = CONFIGS_FOLDER + "/api.json"
	CONFIG_STREAMERS_PATH string = CONFIGS_FOLDER + "/streamers.json"
)

func Init() DB {
	var sys DB
	jsondb.Load(CONFIG_SETTINGS_PATH, &sys.Settings)
	//db.LoadConfig(CONFIG_USERS_PATH, &sys.Users)
	//load(CONFIG_USERS_PATH, &sys.Users)
	jsondb.Load(CONFIG_API_PATH, &sys.APIKeys)
	jsondb.Load(CONFIG_STREAMERS_PATH, &sys.Streamers)
	return sys
}
func load(path string, data any) {
	err := jsondb.Load(path, &data)
	if err != nil {
		fmt.Println(prettyprint.Yellow("Example settings file generated for:"), path)
		jsondb.GenerateDefault(path, data)
	}
}
