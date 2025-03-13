package recorder

import (
	"GoStreamRecord/internal/web/provider"
	"fmt"
	"os"
	"os/exec"

	"github.com/Eyevinn/mp4ff/mp4"
	"github.com/fsnotify/fsnotify"
	// Assume mp4ff is imported and used for MP4 parsing
)

type Recorder struct {
	Website      provider.Provider `json:"website"`
	StopSignal   bool              `json:"-"`
	IsRestarting bool              `json:"restarting"`
	Cmd          *exec.Cmd         `json:"-"`
	IsRecording  bool              `json:"is_recording"`
}

// watchFile monitors the file and calls processFile when new data is appended.
func watchFile(filePath string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		panic(err)
	}
	defer watcher.Close()

	// Start watching the directory containing the file.
	dir := "." // adjust as needed
	err = watcher.Add(dir)
	if err != nil {
		panic(err)
	}

	// Initial processing
	processFile(filePath)

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			// If the file is written to, re-read it.
			if event.Name == filePath && event.Op&fsnotify.Write == fsnotify.Write {
				fmt.Println("File modified, re-processing...")
				processFile(filePath)
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			fmt.Println("Error watching file:", err)
		}
	}
}

// processFile opens the file, reads its content, and tries to parse MP4 structure.
func processFile(filePath string) {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// Option 1: Use a library like mp4ff to parse the file.
	// Note: mp4ff might error out if the file is incomplete.
	parsedFile, err := mp4.DecodeFile(file)
	if err != nil {
		fmt.Println("Error parsing MP4 file:", err)
		// Optionally, handle partial parsing or wait for more data.
	} else {
		// Extract info such as duration, video dimensions, etc.
		fmt.Printf("Parsed MP4 file: Duration = %v, Size = %d bytes\n", parsedFile.Moov.Mvhd.Duration, getFileSize(filePath))
		// You can then extract frames if you have the proper decoder.
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
