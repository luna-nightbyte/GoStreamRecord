package controller

import (
	"GoStreamRecord/modules/db"
	"GoStreamRecord/modules/file"
	"GoStreamRecord/modules/handlers/cookies"
	"GoStreamRecord/modules/handlers/status"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type Video struct {
	URL      string `json:"url"`
	Name     string `json:"name"`
	NoVideos string `json:"error"`
}

func GetVideos(w http.ResponseWriter, r *http.Request) {
	if !cookies.Session.IsLoggedIn(w, r) {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	videos := []Video{}

	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page < 1 {
		page = 1
	}

	err = filepath.Walk(db.Config.Settings.App.Videos_folder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && file.IsVideoFile(info.Name()) {

			videos = append(videos, Video{URL: "/videos/" + filepath.Join(filepath.Base(filepath.Dir(path)), info.Name()), Name: info.Name()})
		}
		return nil
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if len(videos) == 0 {
		videos = append(videos, Video{URL: "", Name: "", NoVideos: fmt.Sprintf("No videos available. Try adding some to '%s'", db.Config.Settings.App.Videos_folder)})

	}

	start := (page - 1) * 999
	end := start + 999
	if start >= len(videos) {
		start = len(videos)
	}
	if end > len(videos) {
		end = len(videos)
	}

	paginatedVideos := videos[start:end]

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(paginatedVideos)
}

// DeleteVideosRequest represents the expected JSON payload.
type DeleteVideosRequest struct {
	Videos []string `json:"videos"`
}

// DeleteVideosResponse is the structure of our JSON response.
type DeleteVideosResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}

// DeleteVideosHandler handles requests to delete videos.
func DeleteVideos(w http.ResponseWriter, r *http.Request) {
	if !cookies.Session.IsLoggedIn(w, r) {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	// Only allow POST requests.
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse the JSON body.
	var req DeleteVideosRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad Request: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Validate that videos were provided.
	if len(req.Videos) == 0 {
		resp := DeleteVideosResponse{
			Success: false,
			Message: "No videos provided",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
		return
	}

	// Process deletion of each video.
	video_erros := 0
	for _, video := range req.Videos {
		video_path := filepath.Join(db.Config.Settings.App.Videos_folder, strings.Replace(video, "/videos/", "", 1))
		fmt.Println("Deleting video:", video_path)
		err := os.Remove(video_path)
		if err != nil {
			video_erros++
			log.Println("error deleting video: ", err)
			status.ResponseHandler(w, r, "Error deleting video"+video, nil)
			continue
		}
	}

	resp := DeleteVideosResponse{
		Success: video_erros == 0,
		Message: fmt.Sprintf("Deleted %d videos", len(req.Videos)-video_erros),
	}
	status.ResponseHandler(w, r, "Videos deleted", resp)

}
