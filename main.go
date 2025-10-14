package main

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"remoteCtrl/internal"
	"remoteCtrl/internal/embedded"
	"remoteCtrl/internal/media/video_download"
	"remoteCtrl/internal/system"
	"remoteCtrl/internal/system/cookies"
	"remoteCtrl/internal/system/version"
	"remoteCtrl/internal/utils"
	"remoteCtrl/internal/web/handlers"
	webController "remoteCtrl/internal/web/handlers/controller"
	"remoteCtrl/internal/web/handlers/login"
	"remoteCtrl/internal/web/handlers/status"
	"remoteCtrl/internal/web/handlers/streamers"
	"remoteCtrl/internal/web/handlers/users"
	"remoteCtrl/internal/web/telegram"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func ini() {

	ytDLP_path := utils.CheckPath("yt-dlp")

	ffmpeg_path := utils.CheckPath("ffmpeg")
	if ytDLP_path == "" {
		fmt.Println("missing yt-dlp. Please install before running this program.")
	}
	if ffmpeg_path == "" {
		fmt.Println("missing yt-dlp. Please install before running this program.")
	}
	if ytDLP_path == "" || ffmpeg_path == "" {
		log.Fatal("missing dependencies")
	}
}
func main() {
	system.System.WaitForNetwork = false
	err := internal.Init()
	if err != nil {
		log.Fatal(err)
	}

	// hls.GetToken("https://chaturbate.com/", "cutebrutalitys")
	// hlsS := hls.New("https://chaturbate.com/", "test.mp4")
	// fmt.Println("hls url", hlsS.URL)
	// return
	go serveHTTP(system.System.Context)

	<-system.System.Context.Done()

	telegram.Bot.SendMsg("Server shutdown")
}

func serveHTTP(ctx context.Context) {

	var app handlers.API
	app.Router = mux.NewRouter()
	app.Router.HandleFunc("/api/download", handlers.DownloadHandler)
	app.Router.HandleFunc("/api/progress", video_download.Handler)
	app.Router.HandleFunc("/api/add-streamer", streamers.AddStreamer)
	app.Router.HandleFunc("/api/get-streamers", streamers.GetStreamers)
	app.Router.HandleFunc("/api/remove-streamer", streamers.RemoveStreamer)
	app.Router.HandleFunc("/api/control", webController.ControlHandler)
	app.Router.HandleFunc("/api/get-online-status", streamers.CheckOnlineStatus)
	app.Router.HandleFunc("/api/import", streamers.UploadHandler)
	app.Router.HandleFunc("/api/export", streamers.DownloadHandler)
	app.Router.HandleFunc("/api/status", status.StatusHandler)
	app.Router.HandleFunc("/api/get-videos", webController.GetFiles)
	app.Router.HandleFunc("/api/logs", webController.HandleLogs)
	app.Router.HandleFunc("/api/delete-videos", webController.DeleteFiles)

	app.Router.PathPrefix("/api/generate-api-key").Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookies.GenAPIKeyHandler(system.System.DB.APIKeys, w, r)
	}))
	app.Router.PathPrefix("/api/delete-api-key").Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookies.DeleteAPIKeyHandler(system.System.DB.APIKeys, w, r)
	}))
	app.Router.PathPrefix("/api/keys").Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookies.GetAPIkeys(system.System.DB.APIKeys, w, r)
	}))
	app.Router.HandleFunc("/api/get-users", users.GetUsers)
	app.Router.HandleFunc("/api/add-user", users.AddUser)
	app.Router.HandleFunc("/api/update-user", users.UpdateUsers)
	app.Router.HandleFunc("/api/health", handlers.HealthCheckHandler)

	// Auth logic
	if cookies.UserStore == nil {
		cookies.UserStore = make(map[string]string)
		for _, u := range system.System.DB.Users.Users {
			cookies.UserStore[u.Name] = u.Key
		}
	}

	app.Router.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {

		if r.Method == http.MethodGet {
			handlers.GetLogin(w, r)
		} else if r.Method == http.MethodPost {
			login.PostLogin(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	app.Router.PathPrefix("/videos/").Handler(http.StripPrefix("/videos/", http.FileServer(http.Dir(system.System.DB.Settings.App.Files_folder))))
	handlers.VideoMux("/api/videos", app.Router)

	// VUE
	frontendFS, _ := fs.Sub(embedded.VueDistFiles, "app/dist")
	app.Router.PathPrefix("/").Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if !cookies.Session.IsLoggedIn(system.System.DB.APIKeys, w, r) {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		filePath := strings.TrimPrefix(r.URL.Path, "/")
		if file, err := frontendFS.Open(filePath); err == nil {
			file.Close()
			http.StripPrefix("/", http.FileServer(http.FS(frontendFS))).ServeHTTP(w, r)
		} else {
			indexFile, _ := frontendFS.Open("index.html")
			indexContent, _ := io.ReadAll(indexFile)
			w.Header().Set("Content-Type", "text/html")
			w.Write(indexContent)
		}
	}))

	// CORS
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:" + fmt.Sprint(system.System.DB.Settings.App.Port), fmt.Sprintf("http://%s:%d", "localhost", system.System.DB.Settings.App.Port), "http://localhost:*"},
		AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
		AllowCredentials: true,
		Debug:            false,
	})

	srv := &http.Server{
		Handler:      c.Handler(app.Router),
		Addr:         ":" + fmt.Sprint(system.System.DB.Settings.App.Port),
		WriteTimeout: 1 * time.Hour,
		ReadTimeout:  60 * time.Second,
	}

	fmt.Println("Starting at", "http://localhost:"+fmt.Sprint(system.System.DB.Settings.App.Port))
	fmt.Println("Verison:", version.Version)

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("listen: %s\n", err)
	}
	<-ctx.Done()
	log.Println("Shutting down server...")
	log.Println("Server exited gracefully")
}
