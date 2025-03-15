package recorder

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/Eyevinn/mp4ff/mp4"
	"github.com/fsnotify/fsnotify"

	"GoStreamRecord/internal/db"
	"GoStreamRecord/internal/web/provider"
)

type Recorder struct {
	Website      provider.Provider `json:"website"`
	IsRestarting bool              `json:"restarting"`
	IsRecording  bool              `json:"is_recording"`
	StopSignal   bool              `json:"-"`
	Cmd          *exec.Cmd         `json:"-"`
}

// Todo: Implement with frontend and update input dir.
//
// WatchFile monitors the file and calls processFile when new data is appended.
func WatchFile(filePath string) {
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
	} else {
		// TODO: impolement to broadcast updates to clients.
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

// writeYoutubeDLdb writes the youtube-dl configuration file.
func WriteYoutubeDLdb() error {
	f, err := os.Create("youtube-dl.config")
	if err != nil {
		return err
	}
	defer f.Close()

	folder := db.Config.Settings.App.Videos_folder
	dbLine := fmt.Sprintf("-o \"%s", folder) + "/%(id)s/%(title)s.%(ext)s\""
	_, err = f.Write([]byte(dbLine))
	return err
}
