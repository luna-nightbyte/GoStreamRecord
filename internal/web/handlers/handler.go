package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type FormData struct {
	Data map[string]any `json:"formData"`
}

type API struct {
	Router *mux.Router
}

var Api API

// HealthResponse represents the JSON structure for health responses.
type HealthResponse struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}

// HasTimePassed checks if a certain duration has passed since the given timestamp.
func HasTimePassed(startTime time.Time, duration time.Duration) bool {
	return time.Since(startTime) >= duration
}

// HealthCheckHandler is the HTTP handler for the health check endpoint.
func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	// if !cookies.Session.IsLoggedIn(system.System.DB.APIKeys, w, r) {
	// 	http.Redirect(w, r, "/login", http.StatusFound)
	// 	return
	// }
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

func RedirectHome(w http.ResponseWriter, r *http.Request) {
	// if !cookies.Session.IsLoggedIn(system.System.DB.APIKeys, w, r) {
	// 	http.Redirect(w, r, "/login", http.StatusFound)
	// 	return
	// }
	http.Redirect(w, r, "/", http.StatusFound) // 302 Found
}
