package video_download

import (
	"fmt"
	"log"
	"os"
)

type Progress struct {
	Running   bool   `json:"running"`
	Total     int    `json:"total"`
	Current   int    `json:"current"`
	Progress  int    `json:"progress"`
	Text      string `json:"progressText"`
	QueueText string `json:"queueText"`
}

func (r Progress) Init(running bool, total, progress, current int, queue, text string) {
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
	// Data = r
}
func (r Progress) SetText(s string) Progress {
	r.Text = s
	return r
}
func (r Progress) ApendText(s string) Progress {
	r.Text = "\n" + s
	return r
}
func PrintError(s any) {
	//Data.Init(Data.Running, Data.Total, Data.Progress, Data.Current, Data.QueueText, fmt.Sprint(s))
}
