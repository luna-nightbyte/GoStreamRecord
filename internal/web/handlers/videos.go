// internal/web/handlers/video.go
package handlers

import (
	"encoding/json"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"remoteCtrl/internal/db"
	"remoteCtrl/internal/system"
	"remoteCtrl/internal/utils"
	"remoteCtrl/internal/web/handlers/cookie"
	"remoteCtrl/internal/web/handlers/login"
	"strings"

	"github.com/gorilla/mux"
)

type Video struct {
	URL      string `json:"url"`
	Name     string `json:"name"`
	NoVideos string `json:"error"`
}

func getVideos(baseDir string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var videos []Video
		c, err := r.Cookie(cookie.SessionCookieName)
		if err != nil {
			return
		}
		username := cookie.UserSessions[c.Value]

		err = filepath.WalkDir(baseDir, func(fp string, d os.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() {
				return nil
			}
			hasAccess, _ := db.DataBase.UserHasAccessToVideo(system.System.Context, username.Name, filepath.Base(fp))
			isAdmin, _ := db.DataBase.IsAdmin(username.Name)
			if hasAccess || isAdmin {
				rel, err := filepath.Rel(baseDir, fp)
				if err != nil {
					return nil
				}
				rel = filepath.ToSlash(rel)
				utils.VideoVerify.Add(filepath.Join(baseDir, rel))

				segs := strings.Split(rel, "/")
				for i, s := range segs {
					segs[i] = url.PathEscape(s)
				}
				encoded := strings.Join(segs, "/")
				videos = append(videos, Video{
					URL:  "/videos/" + encoded,
					Name: rel,
				})
			}
			return nil
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if len(videos) == 0 {
			abs, _ := filepath.Abs(baseDir)
			videos = append(videos, Video{NoVideos: "No videos available. Add files to: " + abs})
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(videos)

	}
}

func VideoMux(api string, r *mux.Router, baseDir string) {
	r.HandleFunc(api, login.RequireAuth(getVideos(baseDir)))
}
