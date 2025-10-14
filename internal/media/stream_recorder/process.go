package stream_recorder

import (
	"fmt"
	"log"
	"os"
	"remoteCtrl/internal/media/stream_recorder/recorder"
	"remoteCtrl/internal/system"
	"syscall"
)

func (b *Controller) AddProcess(provider_type, streamerName string) {
	// Only add if not already present
	next_index := len(b.status)
	for _, rec := range b.status {
		if rec.Website.Username == streamerName {
			return
		}
	}
	b.status = append(b.status, recorder.Recorder{})
	b.status[next_index].Website.New(provider_type, streamerName)
}
func (b *Controller) Status(name string) recorder.Recorder {
	return getProcess(name, b)
}

// ListRecorders returns the current list of recorder statuses.
func (b *Controller) ListRecorders() []recorder.Recorder {
	if b == nil {
		return []recorder.Recorder{}
	}
	b.mux.Lock()
	defer b.mux.Unlock()
	return b.status
}
func (b *Controller) StopRunningEmpty() {

	if b == nil {
		return
	}
	b.checkProcesses()
}

// StopProcess sends a SIGINT to active recording processes and waits for them to finish.
func (b *Controller) StopProcess(processName string) {
	if b == nil {
		return
	}
	b.mux.Lock()
	// Create a copy of status indices to avoid modification during iteration.
	statusCopy := make([]recorder.Recorder, len(b.status))
	copy(statusCopy, b.status)
	b.mux.Unlock()

	for _, rec := range statusCopy {
		// Stop only the specified process (or all if processName is empty).
		if processName != "" && rec.Website.Username != processName {
			continue
		}
		b.stopProcessIfRunning(rec)

	}
}

// checkProcesses looks through the list of processes and removes any that have finished.
func (b *Controller) checkProcesses() int {
	if b == nil {
		return 0
	}
	b.mux.Lock()
	defer b.mux.Unlock()
	for i := 0; i < len(b.status); i++ {
		// Use signal 0 to check if process is still running.
		if !b.status[i].StopSignal || b.status[i].Cmd == nil {
			continue
		}
		if err := b.status[i].Cmd.Process.Signal(syscall.SIGTERM); err != nil {
			log.Printf("Process for %s has stopped", b.status[i].Website.Username)
			b.status = append(b.status[:i], b.status[i+1:]...)
			i--
		}
	}
	return len(b.status)
}

func (b *Controller) stopProcessIfRunning(bs recorder.Recorder) {

	if b == nil {
		return
	}
	for i, s := range b.status {
		if bs.Cmd != nil && s.Website.Username == bs.Website.Username {
			b.status[i].StopSignal = true
			if s.Cmd.Process == nil {
				log.Printf("Process for %s had no process", bs.Website.Username)
				fmt.Printf("Process for %s had no process", bs.Website.Username)
				continue
			}
			if err := s.Cmd.Process.Signal(syscall.SIGINT); err != nil {
				i--
			}

			b.StopBot(bs.Website.Username)
			log.Printf("Command for %s has stopped\n", bs.Website.Username)
			fmt.Printf("Command for %s has stopped\n", bs.Website.Username)
			if len(b.status) == 0 {
				break
			}
			b.status = append(b.status[:i], b.status[i+1:]...)
			break
		}
		if s.Cmd == nil && s.Website.Username == bs.Website.Username {
			b.status[i].StopSignal = true
			log.Printf("Process for %s was already stopped", bs.Website.Username)
			b.status = append(b.status[:i], b.status[i+1:]...)
			break
		}
	}

}

// isRecorderActive returns true if a recorder for the given streamer is already running.
func (b *Controller) isRecorderActive(streamerName string) bool {

	if b == nil {
		return false
	}
	for _, rec := range b.status {
		if rec.Website.Username == streamerName && rec.IsRecording {
			return true
		}
	}
	return false
}

func getProcess(name string, b *Controller) recorder.Recorder {

	if b == nil {
		return recorder.Recorder{}
	}
	b.mux.Lock()
	defer b.mux.Unlock()
	for _, s := range b.status {
		if name == s.Website.Username {
			return s
		}
	}
	return recorder.Recorder{StopSignal: false, IsRecording: false, Cmd: nil}
}

// writeYoutubeDLdb writes the youtube-dl configuration file.
func (b *Controller) writeYoutubeDLdb() error {
	f, err := os.Create("youtube-dl.config")
	if err != nil {
		return err
	}
	defer f.Close()

	folder := system.System.DB.Settings.App.Files_folder
	dbLine := fmt.Sprintf("-o \"%s", folder) + "/%(id)s/%(title)s.%(ext)s\""
	_, err = f.Write([]byte(dbLine))
	return err
}

// StopBot signals the bot to stop starting new recordings and then gracefully stops active processes.
func (b *Controller) StopBot(streamerName string) {
	if b == nil {
		return
	}
	// Signal cancellation.
	b.cancel()
	log.Println("Stopping bot..")
	// Give current recorders time to finish (or exit gracefully).
	for i := range b.status {
		if b.status[i].Website.Username == streamerName {
			b.status[i].StopSignal = true
		}
	}
}
