package streamers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"remoteCtrl/internal/db"
	"remoteCtrl/internal/web/handlers/status"
)

// Handles POST /api/upload.
// It reads an uploaded file and returns a dummy success response.
func UploadHandler(w http.ResponseWriter, r *http.Request) {
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

	if filepath.Ext(handler.Filename) != ".json" {
		return
	}

	// For demonstration, we'll read the file's contents (but not store it)
	fileContent, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Error reading file", http.StatusInternalServerError)
		return
	}
	counter := 0

	resp := status.Response{}

	var import_list []db.Streamer
	err = json.Unmarshal(fileContent, &import_list)
	if err != nil {
		resp = status.Response{
			Status:  status.Status,
			Message: fmt.Sprintf("Failed to import new streamers!"),
		}
		resp.Status.Ok = false
		http.Error(w, "Error reading file", http.StatusInternalServerError)
		return
	}

	shareStr := r.URL.Query().Get("share")
	share := false
	if shareStr == "true" {
		share = true
	}
	for _, streamer := range import_list {
		db.DataBase.NewStreamer(streamer.Name, streamer.Provider, db.DataBase.RequestUserID(r), share)
		counter++
	}

	resp = status.Response{
		Status:  status.Status,
		Message: fmt.Sprintf("Imported %d new streamers!", counter),
	}
	resp.Status.Ok = true
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
