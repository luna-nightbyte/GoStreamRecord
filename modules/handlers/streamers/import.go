package streamers

import (
	"GoStreamRecord/modules/db"
	"GoStreamRecord/modules/handlers/cookies"
	web_status "GoStreamRecord/modules/handlers/status"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strings"
)

// Handles POST /api/upload.
// It reads an uploaded file and returns a dummy success response.
func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if !cookies.Session.IsLoggedIn(w, r) {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	// Limit the size of the incoming request to 10MB
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	// Retrieve file from posted form-data
	file, handler, err := r.FormFile("file")

	if err != nil {
		http.Error(w, "Error retrieving the file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	if filepath.Ext(handler.Filename) != ".txt" {
		return
	}

	// For demonstration, we'll read the file's contents (but not store it)
	fileContent, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Error reading file", http.StatusInternalServerError)
		return
	}
	counter := 0
	for _, line := range strings.Split(string(fileContent), "\n") {
		if db.Config.Streamers.Exist(line) {
			continue
		}
		counter++
		db.Config.AddStreamer(line, r.FormValue("provider"))
	}
	resp := web_status.Response{
		Status:  "success",
		Message: fmt.Sprintf("Added %d new streamers!", counter),
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
