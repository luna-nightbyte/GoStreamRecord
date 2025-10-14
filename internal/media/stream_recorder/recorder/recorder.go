package recorder

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"remoteCtrl/internal/media/stream_recorder/recorder/provider"

	"github.com/Eyevinn/mp4ff/mp4"
	// Assume mp4ff is imported and used for MP4 parsing
)

type Recorder struct {
	Website      provider.Provider `json:"website"`
	StopSignal   bool              `json:"-"`
	IsRestarting bool              `json:"restarting"`
	Cmd          *exec.Cmd         `json:"-"`
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
