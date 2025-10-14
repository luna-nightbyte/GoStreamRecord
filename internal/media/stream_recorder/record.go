package stream_recorder

import (
	"log"
	"remoteCtrl/internal/db"
	"remoteCtrl/internal/media/stream_recorder/recorder"
	"remoteCtrl/internal/system"
	"remoteCtrl/internal/system/settings"
	"sync"
	"time"
)

func (b *Controller) RecordLoop(streamerName string) {
	if b.bots[streamerName] != nil {
		return
	}
	if b.bots[streamerName].IsRecording {
		return
	}
	// Write youtube-dl db.
	if err := b.writeYoutubeDLdb(); err != nil {
		log.Println("Error writing youtube-dl db:", err)
		return
	}
	var streamer settings.Streamer
	for configIndex := range system.System.DB.Streamers.List {
		streamer = system.System.DB.Streamers.List[configIndex]
	}

	b.mux.Lock()
	b.AddProcess(streamer.Provider, streamer.Name)
	b.mux.Unlock()
	var wg sync.WaitGroup
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	b.bots[streamerName].Website.Interface.TrueName(streamerName)

	wg.Add(1)
	// Pass the index and streamer name into the closure to avoid capture issues.
	go func(bot *recorder.Recorder, sName string) {
		defer wg.Done()
		for {
			select {
			case <-b.ctx.Done():
			case <-ticker.C:
				if bot.IsRecording {
					continue
				}
				db.LoadConfig(settings.CONFIG_SETTINGS_PATH, &system.System.DB.Settings)
				log.Printf("Checking %s online status", bot.Website.Username)
				if !bot.Website.Interface.IsOnline(sName) {
					log.Printf("Not online.")
					return
				}
				log.Printf("Online! Starting new recording", sName)
				bot.StartRecording(sName)

				b.mux.Lock()
				bot.IsRecording = false
				bot.Stop()
				b.mux.Unlock()

				log.Printf("Recording for %s finished", sName)
			}
		}
	}(b.bots[streamerName], streamerName)

	wg.Wait()
}
