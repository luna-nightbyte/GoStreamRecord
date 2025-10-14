package recorder

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
)

// startRecording starts a recording for the given streamer.
func (rec *Recorder) StartRecording(streamerName string) {
	log.Printf("Starting recording for %s", streamerName)

	// YT-dl

	args := strings.Fields("/home/thomas/.local/bin/yt-dlp --no-part")

	args = append(args, fmt.Sprintf("%s%s/", rec.Website.Url, streamerName), "--config-location", "youtube-dl.config")

	rec.Cmd = exec.Command(args[0], args[1:]...)

	// Start the recording process.
	if err := rec.Cmd.Start(); err != nil {
		log.Printf("Error starting recording for %s: %v\n", streamerName, err)
	}
	rec.Cmd.Wait()
}
