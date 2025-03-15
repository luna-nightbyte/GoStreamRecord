package streamers

import (
	"GoStreamRecord/internal/db"
	"GoStreamRecord/internal/db/streamers"
	"GoStreamRecord/internal/web/handlers/cookies"
	web_status "GoStreamRecord/internal/web/handlers/status"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
)

// Handles POST /api/upload.
func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if !cookies.Session.IsLoggedIn(w, r) {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	// TODO: Update to modify changes as admin in the web-ui
	// Limit the size of the incoming request to 10MB
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		fmt.Println(err)
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	// Retrieve file from posted form-data
	file, handler, err := r.FormFile("file")

	if err != nil {
		fmt.Println(err)
		http.Error(w, "Error retrieving the file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	if filepath.Ext(handler.Filename) != ".json" {
		return
	}

	fileContent, err := io.ReadAll(file)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Error reading file", http.StatusInternalServerError)
		return
	}
	counter := 0

	resp := web_status.Response{}

	var import_list []streamers.Streamer
	err = json.Unmarshal(fileContent, &import_list)
	if err != nil {

		resp = web_status.Response{
			Status:  "failed",
			Message: fmt.Sprintf("Failed to import new streamers!"),
		}
		fmt.Println(err)
		http.Error(w, "Error reading file", http.StatusInternalServerError)
		return
	}
	for _, streamer := range import_list {
		if db.Config.Streamers.Exist(streamer.Name) {
			continue
		}
		counter++
		db.Config.AddStreamer(streamer.Name, streamer.Provider)
	}

	resp = web_status.Response{
		Status:  "success",
		Message: fmt.Sprintf("Imported %d new streamers!", counter),
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
