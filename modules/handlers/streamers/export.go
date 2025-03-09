package streamers

import (
	"GoStreamRecord/modules/db"
	"GoStreamRecord/modules/handlers/cookies"
	"net/http"
)

// Handles GET /api/download.
// It sends a dummy file for download.
func DownloadHandler(w http.ResponseWriter, r *http.Request) {
	if !cookies.Session.IsLoggedIn(w, r) {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET allowed", http.StatusMethodNotAllowed)
		return
	}

	// Dummy file content
	fileContent := ""

	for _, s := range db.Config.Streamers.Streamers {
		fileContent = fileContent + s.Name + "\n"
	}
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", "attachment; filename=export.txt")
	w.Write([]byte(fileContent))
}
