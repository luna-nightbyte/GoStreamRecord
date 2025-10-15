// internal/web/handlers/video.go
package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"remoteCtrl/internal/utils"
	"strings"

	"github.com/gorilla/mux"
)

type Video struct {
	URL      string `json:"url"`  // URL-encoded, ready to use in <video src>
	Name     string `json:"name"` // display name (relative path with slashes)
	NoVideos string `json:"error"`
}

func getVideos(baseDir string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var videos []Video

		// Walk the SAME directory that FileServer serves:
		err := filepath.WalkDir(baseDir, func(fp string, d os.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() {
				return nil
			}

			// relative path (OS-agnostic → forward slashes for URLs)
			rel, err := filepath.Rel(baseDir, fp)
			if err != nil {
				return nil
			}
			rel = filepath.ToSlash(rel) // e.g. "sub/My File.mp4"
			fmt.Println("adding")
			utils.VideoVerify.Add(filepath.Join(baseDir, rel))
			fmt.Println("added")

			// Build a URL-encoded path segment-by-segment
			segs := strings.Split(rel, "/")
			for i, s := range segs {
				segs[i] = url.PathEscape(s) // encodes spaces, #, +, etc.
			}
			encoded := strings.Join(segs, "/")
			videos = append(videos, Video{
				URL:  "/videos/" + encoded, // maps 1:1 to FileServer
				Name: rel,                  // nice for display
			})
			return nil
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if len(videos) == 0 {
			// Show the absolute dir to help debugging
			abs, _ := filepath.Abs(baseDir)
			videos = append(videos, Video{NoVideos: "No videos available. Add files to: " + abs})
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(videos)

	}
}

// Wire it up with the same baseDir as FileServer:
func VideoMux(api string, r *mux.Router, baseDir string) {
	r.HandleFunc(api, getVideos(baseDir))
}
