package web

import (
	"bytes"
	"context"
	"embed"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"strings"
	"time"

	"remoteCtrl/internal/db"
	"remoteCtrl/internal/system"
	"remoteCtrl/internal/system/cookies"
	"remoteCtrl/internal/system/prettyprint"
	"remoteCtrl/internal/utils"
	"remoteCtrl/internal/web/handlers"
	webController "remoteCtrl/internal/web/handlers/controller"
	"remoteCtrl/internal/web/handlers/cookie"
	"remoteCtrl/internal/web/handlers/login"
	"remoteCtrl/internal/web/handlers/status"
	"remoteCtrl/internal/web/handlers/streamers"
	"remoteCtrl/internal/web/handlers/users"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

var loginFS fs.FS

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		return
	}

	indexFile, err := loginFS.Open("index.html")
	if err != nil {
		http.Error(w, "Login index file not found", http.StatusInternalServerError)
		return
	}
	defer indexFile.Close()
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	io.Copy(w, indexFile)
}
func ServeHTTP(ctx context.Context, eLogin, eApp embed.FS) {

	var app handlers.API
	app.Router = mux.NewRouter()
	app.Router.HandleFunc("/api/download", login.RequireAuth(handlers.DownloadHandler))
	// app.Router.HandleFunc("/api/progress", video_download.Handler))
	app.Router.HandleFunc("/api/add-streamer", login.RequireAuth(streamers.AddStreamer))
	app.Router.HandleFunc("/api/get-streamers", login.RequireAuth(streamers.GetStreamers))
	app.Router.HandleFunc("/api/remove-streamer", login.RequireAuth(streamers.RemoveStreamer))
	app.Router.HandleFunc("/api/control", login.RequireAuth(webController.ControlHandler))
	app.Router.HandleFunc("/api/get-online-status", login.RequireAuth(streamers.CheckOnlineStatus))
	app.Router.HandleFunc("/api/import", login.RequireAuth(streamers.UploadHandler))
	app.Router.HandleFunc("/api/export", login.RequireAuth(streamers.DownloadHandler))
	app.Router.HandleFunc("/api/status", login.RequireAuth(status.StatusHandler))
	app.Router.HandleFunc("/api/get-videos", login.RequireAuth(webController.GetFiles))
	app.Router.HandleFunc("/api/logs", login.RequireAuth(webController.HandleLogs))
	app.Router.HandleFunc("/api/delete-videos", login.RequireAuth(webController.DeleteFiles))

	app.Router.PathPrefix("/api/generate-api-key").Handler(login.RequireAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("gen api")
		//	cookies.GenAPIKeyHandler(system.System.Config.APIKeys, w, r)
	})))
	app.Router.PathPrefix("/api/delete-api-key").Handler(login.RequireAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("del api")
		//	cookies.DeleteAPIKeyHandler(system.System.Config.APIKeys, w, r)
	})))
	app.Router.PathPrefix("/api/keys").Handler(login.RequireAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("get api")
		//	cookies.GetAPIkeys(system.System.Config.APIKeys, w, r)
	})))
	app.Router.HandleFunc("/api/user_info", login.RequireAuth(users.GetUsers))
	app.Router.HandleFunc("/api/add-user", login.RequireAuth(users.AddUser))
	app.Router.HandleFunc("/api/update-user", login.RequireAuth(users.UpdateUsers))
	app.Router.HandleFunc("/api/health", login.RequireAuth(handlers.HealthCheckHandler))

	// Auth logic
	if cookies.UserStore == nil {
		cookies.UserStore = make(map[string]string)
		users, _ := db.DataBase.ListUsers()
		for _, u := range users {
			cookies.UserStore[u.Username] = string(u.PasswordHash)
		}
	}

	//  SPA Setup

	frontendFS, err := fs.Sub(eApp, "vue/app/dist")
	if err != nil {
		log.Println("Error creating main frontend sub-filesystem:", err)
	}

	loginFS, err = fs.Sub(eLogin, "vue/login/dist")
	if err != nil {
		log.Println("Error creating login sub-filesystem:", err)
	}
	rootAssetHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var fsToUse fs.FS
		if _, err := cookie.ValidateSession(r); err != nil {
			fsToUse = loginFS
		} else {
			fsToUse = frontendFS
		}
		filePath := strings.TrimPrefix(r.URL.Path, "/")

		file, err := fsToUse.Open(filePath)
		if err == nil {
			defer file.Close()

			content, readErr := io.ReadAll(file)
			if readErr != nil {
				http.Error(w, "Error reading embedded file content", http.StatusInternalServerError)
				return
			}
			contentReader := bytes.NewReader(content)

			http.ServeContent(w, r, filePath, time.Time{}, contentReader)
			return
		} else {
		}
	})

	app.Router.PathPrefix("/js/").Handler(rootAssetHandler).Methods("GET")
	app.Router.PathPrefix("/css/").Handler(rootAssetHandler).Methods("GET")
	app.Router.HandleFunc("/favicon.ico", rootAssetHandler).Methods("GET")

	app.Router.HandleFunc("/login", login.HandleLogin).Methods("POST")
	//app.Router.HandleFunc("/login/", login.HandleLogin).Methods("POST")
	app.Router.HandleFunc("/login", LoginHandler).Methods("GET")
	// End SPA Setup

	baseDir := system.System.Config.OutputFolder
	app.Router.PathPrefix("/videos/").
		Handler(http.StripPrefix("/videos/", http.FileServer(http.Dir(baseDir))))

	VideoMux("/api/videos", app.Router)

	app.Router.PathPrefix("/").Handler(http.HandlerFunc(login.RequireAuth(func(w http.ResponseWriter, r *http.Request) {
		// if !cookies.Session.IsLoggedIn(system.System.Config.APIKeys, w, r) {
		// 	http.Redirect(w, r, "/login", http.StatusFound)
		// 	return
		// }
		filePath := strings.TrimPrefix(r.URL.Path, "/")
		if file, err := frontendFS.Open(filePath); err == nil {
			file.Close()
			http.StripPrefix("/", http.FileServer(http.FS(frontendFS))).ServeHTTP(w, r)
		} else {
			indexFile, err := frontendFS.Open("index.html")
			if err != nil {
				fmt.Println(err)
			}
			indexContent, err := io.ReadAll(indexFile)
			if err != nil {
				fmt.Println(err)
			}
			w.Header().Set("Content-Type", "text/html")
			w.Write(indexContent)
		}
	})))

	// CORS
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:" + fmt.Sprint(system.System.Config.Port), fmt.Sprintf("http://%s:%d", "localhost", system.System.Config.Port), "http://localhost:*"},
		AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
		AllowCredentials: true,
		Debug:            false,
	})

	srv := &http.Server{
		Handler:      c.Handler(app.Router),
		Addr:         ":" + fmt.Sprint(system.System.Config.Port),
		WriteTimeout: 1 * time.Hour,
		ReadTimeout:  60 * time.Second,
	}

	prettyprint.P.BoldWhite.Print("Local web server avalable at: ")
	prettyprint.P.Green.Println(fmt.Sprintf("http://%s:%d", utils.GetLocalIp(), system.System.Config.Port))

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("listen: %s\n", err)
	}
	<-ctx.Done()
	log.Println("Shutting down server...")
	log.Println("Server exited gracefully")
}
