package bot

import (
	"log"
	"syscall"

	"GoStreamRecord/internal/recorder"
)

func (b *controller) AddProcess(provider_type, streamerName string) {

	b.mux.Lock()
	defer b.mux.Unlock()
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
func (b *controller) Status(name string) recorder.Recorder {
	return getProcess(name, b)
}

// ListRecorders returns the current list of recorder statuses.
func (b *controller) ListRecorders() []recorder.Recorder {
	b.mux.Lock()
	defer b.mux.Unlock()
	return b.status
}
func (b *controller) StopRunningEmpty() {
	b.checkProcesses()
}

// checkProcesses looks through the list of processes and removes any that have finished.
func (b *controller) checkProcesses() int {
	b.mux.Lock()
	defer b.mux.Unlock()
	for i := 0; i < len(b.status); i++ {
		// Use signal 0 to check if process is still running.
		log.Println("index1: ", i)
		log.Println("Checking process:", b.status[i].IsRecording, b.status[i].StopSignal, b.status[i].IsRecording, b.status[i].Cmd)
		if !b.status[i].StopSignal {
			continue
		}
		b.StopProcessIfRunning(&b.status[i])
		log.Printf("Process for %s has stopped", b.status[i].Website.Username)
		b.status = append(b.status[:i], b.status[i+1:]...)
		i--
		log.Println("index2:", i)

	}
	return len(b.status)
}

func (b *controller) StopProcessIfRunning(rec *recorder.Recorder) {
	for i, s := range b.status {
		log.Println("Checking process:", s.IsRecording, s.StopSignal, s.Cmd)
		if rec.Cmd != nil && s.Website.Username == rec.Website.Username {
			b.status[i].StopSignal = true
			if err := s.Cmd.Process.Signal(syscall.SIGINT); err != nil {
				i--
			}

			break
		}
		if s.Cmd == nil {
			b.status[i].StopSignal = true
			log.Printf("Process for %s was already stopped", rec.Website.Username)
			b.status = append(b.status[:i], b.status[i+1:]...)
			break
		}
	}

}

// isRecorderActive returns true if a recorder for the given streamer is already running.
func (b *controller) isRecorderActive(streamerName string) bool {
	for _, rec := range b.status {
		if rec.Website.Username == streamerName && rec.IsRecording {
			return true
		}
	}
	return false
}

func getProcess(name string, b *controller) recorder.Recorder {
	b.mux.Lock()
	defer b.mux.Unlock()
	for _, s := range b.status {
		if name == s.Website.Username {
			return s
		}
	}
	return recorder.Recorder{StopSignal: false, IsRecording: false, Cmd: nil}
}

// StopBot signals the bot to stop starting new recordings and then gracefully stops active processes.
func (b *controller) StopRecorder(streamerName string) {
	// Signal cancellation.
	//b.cancel()
	log.Println("Stopping recorder..")
	// Give current recorders time to finish (or exit gracefully).
	for i := range b.status {
		if b.status[i].Website.Username == streamerName {
			b.status[i].StopSignal = true
		}
	}
}

// StopBot signals the bot to stop starting new recordings and then gracefully stops active processes.
func (b *controller) StopBot(streamerName string) {
	log.Println("Stopping bot..")
	// Signal cancellation.
	b.cancel()

}
