package file

import "strings"

var (
	dbDir                string = "./internal/app/db"
	Settings_json        string = "settings.json"
	Streamers_json       string = "streamers.json"
	API_json             string = "api_keys.json"
	Users_json           string = "users.json"
	YoutubeDL_configPath string = "youtube-dl.config"
	Index_path           string = "./internal/app/web/index.html"
	Login_path           string = "./internal/app/web/login.html"
	Videos_folder        string = "./output/videos"
	Log_path             string = "./output/GoStreamRecord.log"
)

// isVideoFile returns true if the file extension indicates a video file.
func IsVideoFile(filename string) bool {
	extensions := []string{".mp4", ".avi", ".mov", ".mkv"}
	lower := strings.ToLower(filename)
	for _, ext := range extensions {
		if strings.HasSuffix(lower, ext) {
			return true
		}
	}
	return false
}
