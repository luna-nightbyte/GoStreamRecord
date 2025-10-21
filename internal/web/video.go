package web

import (
	"encoding/json"
	"net/http"
	"net/url"
	"path/filepath"
	"remoteCtrl/internal/db"
	"remoteCtrl/internal/media/localfolder"
	"remoteCtrl/internal/system"
	"remoteCtrl/internal/web/handlers/login"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

func getVideos(api, baseDir string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pageNum, _ := strconv.Atoi(r.FormValue("page"))
		pageSize, _ := strconv.Atoi(r.FormValue("pageSize"))
		videos, _ := db.DataBase.ListVisibleVideosForUser(system.System.Context, db.GetUserID(r))
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
			outVideos = append(outVideos, localfolder.Video{
				URL:  "/videos/" + encoded,
				Name: rel,
			})

		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(outVideos)

	}
}

func VideoMux(api string, r *mux.Router, baseDir string) {

	r.HandleFunc(api, login.RequireAuth(getVideos(api, baseDir)))
}
