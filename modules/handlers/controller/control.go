package controller

import (
	"GoStreamRecord/modules/bot"
	"GoStreamRecord/modules/handlers/cookies"
	web_status "GoStreamRecord/modules/handlers/status"
	"encoding/json"
	"fmt"
	"net/http"
)

// dcodes a JSON payload with a "command" field (start, stop, or restart)
// and returns a dummy response.
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
