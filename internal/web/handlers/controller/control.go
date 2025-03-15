package controller

import (
	"GoStreamRecord/internal/bot"
	"GoStreamRecord/internal/web/handlers/connection"
	"GoStreamRecord/internal/web/handlers/cookies"
	web_status "GoStreamRecord/internal/web/handlers/status"
	"encoding/json"
	"fmt"
	"net/http"
)

var ControllerNotifier = connection.NewNotifier()

// dcodes a JSON payload with a "command" field (start, stop, or restart)
func ControlHandler(w http.ResponseWriter, r *http.Request) {
	if !cookies.Session.IsLoggedIn(w, r) {
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

	go bot.Bot.Command(reqData.Command, reqData.Name)
	resp := web_status.Response{
		Message: fmt.Sprintf("Exected command '%s'", reqData.Command),
		Status:  "success",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)

}
