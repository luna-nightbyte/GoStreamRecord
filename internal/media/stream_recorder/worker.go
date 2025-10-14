package stream_recorder

import (
	"context"
	"fmt"
	"log"
	"remoteCtrl/internal/media/stream_recorder/recorder"
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

	if b.bots[name] == nil {
		log.Println("Bot was not started.")
		return
	}
	if len(command) == 0 {
		log.Println("No command provided..")
		return
	}
	switch strings.ToLower(command) {
	case "start":
		fmt.Println("Starting bot for", name)
		// If the bot was previously stopped, reinitialize the context.
		if b.ctx.Err() != nil {
			b.ctx, b.cancel = context.WithCancel(context.Background())
		}

		if b.bots[name] != nil {
			if b.bots[name].Cmd != nil {
				log.Println("Alredy recording video from '%s'", name)
			}
		}
		// for _, s := range b.status {
		// 	if name == s.Website.Username && s.Cmd != nil {
		// 		fmt.Println("Bot already running for", name)
		// 		log.Println("Bot already running..")
		// 		return
		// 	}
		// }
		fmt.Println("Starting bot")
		log.Println("Starting bot")
		go b.RecordLoop(name)
		b.bots[name].IsRestarting = false
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
