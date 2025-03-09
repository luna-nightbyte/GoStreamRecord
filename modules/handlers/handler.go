package handlers

import (
	"GoStreamRecord/modules/db"
	"GoStreamRecord/modules/handlers/controller"
	"GoStreamRecord/modules/handlers/cookies"
	"GoStreamRecord/modules/handlers/login"
	web_status "GoStreamRecord/modules/handlers/status"
	"GoStreamRecord/modules/handlers/streamers"
	"GoStreamRecord/modules/handlers/users"
	"encoding/json"
	"net/http"
	"text/template"
)

func Handle() {
	// API endpoints

	http.Handle("/videos/", http.StripPrefix("/videos/", http.FileServer(http.Dir(db.Config.Settings.App.Videos_folder))))

	http.HandleFunc("/api/add-streamer", streamers.AddStreamer)
	http.HandleFunc("/api/get-streamers", streamers.GetStreamers)
	http.HandleFunc("/api/remove-streamer", streamers.RemoveStreamer)
	http.HandleFunc("/api/control", controller.ControlHandler)
	http.HandleFunc("/api/get-online-status", streamers.CheckOnlineStatus)
	http.HandleFunc("/api/import", streamers.UploadHandler)
	http.HandleFunc("/api/export", streamers.DownloadHandler)
	http.HandleFunc("/api/status", web_status.StatusHandler)
	http.HandleFunc("/api/get-videos", controller.GetVideos)
	http.HandleFunc("/api/logs", controller.HandleLogs)
	http.HandleFunc("/api/delete-videos", controller.DeleteVideos)
	http.HandleFunc("/api/generate-api-key", cookies.GenAPIKeyHandler)
	http.HandleFunc("/api/delete-api-key", cookies.DeleteAPIKeyHandler)
	http.HandleFunc("/api/keys", cookies.GetAPIkeys)

	http.HandleFunc("/api/get-users", users.GetUsers)
	http.HandleFunc("/api/add-user", users.AddUser)
	http.HandleFunc("/api/update-user", users.UpdateUsers)
	http.HandleFunc("/api/health", HealthCheckHandler)

	if cookies.UserStore == nil {
		cookies.UserStore = make(map[string]string)
	}

	for _, u := range db.Config.Users.Users {
		cookies.UserStore[u.Name] = u.Key
	}
	//	fs := http.FileServer(http.Dir(filepath.Dir(file.Index_path)))

	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {

		if r.Method == http.MethodGet {
			GetLogin(w, r)
		} else if r.Method == http.MethodPost {
			login.PostLogin(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if cookies.Session.IsLoggedIn(w, r) {
			GetIndex(w, r)
			return
		} else {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
	})
}

var IndexHTML, LoginHTML string

type Template struct {
	W    http.ResponseWriter
	Tmpl *template.Template
}

func (t *Template) Execute(data any) error {
	return t.Tmpl.Execute(t.W, data)
}

func GetIndex(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.New("index").Parse(IndexHTML))
	indexTemplate := Template{W: w, Tmpl: tmpl}
	if err := indexTemplate.Execute(nil); err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
	}
}

func GetLogin(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.New("login").Parse(LoginHTML))
	loginTemplate := Template{W: w, Tmpl: tmpl}
	if err := loginTemplate.Execute(nil); err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
	}
}

// HealthResponse represents the JSON structure for health responses.
type HealthResponse struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}

// HealthCheckHandler is the HTTP handler for the health check endpoint.
func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	if !cookies.Session.IsLoggedIn(w, r) {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	response := HealthResponse{
		Status:  "ok",
		Message: "Service is up and running",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
	}

}
