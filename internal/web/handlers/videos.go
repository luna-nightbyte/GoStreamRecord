package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gorilla/mux"
)

type Video struct {
	URL      string `json:"url"`
	Name     string `json:"name"`
	NoVideos string `json:"error"`
}

func getVideos(w http.ResponseWriter, r *http.Request) {
	videos := []Video{}
	videoFolder := "./videos"

	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page < 1 {
		page = 1
	}

	err = filepath.Walk(videoFolder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			videos = append(videos, Video{URL: "/videos/" + info.Name(), Name: info.Name()})
		}
		return nil
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	exePath, _ := os.Executable()

	exeDir := filepath.Dir(exePath)
	videosDir := filepath.Join(exeDir, videoFolder)
	if len(videos) == 0 {
		videos = append(videos, Video{URL: "", Name: "", NoVideos: fmt.Sprintf("No videos available. Try adding some to '%s'", videosDir)})

	}
	pageSize := 10
	start := (page - 1) * pageSize
	end := start + pageSize
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

// Define a method on the new type
func VideoMux(api string, r *mux.Router) {
	r.HandleFunc(api, getVideos)
}
