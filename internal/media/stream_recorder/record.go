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

	var wg sync.WaitGroup
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	b.bots[streamerName].Website.Interface.TrueName(streamerName)

	wg.Add(1)
	// Pass the index and streamer name into the closure to avoid capture issues.
	go func(bot *recorder.Recorder, streamerName string) {
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
				if !bot.Website.Interface.IsOnline(streamerName) {
					log.Printf("Not online.")
					return
				}
				log.Println("Online! Starting new recording", streamerName)

				bot.StartRecording(streamerName)

				b.mux.Lock()
				bot.IsRecording = false
				bot.Stop()
				b.mux.Unlock()

				log.Printf("Recording for %s finished", streamerName)
			}
		}
	}(b.bots[streamerName], streamerName)

	wg.Wait()
}
