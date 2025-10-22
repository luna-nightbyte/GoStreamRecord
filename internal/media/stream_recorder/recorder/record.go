package recorder

import (
	"fmt"
	"log"
	"os/exec"
	"remoteCtrl/internal/system"
	"strings"
)

func (rec *Recorder) Stop() {
	rec.stopSignal = true
}
func (rec *Recorder) Start() {
	rec.stopSignal = false
}
func (rec *Recorder) ShouldStop() bool {
	return rec.stopSignal
}

func (rec *Recorder) StopTicker() {
	rec.exitSignal <- true
}

// startRecording starts a recording for the given streamer.
func (rec *Recorder) StartRecording(streamerName string) {
	if rec.IsRecording {
		return
	}
	rec.IsRecording = true

	args := strings.Fields("yt-dlp --no-part")

	args = append(args, fmt.Sprintf("%s%s/", rec.Website.Url, streamerName), "--config-location", "youtube-dl.config")
	if system.DEBUG {
		log.Println("Executing yt-dlp command: ", args)
	}
	rec.Cmd = exec.Command(args[0], args[1:]...)

	if err := rec.Cmd.Start(); err != nil {
		log.Printf("Error starting recording for %s: %v\n", streamerName, err)
	}

	rec.Cmd.Wait()

	// TODO
	//utils.VerifyCodec("")
}
