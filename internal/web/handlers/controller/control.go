package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"remoteCtrl/internal/media/stream_recorder"
	"remoteCtrl/internal/utils"
	"remoteCtrl/internal/web/handlers/status"
)

// dcodes a JSON payload with a "command" field (start, stop, or restart, fix-codec)
// and returns a dummy response.
func ControlHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	type RequestData struct {
		Command string `json:"command"`
		Name    string `json:"name"`
	}
	var reqData RequestData
	if err := json.NewDecoder(r.Body).Decode(&reqData); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	resp := status.Response{
		Message: fmt.Sprintf("Executed command '%s'", reqData.Command),
	}
	if reqData.Name == "" { // All was selected or Video codec fix

		// List all active bots/recorders
		recorders := stream_recorder.Streamer.ListRecorders()
		for _, recorder := range recorders {
			go stream_recorder.Streamer.Execute(reqData.Command, recorder.Website.Username)
		}
	} else {
		currentCodecQueue := len(utils.VideoVerify.Queue)
		if currentCodecQueue == 0 {
			status.Status.Ok = false
			status.Status.Is_Fixing_Codec = false
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(resp)
			return
		}
		go stream_recorder.Streamer.Execute(reqData.Command, reqData.Name)

	}
	status.Status.Ok = true
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)

}
