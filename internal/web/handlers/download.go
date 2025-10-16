package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"remoteCtrl/internal/media/video_download"
	"remoteCtrl/internal/web/handlers/status"
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

var ActiveDownloaders int = 0

func init() {
	request.Timestamp = time.Now()
}

func DownloadHandler(w http.ResponseWriter, r *http.Request) {
	if !HasTimePassed(request.Timestamp, 1*time.Second) {
		return
	}
	var VideoFormData video_download.DownloadForm
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
		VideoFormData.Save = getNextDefaultFileName("default.mp4")
		// Dont overwrite existing defaults.

	}
	var dw video_download.VideoDownloader
	dw.IsDownloading = true

	w.Header().Set("Content-Type", "text/plain")
	if VideoFormData.Option == "onlyfans" {
		// go onlyfans.Download(FormData.Search, "img")
	} else {
		status.Status.Is_Downloading = true
		ActiveDownloaders++
		dw.Download(VideoFormData)
		ActiveDownloaders--
		if ActiveDownloaders == 0 {
			status.Status.Is_Downloading = false
		}
	}

}

func getNextDefaultFileName(filename string) string {
	_, err := os.ReadFile(filename)
	if !os.IsNotExist(err) {
		getNextDefaultFileName(filename)
	}
	return filename
}
