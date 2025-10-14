package stream_recorder

import (
	"fmt"
	"log"
	"os"
	"remoteCtrl/internal/media/stream_recorder/recorder"
	"remoteCtrl/internal/media/stream_recorder/recorder/provider"
	"remoteCtrl/internal/system"
	"syscall"
)

func (b *Controller) AddProcess(provider_type, streamerName string) {
	// Only add if not already present
	if b.bots[streamerName] != nil {
		return
	}
	b.bots[streamerName] = &recorder.Recorder{
		Website:      provider.Provider{},
		IsRestarting: false,
		IsRecording:  true,
	}
	b.bots[streamerName].Website.New(provider_type, streamerName)
	// next_index := len(b.status)
	// for _, rec := range b.status {
	// 	if rec.Website.Username == streamerName {
	// 		return
	// 	}
	// }
	// b.status = append(b.status, recorder.Recorder{})
	// b.status[next_index].Website.New(provider_type, streamerName)
}
func (b *Controller) Status(name string) recorder.Recorder {
	return getProcess(name, b)
}

// ListRecorders returns the current list of recorder statuses, alternatively it initialyzes a new map if empty.
func (b *Controller) ListRecorders() map[string]*recorder.Recorder {
	if b == nil {
		return nil
	}
	if b.bots == nil {
		return make(map[string]*recorder.Recorder)
	}
	b.mux.Lock()
	defer b.mux.Unlock()
	return b.bots
}
func (b *Controller) StopRunningEmpty() {

	if b == nil {
		return
	}
	b.checkProcesses()
}

// StopProcess sends a SIGINT to active recording processes and waits for them to finish.
func (b *Controller) StopProcessByName(processName string) {
	if b == nil || b.bots[processName] == nil {
		return
	}
	b.stopProcessIfRunning(b.bots[processName])

}

// checkProcesses looks through the list of processes and removes any that have finished.
func (b *Controller) checkProcesses() int {
	if b == nil {
		return 0
	}
	b.mux.Lock()
	defer b.mux.Unlock()
	for name := range b.bots {
		if !b.bots[name].ShouldStop() || b.bots[name].Cmd == nil {
			continue
		}
		b.bots[name].Cmd.Process.Signal(syscall.SIGTERM)
		log.Printf("Process for %s has stopped", b.bots[name].Website.Username)
	}
	return len(b.bots)
}

func (b *Controller) stopProcessIfRunning(bot *recorder.Recorder) {
	if b == nil {
		return
	}
	if bot.Cmd != nil {
		bot.Stop()
		if err := bot.Cmd.Process.Signal(syscall.SIGINT); err != nil {

			log.Println("Error signaling stop: ", err)
		}
	}

}

// isRecorderActive returns true if a recorder for the given streamer is already running.
func (b *Controller) isRecorderActive(name string) bool {
	if b == nil {
		return false
	}
	if b.bots[name] == nil {
		return false
	}
	return b.bots[name].IsRecording

	// for _, rec := range b.status {
	// 	if rec.Website.Username == streamerName && rec.IsRecording {
	// 		return true
	// 	}
	// }
	// return false
}

func getProcess(name string, b *Controller) recorder.Recorder {

	if b == nil {
		return recorder.Recorder{}
	}
	b.mux.Lock()
	defer b.mux.Unlock()
	return *b.bots[name]
	// for _, s := range b.status {
	// 	if name == s.Website.Username {
	// 		return s
	// 	}
	// }
	//return recorder.Recorder{StopSignal: false, IsRecording: false, Cmd: nil}
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
func (b *Controller) StopBot(name string) {
	if b == nil || b.bots[name] == nil {
		return
	}
	// Signal cancellation.
	b.cancel()
	b.bots[name].Stop()
}
