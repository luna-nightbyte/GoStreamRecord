package web

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path/filepath"
	"remoteCtrl/internal/db"
	"remoteCtrl/internal/media/localfolder"
	"remoteCtrl/internal/web/handlers/login"
	"sort"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

func getVideos(api, basbaseDireDir string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pageNum, _ := strconv.Atoi(r.FormValue("page"))
		pageSize, _ := strconv.Atoi(r.FormValue("pageSize"))
		videos, _ := db.DataBase.ListAvailableVideosForUser(db.DataBase.RequestUserID(r))
		var outVideos []localfolder.Video

		for i, video := range videos {
			if i < (pageSize*pageNum)-pageSize {
				continue
			}
			if i >= pageSize*pageNum {
				break
			}
			rel, err := filepath.Rel(filepath.Dir(video.Filepath), video.Filepath)
			if err != nil {
				continue
			}
			segs := strings.Split(rel, "/")
			for i, s := range segs {
				segs[i] = url.PathEscape(s)
			}
			encoded := strings.Join(segs, "/")
			fmt.Println("add", encoded)
			outVideos = append(outVideos, localfolder.Video{
				URL:  "/videos/" + encoded,
				Name: rel,
			})

		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(outVideos)

	}
}

func getVideos2() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pageNum, err := strconv.Atoi(r.FormValue("page"))
		if err != nil || pageNum < 1 {
			pageNum = 1
		}

		pageSize, err := strconv.Atoi(r.FormValue("pageSize"))
		if err != nil || pageSize <= 0 {
			pageSize = 10
		}

		videos, err := db.DataBase.ListAvailableVideosForUser(db.DataBase.RequestUserID(r))
		if err != nil {
			http.Error(w, "Failed to retrieve videos from database.", http.StatusInternalServerError)
			return
		}
		sort.Slice(videos, func(i, j int) bool {
			return videos[i].ID < videos[j].ID
		})

		startIndex := (pageNum - 1) * pageSize
		endIndex := startIndex + pageSize

		if startIndex >= len(videos) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(localfolder.Video{})
			return
		}

		// If the end index goes beyond the list, adjust it to the end of the list.
		if endIndex > len(videos) {
			endIndex = len(videos)
		}

		pagedVideos := videos[startIndex:endIndex]

		outVideos := make([]localfolder.Video, 0, len(pagedVideos))
		for _, video := range pagedVideos {
			rel := filepath.Base(video.Filepath)
			encoded := url.PathEscape(rel)
			outVideos = append(outVideos, localfolder.Video{
				URL:  "/videos/" + encoded,
				Name: rel,
			})
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(outVideos)
	}
}

func VideoMux(api string, r *mux.Router) {

	r.HandleFunc(api, login.RequireAuth(getVideos2()))
}
