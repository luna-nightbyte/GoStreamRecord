package recorder

import (
	"fmt"
	"log"
	"os/exec"
	"remoteCtrl/internal/utils"
	"strings"
)

func (rec *Recorder) Stop() {
	rec.stopSignal = true
}
func (rec *Recorder) ShouldStop() bool {
	return rec.stopSignal
}

// startRecording starts a recording for the given streamer.
func (rec *Recorder) StartRecording(streamerName string) {
	if rec.IsRecording {
		return
	}
	rec.IsRecording = true
	log.Printf("Starting recording for %s\n", streamerName)
	fmt.Printf("Starting recording for %s\n", streamerName)

	ytDlpPath := utils.CheckPath("yt-dlp")

	args := strings.Fields(fmt.Sprint(ytDlpPath) + " --no-part")

	args = append(args, fmt.Sprintf("%s%s/", rec.Website.Url, streamerName), "--config-location", "youtube-dl.config")

	rec.Cmd = exec.Command(args[0], args[1:]...)

	if err := rec.Cmd.Start(); err != nil {
		log.Printf("Error starting recording for %s: %v\n", streamerName, err)
	}

	rec.Cmd.Wait()

	// TODO
	utils.VerifyCodec("")
}
