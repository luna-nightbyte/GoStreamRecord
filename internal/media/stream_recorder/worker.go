package stream_recorder

import (
	"context"
	"fmt"
	"log"
	"remoteCtrl/internal/media/stream_recorder/recorder"
	"remoteCtrl/internal/system"
	"remoteCtrl/internal/system/settings"
	"remoteCtrl/internal/utils"
	"remoteCtrl/internal/web/handlers/status"
	"strings"
	"sync"
	"time"
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
		log.Println("No command provided..")
		return
	}
	switch strings.ToLower(command) {
	case "repair":

		status.Status.Is_Fixing_Codec = true
		log.Println("Starting video codec verification. This might take some time depending on how many videos you have and their lenght/quality.")
		utils.VideoVerify.RunCodecVerification()
		
		status.Status.Is_Fixing_Codec = false
	case "start":
		// If the bot was previously stopped, reinitialize the context.
		if b.ctx.Err() != nil {
			b.ctx, b.cancel = context.WithCancel(context.Background())
		}

		if b.bots[name] != nil {
			if b.bots[name].Cmd != nil {
				log.Println("Alredy recording video from '%s'", name)
			}
			log.Println("Starting bot")
			status.Status.Is_Recording = true
			b.bots[name].Start()


			go b.bots[name].StartRecordTicker(b.ctx)

		} else { 
			if err := b.writeYoutubeDLdb(); err != nil {
				log.Println("Error writing youtube-dl db:", err)
				return
			}
			var streamer settings.Streamer
			found := false
			for _, streamConf := range system.System.Config.Streamers.List {
				if streamConf.Name == name {
					streamer = streamConf
					found = true
					break
				}
			}
			if !found {
				log.Println("Error retrieving streamer from the config..")
				return
			}
			b.mux.Lock()
			b.AddProcess(streamer.Provider, streamer.Name)
			b.mux.Unlock()

			b.bots[name].IsRestarting = false
			go b.Execute("start", name)
			return
		}
	case "stop":
		botsAvailable := len(b.bots) != 0
		if !botsAvailable {
			log.Println("No more bots recording")
			break
		}
		log.Println("Stopping recording for", name)
		b.stopProcessIfRunning(b.bots[name])

		if b.bots[name].IsRestarting {
			time.Sleep(1 * time.Second)
			fmt.Println("Starting again")
			b.bots[name].Start()
			break
		}
		b.bots[name].StopTicker()
		b.checkProcesses()
		if len(b.ListRecorders()) == 0 {
			status.Status.Is_Recording = false
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
