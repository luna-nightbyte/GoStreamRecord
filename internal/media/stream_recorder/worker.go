package stream_recorder

import (
	"context"
	"fmt"
	"log"
	"remoteCtrl/internal/media/stream_recorder/recorder"
	"remoteCtrl/internal/system"
	"remoteCtrl/internal/system/settings"
	"remoteCtrl/internal/utils"
	"strings"
	"sync"
)

// Bot encapsulates the recording bot’s state.
type Controller struct {
	mux        sync.Mutex
	bots       map[string]*recorder.Recorder
	_status    []recorder.Recorder
	isFirstRun bool
	logger     *log.Logger
	// ctx is used to signal shutdown.
	ctx    context.Context
	cancel context.CancelFunc
}

var Streamer Controller = NewBot()

// NewBot creates a new Bot, sets up its cancellation context.
func NewBot() Controller {
	ctx, cancel := context.WithCancel(context.Background())
	b := Controller{
		ctx:    ctx,
		cancel: cancel,
		//status:     []recorder.Recorder{},
		isFirstRun: true,
		bots:       make(map[string]*recorder.Recorder),
	}

	return b
}

func (b *Controller) Execute(command string, name string) {

	if len(command) == 0 {
		fmt.Println("No command provided..")
		log.Println("No command provided..")
		return
	}
	switch strings.ToLower(command) {
	case "repair":
		log.Println("Starting video codec verification. This might take some time depending on how many videos you have and their lenght/quality.")
		fmt.Println("Starting video codec verification. This might take some time depending on how many videos you have and their lenght/quality.")
		utils.VideoVerify.RunCodecVerification()
		log.Println("Done!")
		fmt.Println("Done!")
	case "start":
		// If the bot was previously stopped, reinitialize the context.
		if b.ctx.Err() != nil {
			b.ctx, b.cancel = context.WithCancel(context.Background())
		}

		if b.bots[name] != nil {
			if b.bots[name].Cmd != nil {
				log.Println("Alredy recording video from '%s'", name)
			}
			fmt.Println("Starting bot")
			log.Println("Starting bot")
			go b.RecordLoop(name)
			b.bots[name].IsRestarting = false
		} else {
			fmt.Println("Starting bot for", name)
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
			b.Execute("start", name)
			return
		}
		// for _, s := range b.status {
		// 	if name == s.Website.Username && s.Cmd != nil {
		// 		fmt.Println("Bot already running for", name)
		// 		log.Println("Bot already running..")
		// 		return
		// 	}
		// }
	case "stop":
		botsAvailable := len(b.bots) != 0
		if !botsAvailable {
			log.Println("No more bots recording")
			break
		}
		log.Println("Stopping recording for", name)
		b.mux.Lock()
		b.stopProcessIfRunning(b.bots[name])
		for _, s := range b.bots {
			// Stop only the specified process (or all if name is empty).
			if name == "" || s.Website.Username == name {
				b.stopProcessIfRunning(s)
			} else {
				log.Println("Not stopping..")
			}
		}

		b.mux.Unlock()
		b.checkProcesses()
		if b.bots[name].IsRestarting {
			b.Execute("start", name)
		}
	case "restart":
		log.Println("Restarting bot")
		// b.ctx, b.cancel = context.WithCancel(context.Background())
		b.mux.Lock()
		b.bots[name].IsRestarting = true
		b.mux.Unlock()
		b.Execute("stop", name)
		break
	}

}
