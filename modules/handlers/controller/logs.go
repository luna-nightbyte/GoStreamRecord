package controller

import (
	"GoStreamRecord/modules/file"
	"GoStreamRecord/modules/handlers/cookies"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

// handleLogs returns log lines from "logs.txt" via GET /api/logs.
func HandleLogs(w http.ResponseWriter, r *http.Request) {
	if !cookies.Session.IsLoggedIn(w, r) {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	if _, err := os.Stat(file.Log_path); os.IsNotExist(err) {
		// If no log file exists, return an empty array.
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]string{})
		return
	}
	content, err := ioutil.ReadFile(file.Log_path)
	if err != nil {
		http.Error(w, "Error reading log file", http.StatusInternalServerError)
		return
	}

	// Split the file into lines and filter empty lines.
	lines := strings.Split(string(content), "\n")
	var filtered []string
	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			filtered = append(filtered, line)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(filtered)
}
