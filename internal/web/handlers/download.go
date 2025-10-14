package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"remoteCtrl/internal/media/video_download"
	"remoteCtrl/internal/system"
	"strings"
	"time"
)

type videoFormData struct {
	Option string `json:"option"`
	Save   string `json:"save"`
	Bulk   bool   `json:"bulk"`
	Search string `json:"search"`
	URL    string `json:"url"`
}

type req struct {
	I         int
	Timestamp time.Time
}

var request req
var VideoFormData video_download.DownloadForm

func init() {
	request.Timestamp = time.Now()
}

func DownloadHandler(w http.ResponseWriter, r *http.Request) {

	if !HasTimePassed(request.Timestamp, 1*time.Second) {
		return
	}
	request.I++
	request.Timestamp = time.Now()
	http.Redirect(w, r, "/download", http.StatusFound) // 302 Found
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read body", http.StatusBadRequest)
		return
	}
	//decoding JSON array to persons array
	err = json.Unmarshal(body, &VideoFormData)

	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	if VideoFormData.Option == "" || (VideoFormData.URL == "" && VideoFormData.Search == "") {
		w.Write([]byte("Missing input.."))
		return
	}
	if VideoFormData.Search != "" {
		VideoFormData.Bulk = true
		VideoFormData.Save = strings.ReplaceAll(VideoFormData.Search, " ", "_")

	}
	if VideoFormData.Save == "" {
		VideoFormData.Save = "default"

		// Dont overwrite existing defaults.
		_, err = os.ReadFile(system.System.DB.Settings.App.Files_folder + "/default.mp4")
		if !os.IsNotExist(err) {
			i := 0
			for {
				_, err := os.ReadFile(fmt.Sprintf(system.System.DB.Settings.App.Files_folder+"/default_%d.mp4", i))
				if os.IsNotExist(err) {
					VideoFormData.Save = fmt.Sprintf("default_%d", i)
					break
				}
				i++
			}
		}
	}
	video_download.DownloadIsRunning = true

	w.Header().Set("Content-Type", "text/plain")
	if VideoFormData.Option == "onlyfans" {
		// go onlyfans.Download(FormData.Search, "img")
	} else if VideoFormData.Option == "pornone" {

	} else {
		go video_download.Download(VideoFormData)
	}

}
