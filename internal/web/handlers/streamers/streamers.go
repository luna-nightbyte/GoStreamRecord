package streamers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"remoteCtrl/internal/db"
	"remoteCtrl/internal/media/stream_recorder"
	"remoteCtrl/internal/media/stream_recorder/recorder"
	"remoteCtrl/internal/web/handlers/status"
	"sync"
)

// Handles POST /api/add-streamer.
// It decodes a JSON payload with a "data" field and returns a dummy response.
func AddStreamer(w http.ResponseWriter, r *http.Request) {
	//if !cookies.Session.IsLoggedIn(system.System.Config.APIKeys, w, r) {
	//	http.Redirect(w, r, "/login", http.StatusFound)
	//	return
	//}
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
	shareStr := r.URL.Query().Get("share")
	share := false
	if shareStr == "true" {
		share = true
	}
	db.DataBase.NewStreamer(reqData.Data, r.URL.Query().Get("provider"), db.DataBase.Users.HttpRequestID(r), share)
	resp := status.Response{
		Message: "success",
		Data:    reqData.Data,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// Handles POST /api/remove-streamer.
// It decodes a JSON payload with the selected option and returns a dummy response.
func RemoveStreamer(w http.ResponseWriter, r *http.Request) {
	// if !cookies.Session.IsLoggedIn(system.System.Config.APIKeys, w, r) {
	// 	http.Redirect(w, r, "/login", http.StatusFound)
	// 	return
	// }
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

	streamers, _ := db.DataBase.Streamers.List()
	_, err := db.DataBase.Streamers.DeleteForUser(db.DataBase.Users.HttpRequestID(r), streamers[reqData.Selected].ID)
	resp := status.Response{
		Message: err.Error(),
		Data:    reqData.Selected,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// Handles GET /api/get-streamers.
func GetStreamers(w http.ResponseWriter, r *http.Request) {
	// if !cookies.Session.IsLoggedIn(system.System.Config.APIKeys, w, r) {
	// 	http.Redirect(w, r, "/login", http.StatusFound)
	// 	return
	// }
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	streamers, _ := db.DataBase.Streamers.List()
	list := []string{}
	for _, s := range streamers {
		//list[s.Name] = s.Provider
		list = append(list, s.Name)
	}
	json.NewEncoder(w).Encode(list)
}

func CheckOnlineStatus(w http.ResponseWriter, r *http.Request) {

	// if !cookies.Session.IsLoggedIn(system.System.Config.APIKeys, w, r) {
	// 	http.Redirect(w, r, "/login", http.StatusFound)
	// 	return
	// }
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
		status.ResponseHandler(w, r, "Streamer name is required", false, nil)
		return
	}

	if reqData.Provider == "" {

		streamers, _ := db.DataBase.Streamers.List()
		for _, streamer := range streamers {
			if streamer.Name == reqData.Streamer {
				reqData.Provider = streamer.Provider
				break
			}
		}
	}

	var re recorder.Recorder
	err := re.Website.New(reqData.Provider, reqData.Streamer)
	if err != nil {
		http.Error(w, "Internal recorder error!", http.StatusInternalServerError)
		return
	}
	is_online := fmt.Sprintf("%v", re.Website.Interface.IsOnline(reqData.Streamer))
	status.ResponseHandler(w, r, is_online, true, nil)
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
		// status.ResponseHandler(w, r, "Stopping process for "+s,true, nil)
		stream_recorder.Streamer.StopProcessByName(rd.Streamer)
		status.ResponseHandler(w, r, "Stopped process for"+s, true, nil)
		rd.wg.Done()

	}(&stopData)
	stopData.wg.Wait()
}
