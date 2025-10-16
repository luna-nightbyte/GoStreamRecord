package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"remoteCtrl/internal/media/stream_recorder"
	"remoteCtrl/internal/system"
	"remoteCtrl/internal/system/cookies"
	"remoteCtrl/internal/web/handlers/status"
)

// dcodes a JSON payload with a "command" field (start, stop, or restart)
// and returns a dummy response.
func ControlHandler(w http.ResponseWriter, r *http.Request) {
	if !cookies.Session.IsLoggedIn(system.System.DB.APIKeys, w, r) {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
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
	if reqData.Name == "" { // All was selected
		for _, recorder := range stream_recorder.Streamer.ListRecorders() {
			go stream_recorder.Streamer.Execute(reqData.Command, recorder.Website.Username)
		}
	} else {

		go stream_recorder.Streamer.Execute(reqData.Command, reqData.Name)
	}
	resp := status.Response{
		Message: fmt.Sprintf("Exected command '%s'", reqData.Command),
		Status:  "success",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)

}
