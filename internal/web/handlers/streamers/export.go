package streamers

import (
	"encoding/json"
	"net/http"
	"remoteCtrl/internal/db"
)

// Handles GET /api/download.
// It sends a dummy file for download.
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	// if !cookies.Session.IsLoggedIn(system.System.Config.APIKeys, w, r) {
	// 	http.Redirect(w, r, "/login", http.StatusFound)
	// 	return
	// }
	http.Redirect(w, r, "/home", http.StatusFound)
}
func DownloadHandler(w http.ResponseWriter, r *http.Request) {
	// if !cookies.Session.IsLoggedIn(system.System.Config.APIKeys, w, r) {
	// 	http.Redirect(w, r, "/login", http.StatusFound)
	// 	return
	// }
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET allowed", http.StatusMethodNotAllowed)
		return
	}
	streamers, _ := db.DataBase.ListStreamers()
	fileContent, _ := json.Marshal(streamers)
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", "attachment; filename=export.json")
	w.Write([]byte(fileContent))
}
