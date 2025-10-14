package video_download

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"remoteCtrl/internal/system"
	"remoteCtrl/internal/system/cookies"
)

type Runner struct {
	Running   bool   `json:"running"`
	Total     int    `json:"total"`
	Current   int    `json:"current"`
	Progress  int    `json:"progress"`
	Text      string `json:"progressText"`
	QueueText string `json:"queueText"`
}

var Data Runner

func (r Runner) Init(running bool, total, progress, current int, queue, text string) {
	if total == 0 {
		total = 100
		current = 1
	}
	r = r.SetText(text)
	r.Running = running
	r.Total = total
	r.Current = current
	r.Progress = current / total

	f, err := os.ReadDir("videos")
	if err != nil {
		log.Println(err)
	}
	t := fmt.Sprintf(`Local files: %d`, len(f))
	r.QueueText = t
	Data = r
}
func (r Runner) SetText(s string) Runner {
	r.Text = s
	return r
}
func (r Runner) ApendText(s string) Runner {
	r.Text = "\n" + s
	return r
}
func Handler(w http.ResponseWriter, r *http.Request) {

	if !cookies.Session.IsLoggedIn(system.System.DB.APIKeys, w, r) {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	ru := Runner{Total: Data.Total, Progress: Data.Progress, Current: Data.Current, Running: Data.Running, QueueText: Data.QueueText, Text: Data.Text}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ru)
}

func PrintError(s any) {
	Data.Init(Data.Running, Data.Total, Data.Progress, Data.Current, Data.QueueText, fmt.Sprint(s))
}
