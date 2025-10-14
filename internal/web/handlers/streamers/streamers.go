package streamers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"remoteCtrl/internal/media/stream_recorder"
	"remoteCtrl/internal/media/stream_recorder/recorder"
	"remoteCtrl/internal/system"
	"remoteCtrl/internal/system/cookies"
	"remoteCtrl/internal/web/handlers/status"
	"sync"
)

// Handles POST /api/add-streamer.
// It decodes a JSON payload with a "data" field and returns a dummy response.
func AddStreamer(w http.ResponseWriter, r *http.Request) {
	if !cookies.Session.IsLoggedIn(system.System.DB.APIKeys, w, r) {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	type RequestData struct {
		Data string `json:"data"`
	}
	var reqData RequestData
	if err := json.NewDecoder(r.Body).Decode(&reqData); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	resp := status.Response{
		Message: system.System.DB.AddStreamer(reqData.Data, r.URL.Query().Get("provider")),
		Data:    reqData.Data,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// Handles POST /api/remove-streamer.
// It decodes a JSON payload with the selected option and returns a dummy response.
func RemoveStreamer(w http.ResponseWriter, r *http.Request) {
	if !cookies.Session.IsLoggedIn(system.System.DB.APIKeys, w, r) {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	type RequestData struct {
		Selected string `json:"selected"`
	}
	var reqData RequestData
	if err := json.NewDecoder(r.Body).Decode(&reqData); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	resp := status.Response{
		Message: system.System.DB.RemoveStreamer(reqData.Selected),
		Data:    reqData.Selected,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// Handles GET /api/get-streamers.
func GetStreamers(w http.ResponseWriter, r *http.Request) {
	if !cookies.Session.IsLoggedIn(system.System.DB.APIKeys, w, r) {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	list := []string{}
	for _, s := range system.System.DB.Streamers.List {
		list = append(list, s.Name)
	}
	json.NewEncoder(w).Encode(list)
}

func CheckOnlineStatus(w http.ResponseWriter, r *http.Request) {

	if !cookies.Session.IsLoggedIn(system.System.DB.APIKeys, w, r) {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	type RequestData struct {
		Streamer string `json:"streamer"`
		Provider string `json:"provider"`
	}
	var reqData RequestData
	if err := json.NewDecoder(r.Body).Decode(&reqData); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	if reqData.Streamer == "" {
		log.Println("Streamer name is required")
		status.ResponseHandler(w, r, "Streamer name is required", nil)
		return
	}

	if reqData.Provider == "" {
		log.Println("Provider name is required")
		status.ResponseHandler(w, r, "Streamer name is required", nil)
		return
	}

	var re recorder.Recorder
	err := re.Website.New(reqData.Provider, reqData.Streamer)
	if err != nil {
		http.Error(w, "Internal recorder error!", http.StatusInternalServerError)
		return
	}
	is_online := fmt.Sprintf("%v", re.Website.Interface.IsOnline(reqData.Streamer))
	status.ResponseHandler(w, r, is_online, nil)
}

type RequestData struct {
	wg       *sync.WaitGroup `json:"-"`
	mu       sync.Mutex      `json:"-"`
	Streamer string          `json:"streamer"`
}

var stopData RequestData

func StopProcess(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}
	if err := json.NewDecoder(r.Body).Decode(&stopData); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	stopData.wg.Add(1)
	go func(rd *RequestData) {
		rd.mu.Lock()
		s := rd.Streamer
		rd.mu.Unlock()
		status.ResponseHandler(w, r, "Stopping process for "+s, nil)
		stream_recorder.Streamer.StopProcess(rd.Streamer)
		status.ResponseHandler(w, r, "Stopped process for"+s, nil)
		rd.wg.Done()

	}(&stopData)
	stopData.wg.Wait()
}
