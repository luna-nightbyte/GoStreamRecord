package status

import (
	"GoStreamRecord/internal/bot"
	"GoStreamRecord/internal/bot/recorder"
	"GoStreamRecord/internal/db"
	"GoStreamRecord/internal/file"
	"GoStreamRecord/internal/handlers/cookies"
	"encoding/json"
	"net/http"
)

// Response is a generic response structure for our API endpoints.
type Response struct {
	Status    string              `json:"status"`
	Message   string              `json:"message,omitempty"`
	Data      interface{}         `json:"data,omitempty"`
	BotStatus []recorder.Recorder `json:"botStatus"`
}

func StatusHandler(w http.ResponseWriter, r *http.Request) {
	if !cookies.Session.IsLoggedIn(w, r) {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET allowed", http.StatusMethodNotAllowed)
		return
	}

	// Reload streamer list from config file
	file.ReadJson("streamers", file.Streamers_json, &db.Config.Streamers)

	bot.Bot.StopRunningEmpty()
	// Fetch current recording status
	recorderStatus := bot.Bot.ListRecorders()
	isRecording := false
	for _, s := range recorderStatus {
		if s.IsRecording {
			isRecording = true
			break
		}
	}

	// Prepare response
	recorder := Response{
		BotStatus: recorderStatus,
		Status:    "Stopped",
	}

	if isRecording {
		recorder.Status = "Running"
	}

	// Send JSON response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(recorder); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func ResponseHandler(w http.ResponseWriter, r *http.Request, message string, data interface{}) {
	resp := Response{
		Status:  "success",
		Message: message,
		Data:    data,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
