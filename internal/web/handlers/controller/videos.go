package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"remoteCtrl/internal/db"
	"remoteCtrl/internal/system"
	"remoteCtrl/internal/web/handlers/cookie"
	"strconv"
	"strings"
)

type Video struct {
	URL     string `json:"url"`
	Name    string `json:"name"`
	NoFiles string `json:"error"`
}

func GetFiles(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie(cookie.SessionCookieName)
	if err != nil {
		return
	}
	user_id := db.DataBase.Users.NameToID(cookie.UserSessions[c.Value].Name)

	videos, _ := db.DataBase.Videos.ListAvailable(system.System.Context, user_id)
	fmt.Println("Videos for user: ", videos)
	files := []Video{}

	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page < 1 {
		page = 1
	}

	if system.System.Config.OutputFolder == "" {
		system.System.Config.OutputFolder = "videos"
	}
	if _, err := os.Stat(system.System.Config.OutputFolder); os.IsNotExist(err) {
		os.MkdirAll(system.System.Config.OutputFolder, 0755)
	}

	for _, video := range videos {
		files = append(files, Video{URL: "/files/" + filepath.Join(filepath.Base(filepath.Dir(video.Filepath)), video.Name), Name: video.Name})

	}
	err = filepath.Walk(system.System.Config.OutputFolder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && isVideoFile(info.Name()) {

			files = append(files, Video{URL: "/files/" + filepath.Join(filepath.Base(filepath.Dir(path)), info.Name()), Name: info.Name()})
		}
		return nil
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if len(files) == 0 {
		files = append(files, Video{URL: "", Name: "", NoFiles: fmt.Sprintf("No files available. Try adding some to '%s'", system.System.Config.OutputFolder)})

	}

	start := (page - 1) * 999
	end := start + 999
	if start >= len(files) {
		start = len(files)
	}
	if end > len(files) {
		end = len(files)
	}

	paginatedFiles := files[start:end]

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(paginatedFiles)
}

// isVideoFile returns true if the file extension indicates a video file.
func isVideoFile(filename string) bool {
	extensions := []string{".mp4", ".avi", ".mov", ".mkv"}
	lower := strings.ToLower(filename)
	for _, ext := range extensions {
		if strings.HasSuffix(lower, ext) {
			return true
		}
	}
	return false
}

// DeleteFilesRequest represents the expected JSON payload.
type DeleteFilesRequest struct {
	Files []string `json:"files"`
}

// DeleteFilesResponse is the structure of our JSON response.
type DeleteFilesResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}

// DeleteFilesHandler handles requests to delete files.
func DeleteFiles(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse the JSON body.
	var req DeleteFilesRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad Request: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Validate that files were provided.
	if len(req.Files) == 0 {
		resp := DeleteFilesResponse{
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
	for _, video := range req.Files {
		video_path := filepath.Join(system.System.Config.OutputFolder, strings.Replace(video, "/videos/", "", 1))
		log.Println("Deleting video:", video_path)
		err := os.Remove(video_path)
		if err != nil {
			video_erros++
			log.Println("error deleting video: ", err)
			continue
		}
	}

	resp := DeleteFilesResponse{
		Success: video_erros == 0,
		Message: fmt.Sprintf("Deleted %d videos", len(req.Files)-video_erros),
	}

	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(resp)

}
