package recorder

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"remoteCtrl/internal/db/jsondb"
	"remoteCtrl/internal/media/stream_recorder/recorder/provider"
	"remoteCtrl/internal/system"
	"remoteCtrl/internal/system/settings"
	"time"

	"github.com/Eyevinn/mp4ff/mp4"
	// Assume mp4ff is imported and used for MP4 parsing
)

type Recorder struct {
	Cmd          *exec.Cmd         `json:"-"`
	stopSignal   bool              `json:"-"` // Stops the recording, but ticker can still run
	exitSignal   chan bool         `json:"-"` // Completely stop the bot ticker.
	Website      provider.Provider `json:"website"`
	IsRestarting bool              `json:"restarting"`
	IsRecording  bool              `json:"is_recording"`
}

// processFile opens the file, reads its content, and tries to parse MP4 structure.
func processFile(filePath string) {
	file, err := os.Open(filePath)
	if err != nil {
		log.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	parsedFile, err := mp4.DecodeFile(file)
	if err != nil {
		log.Println("Error parsing MP4 file:", err)
	} else {
		fmt.Printf("Parsed MP4 file: Duration = %v, Size = %d bytes\n", parsedFile.Moov.Mvhd.Duration, getFileSize(filePath))
	}
}

// getFileSize returns the size of the file.
func getFileSize(filePath string) int64 {
	info, err := os.Stat(filePath)
	if err != nil {
		return 0
	}
	return info.Size()
}


// StartRecordTicker Starts a new ticker where the bots will wait for the model to be online. 
// Once online it starts a new recording for that bot until either exit sognal is recieved, or theselected bot should be fully stopped.
func (b *Recorder) StartRecordTicker(ctx context.Context) {

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	b.Website.Interface.TrueName(b.Website.Username)
	go b.start()

	for {
		select {
		case <-ctx.Done():
		case exit := <-b.exitSignal:
			if exit {
				log.Println("Bot ticker stopped")
				return
			}
		case <-ticker.C:
			if b.IsRecording || b.ShouldStop() {
				continue
			}
			b.start()
			log.Printf("Recording for %s finished", b.Website.Username)
		}
	}

}
func (b *Recorder) start() {
	jsondb.Load(settings.CONFIG_SETTINGS_PATH, &system.System.DB.Settings)
	if !b.Website.Interface.IsOnline(b.Website.Username) {
		return
	}
	log.Println("Starting new recording for", b.Website.Username)

	b.StartRecording(b.Website.Username)

	b.IsRecording = false
	b.Stop()
}
