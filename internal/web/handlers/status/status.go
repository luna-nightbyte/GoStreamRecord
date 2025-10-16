package status

import (
	"encoding/json"
	"net/http"
	"remoteCtrl/internal/db"
	"remoteCtrl/internal/media/stream_recorder"
	"remoteCtrl/internal/media/stream_recorder/recorder"
	"remoteCtrl/internal/system"
	"remoteCtrl/internal/system/cookies"
	"remoteCtrl/internal/system/settings"
)

// Response is a generic response structure for our API endpoints.
type Response struct {
	Status    string              `json:"status"`
	Message   string              `json:"message,omitempty"`
	Data      interface{}         `json:"data,omitempty"`
	BotStatus []recorder.Recorder `json:"botStatus"`
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

	stream_recorder.Streamer.StopRunningEmpty()
	// Fetch current recording status
	recorderStatus := stream_recorder.Streamer.ListRecorders()
	isRecording := false
	for _, s := range recorderStatus {
		if s.IsRecording {
			isRecording = true
			break
		}
	}
	var recorders []recorder.Recorder
	for _, recorder := range recorderStatus {
		recorders = append(recorders, *recorder)
	}
	// Prepare response
	recorder := Response{
		BotStatus: recorders,
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
