package recorder

import (
	"bytes"
	"log"
	"os/exec"
	"strings"
)

func (rec *Recorder) StartRecording(streamerName string) {
	
	log.Printf("Starting recording for %s", streamerName)

	args := strings.Fields("yt-dlp --no-part")
	args = append(args, rec.Website.StreamURL, "--config-location", "youtube-dl.config")

	rec.Cmd = exec.Command(args[0], args[1:]...)
	var stdout, stderr bytes.Buffer
	rec.Cmd.Stdout = &stdout
	rec.Cmd.Stderr = &stderr

	if err := rec.Cmd.Start(); err != nil {
		log.Printf("Error starting recording for %s: %v\n", streamerName, err)
		return
	}

	// Wait for the process to complete
	err := rec.Cmd.Wait()
	if err != nil {
		log.Printf("Recording process for %s exited with error: %v", streamerName, err)
		log.Printf("yt-dlp stderr: %s", stderr.String()) // Print stderr output
		//fmt.Printf("yt-dlp stderr: %s", stderr.String()) // Print stderr output
	} else {
		log.Printf("yt-dlp stderr: %s", stdout.String()) // Print stderr output
		log.Printf("Recording process for %s finished successfully", streamerName)
	}
}
