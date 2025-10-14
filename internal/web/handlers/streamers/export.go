package streamers

import (
	"encoding/json"
	"net/http"
	"remoteCtrl/internal/system"
	"remoteCtrl/internal/system/cookies"
)

// Handles GET /api/download.
// It sends a dummy file for download.
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	if !cookies.Session.IsLoggedIn(system.System.DB.APIKeys, w, r) {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	http.Redirect(w, r, "/home", http.StatusFound)
}
func DownloadHandler(w http.ResponseWriter, r *http.Request) {
	if !cookies.Session.IsLoggedIn(system.System.DB.APIKeys, w, r) {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET allowed", http.StatusMethodNotAllowed)
		return
	}

	fileContent, _ := json.Marshal(system.System.DB.Streamers.List)
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", "attachment; filename=export.json")
	w.Write([]byte(fileContent))
}
