package status

import (
	"encoding/json"
	"net/http"
	"remoteCtrl/internal/db"
	"remoteCtrl/internal/system"
	"remoteCtrl/internal/system/cookies"
	"remoteCtrl/internal/system/settings"
)

// Response is a generic response structure for our API endpoints.
type Response struct {
	Status  status      `json:"status"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	// BotStatus []recorder.Recorder `json:"botStatus,omitempty"`
}

var Status status

type status struct {
	Ok              bool `json:"ok"`
	IsOnline        bool `json:"is_online"`
	Is_Recording    bool `json:"is_recording"`
	Is_Downloading  bool `json:"is_downloading"`
	Is_Fixing_Codec bool `json:"is_fixing_codec"`
}

func StatusHandler(w http.ResponseWriter, r *http.Request) {
	if !cookies.Session.IsLoggedIn(system.System.DB.APIKeys, w, r) {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET allowed", http.StatusMethodNotAllowed)
		return
	}

	// Reload streamer list from config file
	db.Write(settings.CONFIG_STREAMERS_PATH, &system.System.DB.Streamers)

	response := Response{
		Status: Status,
	}
	// Send JSON response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func ResponseHandler(w http.ResponseWriter, r *http.Request, message string, ok bool, data interface{}) {
	resp := Response{
		Message: message,
		Data:    data,
	}
	resp.Status = Status
	resp.Status.Ok = ok
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
